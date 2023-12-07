package test_integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	routers "github.com/kajikentaro/isucon-middleware/isumid"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestMain(m *testing.M) {
	fmt.Println("test dir:", OUTPUT_DIR)

	// prepare server
	rec := routers.New(models.Setting{OutputDir: OUTPUT_DIR})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(handler)))
	srv := &http.Server{Addr: ":8888", Handler: mux}
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "failed to start server: %s", err)
			os.Exit(1)
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to shutdown server")
			os.Exit(1)
		}
	}()

	m.Run()
}

func handler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte((fmt.Sprintf("failed to read body: %s", err))))
		return
	}
	res := string(b) + " Response"
	_, err = w.Write([]byte(res))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte((fmt.Sprintf("failed to write body: %s", err))))
	}

	w.WriteHeader(200)
}

func TestRecord(t *testing.T) {
	// send request
	requestBody := "Hello World"
	res, err := http.Post("http://localhost:8888/", "text/plain", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}

	// assert response
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	expected := requestBody + " Response"
	actual := string(responseBody)
	if expected != actual {
		t.Fatalf("response body is not correct: expected %s, actual %s", expected, actual)
	}
}

func TestFetchAll(t *testing.T) {
	res, err := http.Get("http://localhost:8888/fetch-all")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	actual := []storages.RecordedDisplayableOutput{}
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		t.Fatal(err)
	}
	actual[0].Ulid = ""

	expected := []storages.RecordedDisplayableOutput{{
		Meta:      storages.Meta{IsReqText: true, IsResText: false, StatusCode: 200, Ulid: ""},
		ResBody:   "",
		ResHeader: map[string][]string{},
		ReqBody:   "Hello World",
		ReqOthers: storages.RequestOthers{
			Url: "/",
			Header: map[string][]string{
				"Accept-Encoding": {"gzip"},
				"Content-Length":  {"11"},
				"Content-Type":    {"text/plain"},
				"User-Agent":      {"Go-http-client/1.1"},
			},
			Method: "POST",
		},
	}}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("response is not expected")
	}
}

func fetchFirstUlid() (string, error) {
	res, err := http.Get("http://localhost:8888/fetch-all")
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code is not 200")
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	actual := []storages.RecordedDisplayableOutput{}
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		return "", err
	}
	return actual[0].Ulid, nil
}

func TestFetch(t *testing.T) {
	ulid, err := fetchFirstUlid()
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://localhost:8888/fetch/" + ulid)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	actual := storages.RecordedByteOutput{}
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := storages.RecordedByteOutput{
		ReqBody: []byte("Hello World"),
		ResBody: []byte("Hello World Response"),
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("response is not expected")
	}
}
