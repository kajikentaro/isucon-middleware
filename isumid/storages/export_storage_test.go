package storages

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/stretchr/testify/assert"
)

type exportedRow struct {
	Url                string         `db:"url"`
	ReqBody            sql.NullString `db:"reqBody"`
	ResBody            sql.NullString `db:"resBody"`
	ReqBodyType        string         `db:"reqBodyType"`
	ResBodyType        string         `db:"resBodyType"`
	Cookie             sql.NullString `db:"cookie"`
	SetCookie          sql.NullString `db:"setCookie"`
	ReqContentType     sql.NullString `db:"reqContentType"`
	ResContentType     sql.NullString `db:"resContentType"`
	ResContentEncoding sql.NullString `db:"resContentEncoding"`
}

func saveForExport(t *testing.T, s Storage, url string, reqHeader, resHeader map[string][]string, reqBody, resBody string) {
	err := s.Save(models.RecordedDataInput{
		Method:     "POST",
		Url:        url,
		Path:       url,
		ReqHeader:  reqHeader,
		ReqBody:    []byte(reqBody),
		StatusCode: 200,
		ResHeader:  resHeader,
		ResBody:    []byte(resBody),
		StartedAt:  time.Now(),
	})
	assert.NoError(t, err)
	// keep the ulid order same as the save order
	time.Sleep(2 * time.Millisecond)
}

func export(t *testing.T, s Storage, maxMetaBytes, maxBodyBytes int64) []exportedRow {
	dst := filepath.Join(t.TempDir(), "export.sqlite")
	assert.NoError(t, s.Export(dst, maxMetaBytes, maxBodyBytes))

	out, err := sqlx.Open("sqlite", dst)
	assert.NoError(t, err)
	defer out.Close()

	var rows []exportedRow
	err = out.Select(&rows, `
		SELECT url, reqBody, resBody, typeof(reqBody) AS reqBodyType, typeof(resBody) AS resBodyType,
			cookie, setCookie, reqContentType, resContentType, resContentEncoding
		FROM metadata ORDER BY ulid`)
	assert.NoError(t, err)
	return rows
}

var textHeader = map[string][]string{"Content-Type": {"text/plain"}}

func TestExportBodies(t *testing.T) {
	s := newTestStorage(t)
	big := strings.Repeat("x", 1000)
	saveForExport(t, s, "/text", textHeader, textHeader, "req", "res")
	saveForExport(t, s, "/binary", map[string][]string{"Content-Type": {"image/png"}}, textHeader, "\x89PNG", "ok")
	saveForExport(t, s, "/empty", textHeader, textHeader, "", "")
	saveForExport(t, s, "/invalid-utf8", textHeader, textHeader, "a\xffb", "ok")
	saveForExport(t, s, "/big", textHeader, textHeader, big, "")
	saveForExport(t, s, "/small-after-big", textHeader, textHeader, "s", "s")

	t.Run("embeds text bodies until the limit", func(t *testing.T) {
		// enough for everything before /big, but not for /big
		rows := export(t, s, 1<<30, 500)
		assert.Len(t, rows, 6)

		assert.Equal(t, sql.NullString{String: "req", Valid: true}, rows[0].ReqBody)
		assert.Equal(t, sql.NullString{String: "res", Valid: true}, rows[0].ResBody)
		assert.Equal(t, "text", rows[0].ReqBodyType)

		assert.False(t, rows[1].ReqBody.Valid, "binary body is not embedded")
		assert.Equal(t, sql.NullString{String: "ok", Valid: true}, rows[1].ResBody)

		assert.Equal(t, sql.NullString{String: "", Valid: true}, rows[2].ReqBody, "empty body is ''")
		assert.Equal(t, sql.NullString{String: "", Valid: true}, rows[2].ResBody)

		assert.False(t, rows[3].ReqBody.Valid, "invalid UTF-8 is not embedded")
		assert.Equal(t, sql.NullString{String: "ok", Valid: true}, rows[3].ResBody)

		// stops at /big, even though the next row is small
		assert.False(t, rows[4].ReqBody.Valid)
		assert.Equal(t, sql.NullString{String: "", Valid: true}, rows[4].ResBody)
		assert.False(t, rows[5].ReqBody.Valid)
		assert.False(t, rows[5].ResBody.Valid)
	})

	t.Run("embeds everything with a large limit", func(t *testing.T) {
		rows := export(t, s, 1<<30, 1<<30)
		assert.Equal(t, big, rows[4].ReqBody.String)
		assert.Equal(t, "s", rows[5].ReqBody.String)
	})

	t.Run("no bodies with a zero limit", func(t *testing.T) {
		rows := export(t, s, 1<<30, 0)
		for _, row := range rows {
			assert.True(t, !row.ReqBody.Valid || row.ReqBody.String == "", row.Url)
			assert.True(t, !row.ResBody.Valid || row.ResBody.String == "", row.Url)
		}
	})
}

func TestExportMetaLimit(t *testing.T) {
	s := newTestStorage(t)
	for _, url := range []string{"/1", "/2", "/3"} {
		saveForExport(t, s, url, textHeader, textHeader, "", "")
	}

	size, err := s.ExportSize()
	assert.NoError(t, err)
	perRow := size.MetaBytes / 3

	assert.Len(t, export(t, s, 0, 0), 0)
	assert.Len(t, export(t, s, perRow-1, 0), 0)

	rows := export(t, s, perRow*2, 0)
	assert.Len(t, rows, 2)
	assert.Equal(t, "/1", rows[0].Url, "the oldest rows are exported")
	assert.Equal(t, "/2", rows[1].Url)

	assert.Len(t, export(t, s, size.MetaBytes, 0), 3)
}

func TestExportHeaderColumns(t *testing.T) {
	s := newTestStorage(t)
	saveForExport(t, s, "/headers",
		map[string][]string{"Cookie": {"a=1", "b=2"}, "Content-Type": {"application/json"}},
		map[string][]string{"Set-Cookie": {"a=1; Path=/", "b=2"}, "Content-Type": {"text/html"}, "Content-Encoding": {"gzip"}},
		"{}", "")
	saveForExport(t, s, "/no-headers", map[string][]string{}, map[string][]string{}, "", "")

	rows := export(t, s, 1<<30, 1<<30)
	assert.Equal(t, sql.NullString{String: "a=1; b=2", Valid: true}, rows[0].Cookie)
	assert.Equal(t, sql.NullString{String: "a=1; Path=/\nb=2", Valid: true}, rows[0].SetCookie)
	assert.Equal(t, sql.NullString{String: "application/json", Valid: true}, rows[0].ReqContentType)
	assert.Equal(t, sql.NullString{String: "text/html", Valid: true}, rows[0].ResContentType)
	assert.Equal(t, sql.NullString{String: "gzip", Valid: true}, rows[0].ResContentEncoding)

	assert.False(t, rows[1].Cookie.Valid)
	assert.False(t, rows[1].SetCookie.Valid)
	assert.False(t, rows[1].ReqContentType.Valid)
	assert.False(t, rows[1].ResContentType.Valid)
	assert.False(t, rows[1].ResContentEncoding.Valid)
}

func TestExportIndexesAndOriginal(t *testing.T) {
	s := newTestStorage(t)
	saveForExport(t, s, "/", textHeader, textHeader, "req", "res")

	dst := filepath.Join(t.TempDir(), "export.sqlite")
	assert.NoError(t, s.Export(dst, 1<<30, 1<<30))
	out, err := sqlx.Open("sqlite", dst)
	assert.NoError(t, err)
	defer out.Close()

	var indexes []string
	assert.NoError(t, out.Select(&indexes, "SELECT name FROM sqlite_master WHERE type = 'index' AND sql IS NOT NULL ORDER BY name"))
	assert.Equal(t, []string{"metadata_cookie", "metadata_path", "metadata_setCookie", "metadata_startedAtUs", "metadata_ulid"}, indexes)

	// the original database is not changed
	var columns int
	assert.NoError(t, s.db.Get(&columns, "SELECT COUNT(*) FROM pragma_table_info('metadata') WHERE name = 'reqBody'"))
	assert.Equal(t, 0, columns)
	var databases int
	assert.NoError(t, s.db.Get(&databases, "SELECT COUNT(*) FROM pragma_database_list WHERE name = 'export'"))
	assert.Equal(t, 0, databases)

	// can export again
	assert.NoError(t, s.Export(filepath.Join(t.TempDir(), "export2.sqlite"), 1<<30, 1<<30))
}

func TestExportSize(t *testing.T) {
	s := newTestStorage(t)
	saveForExport(t, s, "/text", textHeader, textHeader, "req", "res!")
	saveForExport(t, s, "/binary", map[string][]string{"Content-Type": {"image/png"}}, textHeader, "\x89PNG", "ok")

	size, err := s.ExportSize()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), size.Count)
	assert.Equal(t, int64(3+4+2), size.TextBodyBytes)
	assert.Greater(t, size.MetaBytes, int64(0))
}

func gzipString(t *testing.T, s string) string {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write([]byte(s))
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	return buf.String()
}

func TestExportGzipBodies(t *testing.T) {
	s := newTestStorage(t)
	jsonGzip := map[string][]string{"Content-Type": {"application/json"}, "Content-Encoding": {"gzip"}}
	big := strings.Repeat("a", 1000) // compressed to a few dozen bytes

	saveForExport(t, s, "/gzip-res", textHeader, jsonGzip, "req", gzipString(t, `{"ok":true}`))
	saveForExport(t, s, "/gzip-req", jsonGzip, textHeader, gzipString(t, `{"req":1}`), "res")
	saveForExport(t, s, "/broken-gzip", textHeader, jsonGzip, "", "not gzip")
	saveForExport(t, s, "/gzip-big", textHeader, jsonGzip, "", gzipString(t, big))
	saveForExport(t, s, "/after-big", textHeader, textHeader, "s", "s")

	t.Run("decompresses gzip bodies", func(t *testing.T) {
		rows := export(t, s, 1<<30, 1<<30)
		assert.Equal(t, sql.NullString{String: `{"ok":true}`, Valid: true}, rows[0].ResBody)
		assert.Equal(t, "text", rows[0].ResBodyType)
		assert.Equal(t, sql.NullString{String: "gzip", Valid: true}, rows[0].ResContentEncoding, "the header is kept as recorded")
		assert.Equal(t, sql.NullString{String: `{"req":1}`, Valid: true}, rows[1].ReqBody)
		assert.False(t, rows[2].ResBody.Valid, "broken gzip is not embedded")
		assert.Equal(t, big, rows[3].ResBody.String)
	})

	t.Run("counts the decompressed size", func(t *testing.T) {
		// enough for the compressed size of /gzip-big, but not for the decompressed size
		rows := export(t, s, 1<<30, 500)
		assert.True(t, rows[0].ResBody.Valid)
		assert.True(t, rows[1].ReqBody.Valid)
		assert.False(t, rows[3].ResBody.Valid, "stops at /gzip-big")
		assert.False(t, rows[4].ReqBody.Valid, "and does not embed the bodies after it")
	})
}
