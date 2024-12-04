package storages_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
	"github.com/stretchr/testify/assert"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestSave(t *testing.T) {
	// prepqre request
	saveData := models.RecordedDataInput{
		Method:     "GET",
		Url:        "/test-url",
		ReqHeader:  map[string][]string{"Content-Type": {"application/octet-stream"}},
		ReqBody:    []byte("Test Request Body"),
		StatusCode: 200,
		ResHeader:  map[string][]string{"Content-Type": {"text/plain"}},
		ResBody:    []byte("Test Response Body"),
	}

	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := storages.New(setting)
	assert.NoError(t, err)

	// test start
	for i := 0; i < 20; i++ {
		err := storage.Save(saveData)
		assert.NoError(t, err)
	}
}

func TestFetchMetaList(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := storages.New(setting)
	assert.NoError(t, err)

	actual, err := storage.FetchMetaList(0, 2)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(actual))

	actual, err = storage.FetchMetaList(18, 3)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(actual))
}
