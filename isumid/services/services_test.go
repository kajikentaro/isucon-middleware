package services

import (
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
