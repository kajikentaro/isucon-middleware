package services

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kajikentaro/isucon-middleware/isumid/storages"
)

type Service struct {
	storage storages.Storage
}

type RecordedTransaction struct {
	storages.Meta
	ReqBody string
	ResBody string
}

func New(storage storages.Storage) Service {
	return Service{storage: storage}
}

func (s Service) FetchAll(offset, length int) ([]RecordedTransaction, error) {
	MetaList, err := s.storage.FetchMetaList(offset, length)
	if err != nil {
		return nil, err
	}

	result := []RecordedTransaction{}
	for _, meta := range MetaList {
		transaction := RecordedTransaction{Meta: meta}
		if meta.IsReqText {
			body, err := s.storage.FetchReqBody(meta.Ulid)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to read req body", meta.Ulid)
				continue
			}
			transaction.ReqBody = string(body)
		}
		if meta.IsResText {
			body, err := s.storage.FetchResBody(meta.Ulid)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to read res body", meta.Ulid)
				continue
			}
			transaction.ResBody = string(body)
		}

		result = append(result, transaction)
	}

	return result, nil
}

type FetchBodyResponse struct {
	Body   []byte
	Header http.Header
}

func (s Service) FetchReqBody(ulid string) (FetchBodyResponse, error) {
	body, err := s.storage.FetchReqBody(ulid)
	if err != nil {
		return FetchBodyResponse{}, err
	}

	meta, err := s.storage.FetchMeta(ulid)
	if err != nil {
		return FetchBodyResponse{}, err
	}

	res := FetchBodyResponse{
		Header: meta.ReqHeader,
		Body:   body,
	}
	return res, nil
}

func (s Service) FetchResBody(ulid string) (FetchBodyResponse, error) {
	body, err := s.storage.FetchResBody(ulid)
	if err != nil {
		return FetchBodyResponse{}, err
	}

	meta, err := s.storage.FetchMeta(ulid)
	if err != nil {
		return FetchBodyResponse{}, err
	}

	res := FetchBodyResponse{
		Header: meta.ResHeader,
		Body:   body,
	}
	return res, nil
}

func (s Service) FetchReproducedResBody(ulid string) (FetchBodyResponse, error) {
	body, err := s.storage.FetchReproducedBody(ulid)
	if err != nil {
		return FetchBodyResponse{}, err
	}

	header, err := s.storage.FetchReproducedHeader(ulid)
	if err != nil {
		return FetchBodyResponse{}, err
	}

	res := FetchBodyResponse{
		Header: header,
		Body:   body,
	}
	return res, nil
}
