package services

import (
	"github.com/kajikentaro/request-record-middleware/storages"
)

type Service struct {
	storage storages.Storage
}

func New(storage storages.Storage) Service {
	return Service{storage: storage}
}

func (s Service) FetchAll() ([]storages.RecordedDisplayableOutput, error) {
	saved, err := s.storage.FetchAll()
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (s Service) Fetch(ulid string) (storages.RecordedByteOutput, error) {
	saved, err := s.storage.Fetch(ulid)
	if err != nil {
		return storages.RecordedByteOutput{}, err
	}

	return saved, nil
}
