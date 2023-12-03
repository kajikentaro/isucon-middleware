package services

import (
	"encoding/json"

	"github.com/kajikentaro/request-record-middleware/storages"
)

type Service struct {
	storage storages.Storage
}

func New(storage storages.Storage) Service {
	return Service{storage: storage}
}

type RecordedResponse struct {
	ResBody   []byte
	ResHeader string
	ReqBody   []byte
	ReqHeader string
}

func (s Service) FetchAll() ([]byte, error) {
	dataMap, err := s.storage.FetchAll()
	if err != nil {
		return nil, err
	}

	dataList := []storages.SaveData{}
	for _, val := range dataMap {
		dataList = append(dataList, *val)
	}

	res, err := json.Marshal(dataList)
	if err != nil {
		return nil, err
	}
	return res, nil
}
