package storages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kajikentaro/request-record-middleware/types"
	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack/v5"
)

type Storage struct {
	types.Setting
}

type SaveData struct {
	ResBody   []byte
	ResHeader map[string][]string
	ReqBody   []byte
	ReqHeader map[string][]string
}

func New(setting types.Setting) Storage {
	return Storage{Setting: setting}
}

func (s Storage) Save(data SaveData) error {
	// generate ulid
	ulid, err := ulid.New(ulid.Timestamp(time.Now()), nil)
	if err != nil {
		return err
	}
	ulidStr := ulid.String()

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

func (s Storage) FetchAll() (map[string]*SaveData, error) {
	fileList, err := os.ReadDir(s.OutputDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	recordedAll := map[string]*SaveData{}
	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		splited := strings.Split(file.Name(), ".")
		if len(splited) < 3 {
			fmt.Fprintln(os.Stderr, "file name is invalid: "+file.Name())
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.OutputDir, file.Name()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		key := splited[0]
		if _, ok := recordedAll[key]; !ok {
			recordedAll[splited[0]] = &SaveData{}
		}
		if splited[1] == "res" && splited[2] == "body" {
			recordedAll[key].ResBody = data
		}
		if splited[1] == "res" && splited[2] == "header" {
			var header map[string][]string
			err := msgpack.Unmarshal(data, &header)
			if err != nil {
				return nil, err
			}
			recordedAll[key].ResHeader = header
		}
		if splited[1] == "req" && splited[2] == "body" {
			recordedAll[key].ReqBody = data
		}
		if splited[1] == "req" && splited[2] == "header" {
			var header map[string][]string
			err := msgpack.Unmarshal(data, &header)
			if err != nil {
				return nil, err
			}
			recordedAll[key].ResHeader = header
		}
	}
	return recordedAll, nil
}
