package storages

import (
	"fmt"
	"io/fs"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack/v5"
)

type Storage struct {
	models.Setting
}

type RecordedDataInput struct {
	Method    string
	Url       string
	ReqHeader map[string][]string
	ReqBody   []byte

	StatusCode int
	ResHeader  map[string][]string
	ResBody    []byte
}

type Meta struct {
	Method    string
	Url       string
	ReqHeader map[string][]string

	StatusCode int
	ResHeader  map[string][]string

	IsReqText bool
	IsResText bool
	ReqLength int
	ResLength int
	Ulid      string
}

func New(setting models.Setting) Storage {
	return Storage{Setting: setting}
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

func (s Storage) Save(data RecordedDataInput) error {
	err := os.MkdirAll(s.OutputDir, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// generate ulid
	ulidStr := genUlidStr()

	// save metadata
	{
		path := filepath.Join(s.OutputDir, ulidStr+".meta")
		if err != nil {
			return err
		}
		meta := Meta{
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
		data, err := msgpack.Marshal(meta)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			return err
		}
	}

	// save request body data
	{
		path := filepath.Join(s.OutputDir, ulidStr+".req.body")
		if err != nil {
			return err
		}
		err = os.WriteFile(path, data.ReqBody, 0666)
		if err != nil {
			return err
		}
	}

	// save response body data
	{
		path := filepath.Join(s.OutputDir, ulidStr+".res.body")
		if err != nil {
			return err
		}
		err = os.WriteFile(path, data.ResBody, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) FetchMeta(ulid string) (Meta, error) {
	data, err := os.ReadFile(filepath.Join(s.OutputDir, ulid+".meta"))
	if err != nil {
		return Meta{}, err
	}

	var meta Meta
	err = msgpack.Unmarshal(data, &meta)
	if err != nil {
		return Meta{}, err
	}

	return meta, nil
}

func (s Storage) FetchMetaList(offset, length int) ([]Meta, error) {
	fileList, err := os.ReadDir(s.OutputDir)
	if err != nil {
		return nil, err
	}

	metaList := []fs.DirEntry{}
	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".meta" {
			continue
		}
		metaList = append(metaList, file)
	}

	res := []Meta{}
	for idx, file := range metaList {
		if idx < offset {
			continue
		}
		if offset+length <= idx {
			break
		}

		data, err := os.ReadFile(filepath.Join(s.OutputDir, file.Name()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		var meta Meta
		err = msgpack.Unmarshal(data, &meta)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		res = append(res, meta)
	}
	return res, nil
}

func (s Storage) fetchFile(fileName string) ([]byte, error) {
	body, err := os.ReadFile(filepath.Join(s.OutputDir, fileName))
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
	data, err := os.ReadFile(filepath.Join(s.OutputDir, ulid+".reproduced.header"))
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
		path := filepath.Join(s.OutputDir, ulid+".reproduced.header")
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
		path := filepath.Join(s.OutputDir, ulid+".reproduced.body")
		err := os.WriteFile(path, body, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
