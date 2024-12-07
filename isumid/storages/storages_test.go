package storages

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/stretchr/testify/assert"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestMain(m *testing.M) {
	fmt.Println(OUTPUT_DIR)
	m.Run()
	os.RemoveAll(OUTPUT_DIR)
}

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
	storage, err := New(setting)
	assert.NoError(t, err)

	// test start
	err = storage.Save(saveData)
	assert.NoError(t, err)

	storage.Flush()
}

func getUlid(t *testing.T, s Storage) string {
	metaList, err := s.FetchMetaList(0, 1)
	assert.NoError(t, err)
	return metaList[0].Ulid
}

func TestFetchMeta(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	ulid := getUlid(t, storage)
	actual, err := storage.FetchMeta(ulid)
	assert.NoError(t, err)

	expected := models.Meta{
		Method:     "GET",
		Url:        "/test-url",
		ReqHeader:  map[string][]string{"Content-Type": {"application/octet-stream"}},
		StatusCode: 200,
		ResHeader:  map[string][]string{"Content-Type": {"text/plain"}},
		IsReqText:  false,
		IsResText:  true,
		Ulid:       ulid,
		ReqLength:  17,
		ResLength:  18,
	}
	assert.Exactly(t, expected, actual)
}

func TestFetchMetaList(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	actual, err := storage.FetchMetaList(0, 1)
	assert.NoError(t, err)
	// ignore ulid
	actual[0].Ulid = ""

	expected := []models.Meta{{
		Method:     "GET",
		Url:        "/test-url",
		ReqHeader:  map[string][]string{"Content-Type": {"application/octet-stream"}},
		StatusCode: 200,
		ResHeader:  map[string][]string{"Content-Type": {"text/plain"}},
		IsReqText:  false,
		IsResText:  true,
		Ulid:       "",
		ReqLength:  17,
		ResLength:  18,
	}}
	assert.Exactly(t, expected, actual)
}

func TestFetchReqBody(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	ulid := getUlid(t, storage)
	actual, err := storage.FetchReqBody(ulid)
	assert.NoError(t, err)

	expected := []byte("Test Request Body")
	assert.Exactly(t, expected, actual)
}

func TestFetchResBody(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	ulid := getUlid(t, storage)
	actual, err := storage.FetchResBody(ulid)
	assert.NoError(t, err)

	expected := []byte("Test Response Body")
	assert.Exactly(t, expected, actual)
}

func TestSaveReproduced(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	ulid := getUlid(t, storage)
	err = storage.SaveReproduced(
		ulid,
		[]byte("Test Reproduced Body"),
		map[string][]string{"Content-Type": {"text/plain"}},
	)
	assert.NoError(t, err)
}

func TestFetchReproducedHeader(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	ulid := getUlid(t, storage)
	actual, err := storage.FetchReproducedHeader(ulid)
	assert.NoError(t, err)

	expected := map[string][]string{"Content-Type": {"text/plain"}}
	assert.Exactly(t, expected, actual)
}

func TestFetchReproducedBody(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	ulid := getUlid(t, storage)
	actual, err := storage.FetchReproducedBody(ulid)
	assert.NoError(t, err)

	expected := []byte("Test Reproduced Body")
	assert.Exactly(t, expected, actual)
}

func TestIsText(t *testing.T) {
	{
		actual := IsText(map[string][]string{
			"foo": {"bar"},
		}, nil)
		expected := false

		assert.Equal(t, expected, actual)
	}

	{
		actual := IsText(map[string][]string{
			"Content-Type": {"text/html; charset=utf-8"},
		}, nil)
		expected := true

		assert.Equal(t, expected, actual)
	}

	{
		actual := IsText(map[string][]string{
			"Content-Type": {"text/html", "charset=utf-8"},
		}, nil)
		expected := true

		assert.Equal(t, expected, actual)
	}

	{
		actual := IsText(map[string][]string{
			"Content-Type": {"video/mp4", "text/html"},
		}, nil)
		expected := true

		assert.Equal(t, expected, actual)
	}
}

func TestFetchTotalTransactions(t *testing.T) {
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	metaList, err := storage.FetchMetaList(0, 100)
	assert.NoError(t, err)

	actual, err := storage.FetchTotalTransactions()
	assert.NoError(t, err)

	expected := len(metaList)
	assert.NotEqual(t, expected, 0, "expected should be more than 0 for better testing")
	assert.Equal(t, expected, actual)
}

func TestRemove(t *testing.T) {
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage, err := New(setting)
	assert.NoError(t, err)

	metaList, err := storage.FetchMetaList(0, 100)
	assert.NoError(t, err)

	if err := storage.Remove(metaList[0].Ulid); err != nil {
		t.Fatal(err)
	}

	actual, err := storage.FetchMetaList(0, 100)
	assert.NoError(t, err)

	expected := metaList[1:]
	assert.Exactly(t, expected, actual)
}
