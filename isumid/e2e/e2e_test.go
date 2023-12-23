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
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kajikentaro/isucon-middleware/isumid"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/services"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
	"github.com/stretchr/testify/assert"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestMain(m *testing.M) {
	fmt.Println("test dir:", OUTPUT_DIR)

	// prepare server
	rec := isumid.New(models.Setting{OutputDir: OUTPUT_DIR})
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

	w.Header().Add("sample header", "sample header")

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
	res, err := http.Get("http://localhost:8888/all")
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

	actual := []services.RecordedTransaction{}
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		t.Fatal(err)
	}
	actual[0].Ulid = ""

	expected := []services.RecordedTransaction{{
		ResBody: "",
		ReqBody: "Hello World",
		Meta: storages.Meta{
			Url: "/",
			ReqHeader: map[string][]string{
				"Accept-Encoding": {"gzip"},
				"Content-Length":  {"11"},
				"Content-Type":    {"text/plain"},
				"User-Agent":      {"Go-http-client/1.1"},
			},
			Method: "POST",
			ResHeader: map[string][]string{
				"sample header": {"sample header"},
			},
			IsReqText:  true,
			IsResText:  false,
			StatusCode: 200,
			Ulid:       "",
		},
	}}

	assert.Exactly(t, expected, actual)
}

func fetchFirstUlid() (string, error) {
	res, err := http.Get("http://localhost:8888/all")
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

	actual := []services.RecordedTransaction{}
	err = json.Unmarshal(responseBody, &actual)
	if err != nil {
		return "", err
	}
	return actual[0].Ulid, nil
}

func TestFetchResBody(t *testing.T) {
	ulid, err := fetchFirstUlid()
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://localhost:8888/res-body/" + ulid)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}

	actualBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := []byte("Hello World Response")
	assert.Exactly(t, expectedBody, actualBody)

	actualHeader := res.Header
	actualHeader.Del("Date")
	expectedHeader := http.Header{"Access-Control-Allow-Origin": []string{"*"}, "Content-Length": []string{"20"}, "Content-Type": []string{"text/plain; charset=utf-8"}}
	assert.Exactly(t, expectedHeader, actualHeader)
}

func TestFetchReqBody(t *testing.T) {
	ulid, err := fetchFirstUlid()
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://localhost:8888/req-body/" + ulid)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}

	actualBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := []byte("Hello World")
	assert.Exactly(t, expectedBody, actualBody)

	actualHeader := res.Header
	actualHeader.Del("Date")
	expectedHeader := http.Header{"Accept-Encoding": []string{"gzip"}, "Access-Control-Allow-Origin": []string{"*"}, "Content-Length": []string{"11"}, "Content-Type": []string{"text/plain"}, "User-Agent": []string{"Go-http-client/1.1"}}
	assert.Exactly(t, expectedHeader, actualHeader)
}

func TestReproduce(t *testing.T) {
	ulid, err := fetchFirstUlid()
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://localhost:8888/reproduce/" + ulid)
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
	expected := `{"IsSameResBody":true,"IsSameResHeader":true,"IsSameStatusCode":true,"ActualResHeader":{"sample header":["sample header"]},"ActualResBody":"","IsBodyText":false,"StatusCode":200}`
	actual := string(responseBody)
	if expected != actual {
		t.Fatalf("response body is not correct: expected %s, actual %s", expected, actual)
	}
}
