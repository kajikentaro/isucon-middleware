package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kajikentaro/request-record-middleware/types"
)

type Service struct {
	setting types.Setting
}

func New(setting types.Setting) Service {
	return Service{setting: setting}
}

type RecordedResponse struct {
	Body   []byte
	Header []byte
}

func (s Service) FetchAll() (string, error) {
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
		if len(splited) < 2 {
			fmt.Fprintln(os.Stderr, "file name is invalid: "+file.Name())
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.setting.OutputDir, file.Name()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if _, ok := recordedAll[splited[0]]; !ok {
			recordedAll[splited[0]] = &RecordedResponse{}
		}
		if splited[1] == "body" {
			recordedAll[splited[0]].Body = data
		}
		if splited[1] == "header" {
			recordedAll[splited[0]].Header = data
		}
	}

	resList := []*RecordedResponse{}
	for _, val := range recordedAll {
		resList = append(resList, val)
	}

	res, err := json.Marshal(resList)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
