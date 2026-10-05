package services

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kajikentaro/isucon-middleware/isumid/models"
	mock_storage "github.com/kajikentaro/isucon-middleware/isumid/services/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSearch(t *testing.T) {
	t.Run("Should return body if it is text", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeta.IsReqText = true
		mockMeta.IsResText = true

		mockStorage := mock_storage.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().SearchMetaList(
			"sample%query%",
			10,
			100,
		).Return([]models.Meta{mockMeta}, 1, nil)
		mockStorage.EXPECT().FetchReqBody("sample-ulid").Return([]byte("sample-req-body"), nil)
		mockStorage.EXPECT().FetchResBody("sample-ulid").Return([]byte("sample-res-body"), nil)

		service := New(mockStorage)
		res, err := service.Search("sample*query*", 10, 100)
		assert.NoError(t, err)

		expectedRes := &SearchResponse{
			Transactions: []models.RecordedTransaction{
				{
					Meta:    mockMeta,
					ReqBody: "sample-req-body",
					ResBody: "sample-res-body",
				},
			},
			TotalHit: 1,
		}

		assert.Exactly(t, expectedRes, res)
	})

	t.Run("Should not return body if it not text", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeta.IsReqText = false
		mockMeta.IsResText = false

		mockStorage := mock_storage.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().SearchMetaList(
			"sample%query%",
			10,
			100,
		).Return([]models.Meta{mockMeta}, 1, nil)

		service := New(mockStorage)
		res, err := service.Search("sample*query*", 10, 100)
		assert.NoError(t, err)

		expectedRes := &SearchResponse{
			Transactions: []models.RecordedTransaction{
				{
					Meta:    mockMeta,
					ReqBody: "",
					ResBody: "",
				},
			},
			TotalHit: 1,
		}

		assert.Exactly(t, expectedRes, res)
	})
}

var mockMeta = models.Meta{
	Method:     "GET",
	Url:        "/test-url",
	ReqHeader:  map[string][]string{"Content-Type": {"application/octet-stream"}},
	StatusCode: 200,
	ResHeader:  map[string][]string{"Content-Type": {"text/plain"}},
	IsReqText:  true,
	IsResText:  true,
	Ulid:       "sample-ulid",
	ReqLength:  17,
	ResLength:  18,
}

func TestExport(t *testing.T) {
	t.Run("Should return the exported file and remove it on Close", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var dst string
		mockStorage := mock_storage.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().Export(gomock.Any(), int64(10), int64(20)).DoAndReturn(
			func(path string, _, _ int64) error {
				dst = path
				return os.WriteFile(path, []byte("exported"), 0666)
			})

		service := New(mockStorage)
		file, err := service.Export(10, 20)
		assert.NoError(t, err)

		assert.Regexp(t, `^isumid-\d{8}-\d{6}\.sqlite$`, filepath.Base(file.Name()))

		data, err := io.ReadAll(file)
		assert.NoError(t, err)
		assert.Equal(t, "exported", string(data))

		assert.NoError(t, file.Close())
		assert.NoDirExists(t, filepath.Dir(dst))
	})

	t.Run("Should remove the temporary directory on error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var dst string
		mockStorage := mock_storage.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().Export(gomock.Any(), int64(10), int64(20)).DoAndReturn(
			func(path string, _, _ int64) error {
				dst = path
				return errors.New("failed")
			})

		service := New(mockStorage)
		file, err := service.Export(10, 20)
		assert.Error(t, err)
		assert.Nil(t, file)
		assert.NoDirExists(t, filepath.Dir(dst))
	})
}
