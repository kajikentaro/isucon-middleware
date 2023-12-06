package storages

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/kajikentaro/request-record-middleware/models"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestMain(m *testing.M) {
	fmt.Println(OUTPUT_DIR)
	m.Run()
}

func TestSave(t *testing.T) {
	// prepare data
	header := map[string][]string{
		"Content-Type": {"text/plain"},
	}
	saveData := SaveDataInput{
		ReqBody:    []byte("Test Request Body"),
		ReqHeader:  header,
		ResBody:    []byte("Test Response Body"),
		ResHeader:  map[string][]string{},
		StatusCode: 200,
	}

	// prepare storage
	setting := models.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	// test start
	err := storage.Save(saveData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchAllMetadata(t *testing.T) {
	// prepare storage
	setting := models.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	// test start
	actual, err := storage.fetchAllMetaData()
	if err != nil {
		t.Fatal(err)
	}
	// ignore ulid
	actual[0].Ulid = ""

	expected := []Meta{{IsReqText: true, IsResText: false, StatusCode: 200, Ulid: ""}}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: \t\n%#v, actual: \t\n%#v", expected, actual)
	}
}

func TestFetchAll(t *testing.T) {
	// prepare storage
	setting := models.Setting{OutputDir: OUTPUT_DIR}
	storage := New(setting)

	// test start
	actual, err := storage.FetchAll()
	if err != nil {
		t.Fatal(err)
	}
	// ignore ulid
	actual[0].Ulid = ""

	expected := []SaveDataOutput{{Meta: Meta{IsReqText: true, IsResText: false, StatusCode: 200, Ulid: ""}, ResBody: "", ResHeader: map[string][]string{}, ReqBody: "Test Request Body", ReqHeader: map[string][]string{"Content-Type": {"text/plain"}}}}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: \t\n%#v, \nactual: \t\n%#v", expected, actual)
	}

}
