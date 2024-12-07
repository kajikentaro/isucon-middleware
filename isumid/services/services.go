package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/kajikentaro/isucon-middleware/isumid/models"
)

type Service struct {
	storage StorageInterface
}

func New(storage StorageInterface) Service {
	return Service{storage: storage}
}

type SearchResponse struct {
	Transactions []models.RecordedTransaction `json:"transactions"`
	TotalHit     int                          `json:"totalHit"`
}

func (s Service) Search(query string, offset, length int) (*SearchResponse, error) {
	if query == "" {
		// don't filter for the performance
		return s.fetchList(offset, length)
	}

	query = strings.ReplaceAll(query, "*", "%")

	metaList, totalHit, err := s.storage.SearchMetaList(query, offset, length)
	if err != nil {
		return nil, err
	}

	transactions := []models.RecordedTransaction{}
	for _, meta := range metaList {
		transaction := models.RecordedTransaction{Meta: meta}

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

		transactions = append(transactions, transaction)
	}

	return &SearchResponse{
		Transactions: transactions,
		TotalHit:     totalHit,
	}, nil
}

func (s Service) fetchList(offset, length int) (*SearchResponse, error) {
	MetaList, err := s.storage.FetchMetaList(offset, length)
	if err != nil {
		return nil, err
	}

	transactions := []models.RecordedTransaction{}
	for _, meta := range MetaList {
		transaction := models.RecordedTransaction{Meta: meta}
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

		transactions = append(transactions, transaction)
	}

	totalHit, err := s.storage.FetchTotalTransactions()
	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Transactions: transactions,
		TotalHit:     totalHit,
	}, nil
}

func (s Service) FetchReqBody(ulid string) (models.FetchBodyResponse, error) {
	body, err := s.storage.FetchReqBody(ulid)
	if err != nil {
		return models.FetchBodyResponse{}, err
	}

	meta, err := s.storage.FetchMeta(ulid)
	if err != nil {
		return models.FetchBodyResponse{}, err
	}

	res := models.FetchBodyResponse{
		Header: meta.ReqHeader,
		Body:   body,
	}
	return res, nil
}

func (s Service) FetchResBody(ulid string) (models.FetchBodyResponse, error) {
	body, err := s.storage.FetchResBody(ulid)
	if err != nil {
		return models.FetchBodyResponse{}, err
	}

	meta, err := s.storage.FetchMeta(ulid)
	if err != nil {
		return models.FetchBodyResponse{}, err
	}

	res := models.FetchBodyResponse{
		Header: meta.ResHeader,
		Body:   body,
	}
	return res, nil
}

func (s Service) FetchReproducedResBody(ulid string) (models.FetchBodyResponse, error) {
	body, err := s.storage.FetchReproducedBody(ulid)
	if err != nil {
		return models.FetchBodyResponse{}, err
	}

	header, err := s.storage.FetchReproducedHeader(ulid)
	if err != nil {
		return models.FetchBodyResponse{}, err
	}

	res := models.FetchBodyResponse{
		Header: header,
		Body:   body,
	}
	return res, nil
}

func (s Service) Remove(ulid string) error {
	return s.storage.Remove(ulid)
}

func (s Service) RemoveAll() error {
	err := s.storage.RemoveAll()
	if err != nil {
		return err
	}

	return nil
}

func (s Service) FetchTotalTransactions() (models.FetchTotalTransactionsResponse, error) {
	count, err := s.storage.FetchTotalTransactions()
	if err != nil {
		return models.FetchTotalTransactionsResponse{}, err
	}

	res := models.FetchTotalTransactionsResponse{
		Count: count,
	}
	return res, nil
}
