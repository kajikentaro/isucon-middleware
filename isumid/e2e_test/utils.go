package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/kajikentaro/isucon-middleware/isumid/services"
)

func FetchList(t *testing.T, portNum int) []services.RecordedTransaction {
	u, err := url.Parse(fmt.Sprintf("http://localhost:%d/isumid/list", portNum))
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("offset", "0")
	q.Set("length", "20")
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
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

	recordedTransactions := []services.RecordedTransaction{}
	err = json.Unmarshal(responseBody, &recordedTransactions)
	if err != nil {
		t.Fatal(err)
	}
	return recordedTransactions
}

func StartServer(srv *http.Server) {
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(os.Stderr, "failed to start server: %s", err)
		os.Exit(1)
	}
}

func StopServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to shutdown server")
		os.Exit(1)
	}
}

func SampleHandler(w http.ResponseWriter, r *http.Request) {
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
}
