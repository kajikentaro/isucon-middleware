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

	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/services"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
	"github.com/stretchr/testify/assert"
)

func FetchAllTransactions(t *testing.T, portNum int) []models.RecordedTransaction {
	res := SearchTransactions(t, portNum, "")
	return res.Transactions
}

func SearchTransactions(t *testing.T, portNum int, query string) services.SearchResponse {
	// wait db update
	time.Sleep(storages.AUTO_FLASH_INTERVAL + 1*time.Second)

	requestUrl := GetUrlList(portNum).Search
	u, err := url.Parse(requestUrl)
	assert.NoError(t, err)
	q := u.Query()
	q.Set("offset", "0")
	q.Set("length", "20")
	q.Set("query", query)
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	assert.NoError(t, err)
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}

	responseBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var searchResponse services.SearchResponse
	err = json.Unmarshal(responseBody, &searchResponse)
	assert.NoError(t, err)
	return searchResponse
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

type UrlList struct {
	StartRecording    string
	StopRecording     string
	IsRecording       string
	ReqBody           string
	ResBody           string
	Remove            string
	RemoveAll         string
	ReproducesResBody string
	Search            string
	Reproduce         string

	UrlPrefix string
	UrlOrigin string
}

func GetUrlList(portNum int) UrlList {
	prefix := fmt.Sprintf("http://localhost:%d/isumid", portNum)
	res := UrlList{
		StartRecording:    prefix + "/start-recording",
		StopRecording:     prefix + "/stop-recording",
		IsRecording:       prefix + "/is-recording",
		ReqBody:           prefix + "/req-body/",
		ResBody:           prefix + "/res-body/",
		Remove:            prefix + "/remove/",
		RemoveAll:         prefix + "/remove-all",
		ReproducesResBody: prefix + "/reproduces-res-body/",
		Search:            prefix + "/search",
		Reproduce:         prefix + "/reproduce/",
		UrlPrefix:         prefix,
		UrlOrigin:         fmt.Sprintf("http://localhost:%d", portNum),
	}
	return res
}
