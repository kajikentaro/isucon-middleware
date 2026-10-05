package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	return &SearchResponse{
		Transactions: s.withTextBodies(metaList),
		TotalHit:     totalHit,
	}, nil
}

func (s Service) fetchList(offset, length int) (*SearchResponse, error) {
	metaList, err := s.storage.FetchMetaList(offset, length)
	if err != nil {
		return nil, err
	}

	totalHit, err := s.storage.FetchTotalTransactions()
	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Transactions: s.withTextBodies(metaList),
		TotalHit:     totalHit,
	}, nil
}

// withTextBodies attaches the bodies which are text. A transaction whose body cannot be read is skipped.
func (s Service) withTextBodies(metaList []models.Meta) []models.RecordedTransaction {
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
	return transactions
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

// ExportedFile is a temporary SQLite file created by Export. Close removes it.
type ExportedFile struct {
	*os.File
	dir string
}

func (f *ExportedFile) Close() error {
	return errors.Join(f.File.Close(), os.RemoveAll(f.dir))
}

// Export writes the recorded data to a temporary SQLite file for analysis, and returns it opened.
// The caller must Close it to remove the file.
func (s Service) Export(maxMetaBytes, maxBodyBytes int64) (*ExportedFile, error) {
	dir, err := os.MkdirTemp("", "isumid-export-")
	if err != nil {
		return nil, err
	}

	// the file name is also used as the name of the download, e.g. isumid-20240101-123456.sqlite
	name := fmt.Sprintf("isumid-%s.sqlite", time.Now().Format("20060102-150405"))
	dst := filepath.Join(dir, name)
	if err := s.storage.Export(dst, maxMetaBytes, maxBodyBytes); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}

	file, err := os.Open(dst)
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return &ExportedFile{File: file, dir: dir}, nil
}

func (s Service) ExportSize() (models.ExportSize, error) {
	return s.storage.ExportSize()
}
