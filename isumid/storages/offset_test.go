package storages_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
	"github.com/stretchr/testify/assert"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestSave(t *testing.T) {
	// prepqre request
	saveData := storages.RecordedDataInput{
		Method:     "GET",
		Url:        "/test-url",
		ReqHeader:  map[string][]string{"Content-Type": {"application/octet-stream"}},
		ReqBody:    []byte("Test Request Body"),
		StatusCode: 200,
		ResHeader:  map[string][]string{"Content-Type": {"text/plain"}},
		ResBody:    []byte("Test Response Body"),
	}

	// prepare storage
	setting := models.Setting{OutputDir: OUTPUT_DIR}
	storage := storages.New(setting)

	// test start
	for i := 0; i < 20; i++ {
		err := storage.Save(saveData)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFetchMetaList(t *testing.T) {
	// prepare storage
	setting := models.Setting{OutputDir: OUTPUT_DIR}
	storage := storages.New(setting)

	actual, err := storage.FetchMetaList(0, 2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(actual))

	actual, err = storage.FetchMetaList(18, 3)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(actual))
}
