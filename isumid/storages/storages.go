package storages

import (
	"fmt"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack/v5"
	_ "modernc.org/sqlite"
)

type Storage struct {
	outputDir string
	db        *sqlx.DB
}

func New(setting settings.Setting) (Storage, error) {
	outputDir := path.Join(setting.OutputDir, "body")

	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if err != nil {
		return Storage{}, nil
	}

	db, err := sqlx.Open("sqlite", filepath.Join(setting.OutputDir, "isumid.sqlite"))
	if err != nil {
		return Storage{}, err
	}

	query := `
    CREATE TABLE IF NOT EXISTS metadata (
        method TEXT,
        url TEXT,
        reqHeader BLOB,
        statusCode INTEGER,
        resHeader BLOB,
        isReqText BOOLEAN,
        isResText BOOLEAN,
        ulid TEXT PRIMARY KEY,
        reqLength INTEGER,
        resLength INTEGER
    );
	`
	_, err = db.Exec(query)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to create metadata table: %w", err)
	}

	return Storage{outputDir: outputDir, db: db}, nil
}

func (s Storage) Close() error {
	return s.db.Close()
}

func IsText(header map[string][]string, body []byte) bool {
	mediaTypeExpected := []string{"text/plain", "text/csv", "text/html", "text/css", "text/javascript", "application/json", "application/x-www-form-urlencoded"}

	contentType, ok := header["Content-Type"]
	if !ok {
		if len(body) == 0 {
			return false
		}
		contentType = []string{http.DetectContentType(body)}
	}
	for _, c := range contentType {
		mediaTypeActual, _, err := mime.ParseMediaType(c)
		if err != nil {
			continue
		}

		for _, e := range mediaTypeExpected {
			if mediaTypeActual == e {
				return true
			}
		}
	}
	return false
}

func genUlidStr() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

func serializeMap(header map[string][]string) ([]byte, error) {
	return msgpack.Marshal(header)
}

func deserializeMap(data []byte) (map[string][]string, error) {
	var header map[string][]string
	if err := msgpack.Unmarshal(data, &header); err != nil {
		return nil, err
	}
	return header, nil
}

type serializedMeta struct {
	Method     string `db:"method"`
	Url        string `db:"url"`
	ReqHeader  []byte `db:"reqHeader"`
	StatusCode int    `db:"statusCode"`
	ResHeader  []byte `db:"resHeader"`
	IsReqText  bool   `db:"isReqText"`
	IsResText  bool   `db:"isResText"`
	ReqLength  int    `db:"reqLength"`
	ResLength  int    `db:"resLength"`
	Ulid       string `db:"ulid"`
}

func serializeMeta(meta models.Meta) (serializedMeta, error) {
	reqHeader, err := serializeMap(meta.ReqHeader)
	if err != nil {
		return serializedMeta{}, err
	}
	resHeader, err := serializeMap(meta.ResHeader)
	if err != nil {
		return serializedMeta{}, err
	}
	return serializedMeta{
		Method:     meta.Method,
		Url:        meta.Url,
		ReqHeader:  reqHeader,
		StatusCode: meta.StatusCode,
		ResHeader:  resHeader,
		IsReqText:  meta.IsReqText,
		IsResText:  meta.IsResText,
		ReqLength:  meta.ReqLength,
		ResLength:  meta.ResLength,
		Ulid:       meta.Ulid,
	}, nil
}

func deserializeMeta(data serializedMeta) (models.Meta, error) {
	reqHeader, err := deserializeMap(data.ReqHeader)
	if err != nil {
		return models.Meta{}, err
	}
	resHeader, err := deserializeMap(data.ResHeader)
	if err != nil {
		return models.Meta{}, err
	}
	return models.Meta{
		Method:     data.Method,
		Url:        data.Url,
		ReqHeader:  reqHeader,
		StatusCode: data.StatusCode,
		ResHeader:  resHeader,
		IsReqText:  data.IsReqText,
		IsResText:  data.IsResText,
		Ulid:       data.Ulid,
		ReqLength:  data.ReqLength,
		ResLength:  data.ResLength,
	}, nil
}

func (s Storage) Save(data models.RecordedDataInput) error {
	// generate ulid
	ulidStr := genUlidStr()

	// save metadata
	{
		meta := models.Meta{
			Method:     data.Method,
			Url:        data.Url,
			ReqHeader:  data.ReqHeader,
			StatusCode: data.StatusCode,
			ResHeader:  data.ResHeader,
			IsReqText:  IsText(data.ReqHeader, data.ReqBody),
			IsResText:  IsText(data.ResHeader, data.ResBody),
			Ulid:       ulidStr,
			ReqLength:  len(data.ReqBody),
			ResLength:  len(data.ResBody),
		}
		query := `
			INSERT INTO metadata (method, url, reqHeader, statusCode, resHeader, isReqText, isResText, ulid, reqLength, resLength)
			VALUES (:method, :url, :reqHeader, :statusCode, :resHeader, :isReqText, :isResText, :ulid, :reqLength, :resLength);
		`
		serialized, err := serializeMeta(meta)
		if err != nil {
			return err
		}
		if _, err := s.db.NamedExec(query, serialized); err != nil {
			return err
		}
	}

	// save request body data
	{
		path := filepath.Join(s.outputDir, ulidStr+".req.body")
		err := os.WriteFile(path, data.ReqBody, 0666)
		if err != nil {
			return err
		}
	}

	// save response body data
	{
		path := filepath.Join(s.outputDir, ulidStr+".res.body")
		err := os.WriteFile(path, data.ResBody, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) FetchMeta(ulid string) (models.Meta, error) {
	var serializedMeta serializedMeta
	query := `SELECT * FROM metadata WHERE ulid = ?`
	err := s.db.Get(&serializedMeta, query, ulid)
	if err != nil {
		return models.Meta{}, err
	}

	meta, err := deserializeMeta(serializedMeta)
	if err != nil {
		return models.Meta{}, err
	}

	return meta, nil
}

func (s Storage) FetchMetaList(offset, length int) ([]models.Meta, error) {
	query := `SELECT * FROM metadata LIMIT ? OFFSET ?`
	var serializedMetaList []serializedMeta
	err := s.db.Select(&serializedMetaList, query, length, offset)
	if err != nil {
		return nil, err
	}

	metaList := []models.Meta{}
	for _, m := range serializedMetaList {
		meta, err := deserializeMeta(m)
		if err != nil {
			return nil, err
		}
		metaList = append(metaList, meta)
	}

	return metaList, nil
}

func (s Storage) SearchMetaList(urlQuery string, offset int, length int) ([]models.Meta, int, error) {
	metaList := []models.Meta{}
	{
		query := `SELECT * FROM metadata WHERE url LIKE ? LIMIT ? OFFSET ?`
		var serializedMetaList []serializedMeta
		err := s.db.Select(&serializedMetaList, query, urlQuery, length, offset)
		if err != nil {
			return nil, 0, err
		}

		for _, m := range serializedMetaList {
			meta, err := deserializeMeta(m)
			if err != nil {
				return nil, 0, err
			}
			metaList = append(metaList, meta)
		}
	}

	var totalHit int
	{
		query := `SELECT COUNT(*) FROM metadata WHERE url LIKE ?`
		err := s.db.Get(&totalHit, query, urlQuery)
		if err != nil {
			return nil, 0, err
		}
	}

	return metaList, totalHit, nil
}

func (s Storage) fetchFile(fileName string) ([]byte, error) {
	body, err := os.ReadFile(filepath.Join(s.outputDir, fileName))
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s Storage) FetchReqBody(ulid string) ([]byte, error) {
	return s.fetchFile(ulid + ".req.body")
}

func (s Storage) FetchResBody(ulid string) ([]byte, error) {
	return s.fetchFile(ulid + ".res.body")
}

func (s Storage) FetchReproducedBody(ulid string) ([]byte, error) {
	return s.fetchFile(ulid + ".reproduced.body")
}

func (s Storage) FetchReproducedHeader(ulid string) (map[string][]string, error) {
	data, err := os.ReadFile(filepath.Join(s.outputDir, ulid+".reproduced.header"))
	if err != nil {
		return nil, err
	}

	var header map[string][]string
	err = msgpack.Unmarshal(data, &header)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (s Storage) SaveReproduced(ulid string, body []byte, header map[string][]string) error {
	{
		path := filepath.Join(s.outputDir, ulid+".reproduced.header")
		data, err := msgpack.Marshal(header)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			return err
		}
	}
	{
		path := filepath.Join(s.outputDir, ulid+".reproduced.body")
		err := os.WriteFile(path, body, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) CreateDir() error {
	return os.MkdirAll(s.outputDir, 0777)
}

func (s Storage) RemoveAll() error {
	query := `DELETE FROM metadata`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return os.RemoveAll(s.outputDir)
}

func (s Storage) Remove(ulid string) error {
	query := `DELETE FROM metadata WHERE ulid = ?`
	_, err := s.db.Exec(query, ulid)
	if err != nil {
		return err
	}

	fileList, err := filepath.Glob(filepath.Join(s.outputDir, ulid+"*"))
	if err != nil {
		return err
	}

	for _, filePath := range fileList {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) FetchTotalTransactions() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM metadata`
	err := s.db.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
