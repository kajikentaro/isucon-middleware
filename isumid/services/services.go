package services

import (
	"net/http"

	"github.com/kajikentaro/isucon-middleware/isumid/storages"
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

type FetchBodyResponse struct {
	body   []byte
	header http.Header
}

func (s Service) FetchReqBody(ulid string) (FetchBodyResponse, error) {
	saved, err := s.storage.FetchReqBody(ulid)
	if err != nil {
		return FetchBodyResponse{}, err
	}

	s.storage.
}

func (s Service) FetchReqBody(ulid string) (storages.RecordedByteOutput, error) {
	saved, err := s.storage.Fetch(ulid)
	if err != nil {
		return storages.RecordedByteOutput{}, err
	}

	return saved, nil
}
