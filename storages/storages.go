package storages

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kajikentaro/request-record-middleware/types"
	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack/v5"
)

type Storage struct {
	types.Setting
}

type SaveDataInput struct {
	ResBody    []byte
	ResHeader  map[string][]string
	ReqBody    []byte
	ReqHeader  map[string][]string
	StatusCode int
}

type SaveDataOutput struct {
	Meta
	ResBody   string
	ResHeader map[string][]string
	ReqBody   string
	ReqHeader map[string][]string
}

type Meta struct {
	IsReqText  bool
	IsResText  bool
	StatusCode int
	Ulid       string
}

func New(setting types.Setting) Storage {
	return Storage{Setting: setting}
}

func isText(header map[string][]string) bool {
	contentType, ok := header["Content-Type"]
	if !ok {
		return false
	}
	if len(contentType) >= 2 || len(contentType) <= 0 {
		return false
	}
	contentTypeText := []string{"text/plain", "text/csv", "text/html", "text/css", "text/javascript", "application/json"}
	for _, c := range contentTypeText {
		if contentType[0] == c {
			return true
		}
	}
	return false
}

func (s Storage) Save(data SaveDataInput) error {
	err := os.MkdirAll(s.OutputDir, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// generate ulid
	ulid, err := ulid.New(ulid.Timestamp(time.Now()), nil)
	if err != nil {
		return err
	}
	ulidStr := ulid.String()

	// save metadata
	{
		path := filepath.Join(s.OutputDir, ulidStr+".meta")
		if err != nil {
			return err
		}
		meta := Meta{StatusCode: data.StatusCode, IsReqText: isText(data.ReqHeader), IsResText: isText(data.ResHeader), Ulid: ulidStr}
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

	// save request header data
	{
		path := filepath.Join(s.OutputDir, ulidStr+".req.header")
		data, err := msgpack.Marshal(data.ReqHeader)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			return err
		}
	}

	// save response header data
	{
		path := filepath.Join(s.OutputDir, ulidStr+".res.header")
		data, err := msgpack.Marshal(data.ResHeader)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) fetchAllMetaData() ([]Meta, error) {
	fileList, err := os.ReadDir(s.OutputDir)
	if err != nil {
		return nil, err
	}

	res := []Meta{}
	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".meta" {
			continue
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

func (s Storage) FetchAll() ([]SaveDataOutput, error) {
	metaList, err := s.fetchAllMetaData()
	if err != nil {
		return nil, err
	}

	res := []SaveDataOutput{}
	for _, meta := range metaList {
		saveData := SaveDataOutput{Meta: meta}
		{
			data, err := os.ReadFile(filepath.Join(s.OutputDir, meta.Ulid+".req.header"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			msgpack.Unmarshal(data, &saveData.ReqHeader)
		}
		{
			data, err := os.ReadFile(filepath.Join(s.OutputDir, meta.Ulid+".res.header"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			msgpack.Unmarshal(data, &saveData.ResHeader)
		}
		if meta.IsReqText {
			data, err := os.ReadFile(filepath.Join(s.OutputDir, meta.Ulid+".req.body"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			saveData.ReqBody = string(data)
		}
		if meta.IsResText {
			data, err := os.ReadFile(filepath.Join(s.OutputDir, meta.Ulid+".res.body"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			saveData.ResBody = string(data)
		}
		res = append(res, saveData)
	}

	return res, nil
}
