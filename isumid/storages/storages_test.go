package storages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
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
	saveData := RecordedDataInput{
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
	storage := New(setting)

	// test start
	err := storage.Save(saveData)
	if err != nil {
		t.Fatal(err)
	}
}

func getUlid(t *testing.T) string {
	fileList, err := os.ReadDir(OUTPUT_DIR)
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(fileList[0].Name(), ".")
	return parts[0]
}

func TestFetchMeta(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	ulid := getUlid(t)
	actual, err := storage.FetchMeta(ulid)
	if err != nil {
		t.Fatal(err)
	}

	expected := Meta{
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
	storage := New(setting)

	actual, err := storage.FetchMetaList(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	// ignore ulid
	actual[0].Ulid = ""

	expected := []Meta{{
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
	storage := New(setting)

	ulid := getUlid(t)
	actual, err := storage.FetchReqBody(ulid)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("Test Request Body")
	assert.Exactly(t, expected, actual)
}

func TestFetchResBody(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	ulid := getUlid(t)
	actual, err := storage.FetchResBody(ulid)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("Test Response Body")
	assert.Exactly(t, expected, actual)
}

func TestSaveReproduced(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	ulid := getUlid(t)
	err := storage.SaveReproduced(
		ulid,
		[]byte("Test Reproduced Body"),
		map[string][]string{"Content-Type": {"text/plain"}},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchReproducedHeader(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	ulid := getUlid(t)
	actual, err := storage.FetchReproducedHeader(ulid)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string][]string{"Content-Type": {"text/plain"}}
	assert.Exactly(t, expected, actual)
}

func TestFetchReproducedBody(t *testing.T) {
	// prepare storage
	setting := settings.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	ulid := getUlid(t)
	actual, err := storage.FetchReproducedBody(ulid)
	if err != nil {
		t.Fatal(err)
	}

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
