package services

import "github.com/kajikentaro/isucon-middleware/isumid/models"

type StorageInterface interface {
	Close() error
	CreateDir() error
	FetchMeta(ulid string) (models.Meta, error)
	FetchMetaList(offset int, length int) ([]models.Meta, error)
	FetchReproducedBody(ulid string) ([]byte, error)
	FetchReproducedHeader(ulid string) (map[string][]string, error)
	FetchReqBody(ulid string) ([]byte, error)
	FetchResBody(ulid string) ([]byte, error)
	FetchTotalTransactions() (int, error)
	Remove(ulid string) error
	RemoveAll() error
	Save(data models.RecordedDataInput) error
	SaveReproduced(ulid string, body []byte, header map[string][]string) error
	SearchMetaList(urlQuery string, offset int, length int) ([]models.Meta, int, error)
}
