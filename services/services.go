package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kajikentaro/request-record-middleware/types"
	"github.com/vmihailenco/msgpack/v5"
)

type Service struct {
	setting types.Setting
}

func New(setting types.Setting) Service {
	return Service{setting: setting}
}

type RecordedResponse struct {
	ResBody   []byte
	ResHeader string
	ReqBody   []byte
	ReqHeader string
}

func (s Service) FetchAll() ([]byte, error) {
	fileList, err := os.ReadDir(s.setting.OutputDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	recordedAll := map[string]*RecordedResponse{}
	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		splited := strings.Split(file.Name(), ".")
		if len(splited) < 3 {
			fmt.Fprintln(os.Stderr, "file name is invalid: "+file.Name())
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.setting.OutputDir, file.Name()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		key := splited[0]
		if _, ok := recordedAll[key]; !ok {
			recordedAll[splited[0]] = &RecordedResponse{}
		}
		if splited[1] == "res" && splited[2] == "body" {
			recordedAll[key].ResBody = data
		}
		if splited[1] == "res" && splited[2] == "header" {
			var header map[string][]string
			msgpack.Unmarshal(data, &header)
			json, err := json.Marshal(header)
			recordedAll[key].ResHeader = string(json)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
		if splited[1] == "req" && splited[2] == "body" {
			recordedAll[key].ReqBody = data
		}
		if splited[1] == "req" && splited[2] == "header" {
			var header map[string][]string
			msgpack.Unmarshal(data, &header)
			json, err := json.Marshal(header)
			recordedAll[key].ReqHeader = string(json)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
	}

	resList := []*RecordedResponse{}
	for _, val := range recordedAll {
		resList = append(resList, val)
	}

	res, err := json.Marshal(resList)
	if err != nil {
		return nil, err
	}
	return res, nil
}
