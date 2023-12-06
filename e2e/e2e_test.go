package test_integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	routers "github.com/kajikentaro/request-record-middleware"
	"github.com/kajikentaro/request-record-middleware/models"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestMain(m *testing.M) {
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

func TestE2E(t *testing.T) {
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
			t.Fatal(err)
		}
	}()

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
