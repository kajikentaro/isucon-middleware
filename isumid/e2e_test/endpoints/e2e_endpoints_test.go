package test_e2e_endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/kajikentaro/isucon-middleware/isumid"
	utils "github.com/kajikentaro/isucon-middleware/isumid/e2e_test"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/stretchr/testify/assert"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())
var PORT_NUMBER = 8081
var HTTP_ADDRESS = fmt.Sprintf(":%d", PORT_NUMBER)
var URL_LIST = utils.GetUrlList(PORT_NUMBER)

func TestMain(m *testing.M) {
	fmt.Println("test dir:", OUTPUT_DIR)

	// prepare server
	rec := isumid.New(&settings.Setting{OutputDir: OUTPUT_DIR, RecordOnStart: true})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(utils.SampleHandler)))
	srv := &http.Server{Addr: HTTP_ADDRESS, Handler: mux}
	go utils.StartServer(srv)
	defer utils.StopServer(srv)

	m.Run()
}

func TestRecord(t *testing.T) {
	// send request
	requestBody := "Hello World"
	url := URL_LIST.UrlOrigin + "/"
	res, err := http.Post(url, "text/plain", bytes.NewBufferString(requestBody))
	assert.NoError(t, err)
	defer res.Body.Close()
	assert.Exactly(t, 200, res.StatusCode, "status code should be 200")

	// assert response
	responseBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	expected := requestBody + " Response"
	actual := string(responseBody)
	assert.Exactly(t, expected, actual)
}

func TestFetchList(t *testing.T) {
	actual := utils.FetchList(t, PORT_NUMBER)
	actual[0].Ulid = ""

	expected := []models.RecordedTransaction{{
		ResBody: "Hello World Response",
		ReqBody: "Hello World",
		Meta: models.Meta{
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

	res, err := http.Get(URL_LIST.ResBody + ulid)
	assert.NoError(t, err)
	assert.Exactly(t, 200, res.StatusCode, "status code should be 200")

	actualBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	expectedBody := []byte("Hello World Response")
	assert.Exactly(t, expectedBody, actualBody)

	actualHeader := res.Header
	actualHeader.Del("Date")
	expectedHeader := http.Header{"Access-Control-Allow-Origin": []string{"*"}, "Content-Length": []string{"20"}, "Content-Type": []string{"text/plain; charset=utf-8"}}
	assert.Exactly(t, expectedHeader, actualHeader)
}

func TestFetchReqBody(t *testing.T) {
	ulid := fetchFirstUlid(t)

	res, err := http.Get(URL_LIST.ReqBody + ulid)
	assert.NoError(t, err)
	assert.Exactly(t, 200, res.StatusCode, "status code should be 200")

	actualBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	expectedBody := []byte("Hello World")
	assert.Exactly(t, expectedBody, actualBody)

	actualHeader := res.Header
	actualHeader.Del("Date")
	expectedHeader := http.Header{"Accept-Encoding": []string{"gzip"}, "Access-Control-Allow-Origin": []string{"*"}, "Content-Length": []string{"11"}, "Content-Type": []string{"text/plain"}, "User-Agent": []string{"Go-http-client/1.1"}}
	assert.Exactly(t, expectedHeader, actualHeader)
}

func TestReproduce(t *testing.T) {
	ulid := fetchFirstUlid(t)

	res, err := http.Get(URL_LIST.Reproduce + ulid)
	assert.NoError(t, err)
	assert.Exactly(t, 200, res.StatusCode, "status code should be 200")

	responseBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	expected := models.ReproducerResponse{
		IsSameResBody:    true,
		IsSameResHeader:  true,
		IsSameStatusCode: true,
		ActualResHeader:  http.Header{"sample header": []string{"sample header"}},
		ActualResBody:    "Hello World Response",
		IsBodyText:       true,
		StatusCode:       200,
		ActualResLength:  20,
	}
	actual := models.ReproducerResponse{}
	err = json.Unmarshal(responseBody, &actual)
	assert.NoError(t, err)
	assert.Exactly(t, expected, actual)
}

func TestRemove(t *testing.T) {
	// add recorded data
	TestRecord(t)
	TestRecord(t)
	TestRecord(t)

	transactions := utils.FetchList(t, PORT_NUMBER)

	res, err := http.Get(URL_LIST.Remove + transactions[0].Ulid)
	assert.NoError(t, err)
	assert.Exactly(t, 200, res.StatusCode, "status code should be 200")

	actual := utils.FetchList(t, PORT_NUMBER)
	expected := transactions[1:]

	assert.True(t, reflect.DeepEqual(actual, expected))
}

func TestRemoveAll(t *testing.T) {
	// add recorded data
	TestRecord(t)
	TestRecord(t)
	TestRecord(t)

	res, err := http.Get(URL_LIST.RemoveAll)
	assert.NoError(t, err)
	assert.Exactly(t, 200, res.StatusCode, "status code should be 200")

	actual := utils.FetchList(t, PORT_NUMBER)
	expected := []models.RecordedTransaction{}

	assert.True(t, reflect.DeepEqual(actual, expected))
}

func TestFetchTotalTransactions(t *testing.T) {
	// add recorded data
	TestRecord(t)
	TestRecord(t)
	TestRecord(t)

	res, err := http.Get(URL_LIST.TotalTransactions)
	assert.NoError(t, err)
	assert.Exactly(t, 200, res.StatusCode, "status code should be 200")

	actual := models.FetchTotalTransactionsResponse{}
	err = json.NewDecoder(res.Body).Decode(&actual)
	assert.NoError(t, err)
	expected := models.FetchTotalTransactionsResponse{
		Count: 3,
	}

	assert.True(t, reflect.DeepEqual(actual, expected))
}
