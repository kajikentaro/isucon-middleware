package storages

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
)

// metaSizeExpr estimates the size of a metadata row in bytes.
// 40 bytes are for the numeric columns.
const metaSizeExpr = `(length(CAST(method AS BLOB)) + length(CAST(url AS BLOB)) + length(CAST(path AS BLOB))
	+ length(CAST(ulid AS BLOB)) + length(CAST(reqHeader AS BLOB)) + length(CAST(resHeader AS BLOB)) + 40)`

// Columns only the exported file has, for analysis. They are filled after the rows are copied.
const exportExtraColumns = `
	NULL AS reqBody,
	NULL AS resBody,
	NULL AS cookie,
	NULL AS setCookie,
	NULL AS reqContentType,
	NULL AS resContentType,
	NULL AS resContentEncoding
`

const exportFillHeaderColumns = `
	UPDATE metadata SET
		cookie             = (SELECT group_concat(value, '; ')     FROM json_each(reqHeader, '$.Cookie')),
		setCookie          = (SELECT group_concat(value, char(10)) FROM json_each(resHeader, '$.Set-Cookie')),
		reqContentType     = json_extract(reqHeader, '$.Content-Type[0]'),
		resContentType     = json_extract(resHeader, '$.Content-Type[0]'),
		resContentEncoding = json_extract(resHeader, '$.Content-Encoding[0]');
	UPDATE metadata SET reqBody = '' WHERE reqLength = 0;
	UPDATE metadata SET resBody = '' WHERE resLength = 0;
`

const exportIndexes = `
	CREATE INDEX metadata_path ON metadata (path, statusCode);
	CREATE INDEX metadata_startedAtUs ON metadata (startedAtUs);
	CREATE INDEX metadata_cookie ON metadata (cookie, startedAtUs);
	CREATE INDEX metadata_setCookie ON metadata (setCookie);
`

// ExportSize returns the size of the whole recorded data, to choose the limits of Export.
func (s Storage) ExportSize() (models.ExportSize, error) {
	s.Flush()

	var size models.ExportSize
	query := `
		SELECT
			COUNT(*) AS count,
			COALESCE(SUM(` + metaSizeExpr + `), 0) AS metaBytes,
			COALESCE(SUM(CASE WHEN isReqText THEN reqLength ELSE 0 END + CASE WHEN isResText THEN resLength ELSE 0 END), 0) AS textBodyBytes
		FROM metadata`
	if err := s.db.Get(&size, query); err != nil {
		return models.ExportSize{}, err
	}
	return size, nil
}

// Export writes the recorded data to a new SQLite file dst, for analysis.
// Metadata rows are written from the oldest until their estimated size reaches maxMetaBytes.
// Then text bodies of the written rows are embedded from the oldest until their total size reaches maxBodyBytes.
// The original database is not modified.
func (s Storage) Export(dst string, maxMetaBytes, maxBodyBytes int64) error {
	s.Flush()

	if err := s.exportMeta(dst, maxMetaBytes); err != nil {
		return err
	}

	out, err := sqlx.Open("sqlite", dst)
	if err != nil {
		return err
	}
	defer out.Close()
	out.SetMaxOpenConns(1)

	if _, err := out.Exec(exportFillHeaderColumns); err != nil {
		return err
	}
	if err := s.exportBodies(out, maxBodyBytes); err != nil {
		return err
	}
	if _, err := out.Exec(exportIndexes); err != nil {
		return err
	}
	return out.Close()
}

// exportMeta creates the metadata table in dst (a new database), and copies the oldest rows that fit in maxMetaBytes.
func (s Storage) exportMeta(dst string, maxMetaBytes int64) error {
	ctx := context.Background()
	// ATTACH is per connection, so run everything on one connection
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	var lastUlid sql.NullString
	query := `
		SELECT ulid FROM (
			SELECT ulid, SUM(` + metaSizeExpr + `) OVER (ORDER BY ulid) AS acc FROM metadata
		) WHERE acc <= ? ORDER BY ulid DESC LIMIT 1`
	if err := conn.GetContext(ctx, &lastUlid, query, maxMetaBytes); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if _, err := conn.ExecContext(ctx, "ATTACH DATABASE ? AS export", dst); err != nil {
		return err
	}
	copyErr := copyMeta(ctx, conn, lastUlid)
	// always detach, as the connection goes back to the pool
	_, detachErr := conn.ExecContext(ctx, "DETACH DATABASE export")
	return errors.Join(copyErr, detachErr)
}

// copyMeta copies the rows up to lastUlid to the attached database "export".
// If lastUlid is NULL (no row fits), the table is created with no row.
func copyMeta(ctx context.Context, conn *sqlx.Conn, lastUlid sql.NullString) error {
	// copy the columns as they are in the original table, so that this does not depend on its schema.
	query := "CREATE TABLE export.metadata AS SELECT *," + exportExtraColumns + "FROM main.metadata WHERE ulid <= ? ORDER BY ulid"
	if _, err := conn.ExecContext(ctx, query, lastUlid); err != nil {
		return err
	}
	// CREATE TABLE AS does not copy the primary key. It is needed to update the bodies by ulid.
	_, err := conn.ExecContext(ctx, "CREATE UNIQUE INDEX export.metadata_ulid ON metadata (ulid)")
	return err
}

// exportBodies embeds text bodies from the oldest row, and stops at the first row that does not fit in maxBodyBytes.
// Bodies compressed with gzip (Content-Encoding: gzip) are decompressed, and counted by the decompressed size.
// Bodies not embedded (binary, or after the limit) are left NULL.
func (s Storage) exportBodies(out *sqlx.DB, maxBodyBytes int64) error {
	var rows []struct {
		Ulid               string `db:"ulid"`
		IsReqText          bool   `db:"isReqText"`
		IsResText          bool   `db:"isResText"`
		ReqLength          int64  `db:"reqLength"`
		ResLength          int64  `db:"resLength"`
		ReqContentEncoding string `db:"reqContentEncoding"`
		ResContentEncoding string `db:"resContentEncoding"`
	}
	query := `
		SELECT ulid, isReqText, isResText, reqLength, resLength,
			IFNULL(json_extract(reqHeader, '$.Content-Encoding[0]'), '') AS reqContentEncoding,
			IFNULL(resContentEncoding, '') AS resContentEncoding
		FROM metadata ORDER BY ulid`
	if err := out.Select(&rows, query); err != nil {
		return err
	}

	tx, err := out.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	remaining := maxBodyBytes
	for _, row := range rows {
		// the request and the response of a row are embedded together, or not at all
		var reqBody, resBody *string
		if row.IsReqText && row.ReqLength > 0 {
			body, fits, err := s.readTextBody(row.Ulid+".req.body", row.ReqLength, isGzip(row.ReqContentEncoding), remaining)
			if err != nil {
				return err
			}
			if !fits {
				break
			}
			reqBody = body
		}
		if row.IsResText && row.ResLength > 0 {
			body, fits, err := s.readTextBody(row.Ulid+".res.body", row.ResLength, isGzip(row.ResContentEncoding), remaining-strLen(reqBody))
			if err != nil {
				return err
			}
			if !fits {
				break
			}
			resBody = body
		}
		if reqBody == nil && resBody == nil {
			continue
		}

		query := "UPDATE metadata SET reqBody = COALESCE(?, reqBody), resBody = COALESCE(?, resBody) WHERE ulid = ?"
		if _, err := tx.Exec(query, reqBody, resBody, row.Ulid); err != nil {
			return err
		}
		remaining -= strLen(reqBody) + strLen(resBody)
	}
	return tx.Commit()
}

var errTooLarge = errors.New("too large")

// readTextBody reads a body to be stored as TEXT, decompressing it if isGzip.
// fits is false if the body is larger than limit (then the file is not read, unless it is gzip).
// body is nil if it cannot be embedded: missing file, broken gzip, or invalid UTF-8.
func (s Storage) readTextBody(fileName string, length int64, isGzip bool, limit int64) (body *string, fits bool, err error) {
	if !isGzip && length > limit {
		return nil, false, nil
	}

	data, err := s.fetchFile(fileName)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to read %s: %w", fileName, err)
	}

	if isGzip {
		data, err = gunzip(data, limit)
		if errors.Is(err, errTooLarge) {
			return nil, false, nil
		}
		if err != nil {
			return nil, true, nil
		}
	}

	if !utf8.Valid(data) {
		return nil, true, nil
	}
	str := string(data)
	return &str, true, nil
}

// gunzip decompresses data. It reads at most limit bytes, and returns errTooLarge if the result is larger,
// so that a highly compressed body does not use up the memory.
func gunzip(data []byte, limit int64) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	decoded, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(decoded)) > limit {
		return nil, errTooLarge
	}
	return decoded, nil
}

func isGzip(contentEncoding string) bool {
	return strings.EqualFold(strings.TrimSpace(contentEncoding), "gzip")
}

func strLen(s *string) int64 {
	if s == nil {
		return 0
	}
	return int64(len(*s))
}
