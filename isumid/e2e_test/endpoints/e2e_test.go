package test_e2e_endpoints

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/kajikentaro/isucon-middleware/isumid"
	utils "github.com/kajikentaro/isucon-middleware/isumid/e2e_test"
	"github.com/kajikentaro/isucon-middleware/isumid/services"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
	"github.com/stretchr/testify/assert"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestMain(m *testing.M) {
	fmt.Println("test dir:", OUTPUT_DIR)

	// prepare server
	rec := isumid.New(&settings.Setting{OutputDir: OUTPUT_DIR, RecordOnStart: true})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(utils.SampleHandler)))
	srv := &http.Server{Addr: ":8081", Handler: mux}
	go utils.StartServer(srv)
	defer utils.StopServer(srv)

	m.Run()
}

func TestRecord(t *testing.T) {
	// send request
	requestBody := "Hello World"
	res, err := http.Post("http://localhost:8081/", "text/plain", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
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

func TestFetchList(t *testing.T) {
	actual := utils.FetchList(t, 8081)
	actual[0].Ulid = ""

	expected := []services.RecordedTransaction{{
		ResBody: "Hello World Response",
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
			IsResText:  true,
			StatusCode: 200,
			Ulid:       "",
			ReqLength:  11,
			ResLength:  20,
		},
	}}

	assert.Exactly(t, expected, actual)
}

func fetchFirstUlid(t *testing.T) string {
	return utils.FetchList(t, 8081)[0].Ulid
}

func TestFetchResBody(t *testing.T) {
	ulid := fetchFirstUlid(t)

	res, err := http.Get("http://localhost:8081/isumid/res-body/" + ulid)
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
	ulid := fetchFirstUlid(t)

	res, err := http.Get("http://localhost:8081/isumid/req-body/" + ulid)
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
	ulid := fetchFirstUlid(t)

	res, err := http.Get("http://localhost:8081/isumid/reproduce/" + ulid)
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
	expected := `{"IsSameResBody":true,"IsSameResHeader":true,"IsSameStatusCode":true,"ActualResHeader":{"sample header":["sample header"]},"ActualResBody":"Hello World Response","IsBodyText":true,"StatusCode":200,"ActualResLength":20}`
	actual := string(responseBody)
	if expected != actual {
		t.Fatalf("response body is not correct: expected %s, actual %s", expected, actual)
	}
}
