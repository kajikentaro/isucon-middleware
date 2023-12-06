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

func (s Service) FetchAll() ([]byte, error) {
	saved, err := s.storage.FetchAll()
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(saved)
	if err != nil {
		return nil, err
	}
	return res, nil
}
