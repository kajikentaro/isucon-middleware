package test_e2e_switch

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kajikentaro/isucon-middleware/isumid"
	utils "github.com/kajikentaro/isucon-middleware/isumid/e2e_test"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/stretchr/testify/assert"
)

var OUTPUT_DIR = filepath.Join(os.TempDir(), uuid.NewString())

func TestMain(m *testing.M) {
	fmt.Println("test dir:", OUTPUT_DIR)
	m.Run()
}

func TestRecordOnStart(t *testing.T) {
	// start server
	rec := isumid.New(&settings.Setting{
		OutputDir:     OUTPUT_DIR,
		RecordOnStart: true,
	})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(utils.SampleHandler)))
	srv := &http.Server{Addr: ":8888", Handler: mux}
	go utils.StartServer(srv)
	defer utils.StopServer(srv)

	previousExpected := len(utils.FetchList(t))

	sendSampleRequest(t)

	// fetch result
	actual := len(utils.FetchList(t))
	expected := previousExpected + 1
	assert.Equal(t, expected, actual)
}

func TestDoNotRecordOnStart(t *testing.T) {
	// start server
	rec := isumid.New(&settings.Setting{
		OutputDir:     OUTPUT_DIR,
		RecordOnStart: false,
	})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(utils.SampleHandler)))
	srv := &http.Server{Addr: ":8888", Handler: mux}
	go utils.StartServer(srv)
	defer utils.StopServer(srv)

	previousExpected := len(utils.FetchList(t))

	sendSampleRequest(t)

	// fetch result
	actual := len(utils.FetchList(t))
	expected := previousExpected
	assert.Equal(t, expected, actual)
}

func TestStartAndStopRecording(t *testing.T) {
	// start server
	rec := isumid.New(&settings.Setting{
		OutputDir:     OUTPUT_DIR,
		RecordOnStart: false,
	})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(utils.SampleHandler)))
	srv := &http.Server{Addr: ":8888", Handler: mux}
	go utils.StartServer(srv)
	defer utils.StopServer(srv)

	{
		previousExpected := len(utils.FetchList(t))

		// turn on recording
		res, err := http.Post("http://localhost:8888/isumid/start-recording", "text/plain", nil)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Fatal("status code is not 200")
		}

		sendSampleRequest(t)

		// fetch result
		actual := len(utils.FetchList(t))
		expected := previousExpected + 1
		assert.Equal(t, expected, actual)
	}

	{
		previousExpected := len(utils.FetchList(t))

		// turn off recording
		res, err := http.Post("http://localhost:8888/isumid/stop-recording", "text/plain", nil)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Fatal("status code is not 200")
		}

		sendSampleRequest(t)

		// fetch result
		actual := len(utils.FetchList(t))
		expected := previousExpected
		assert.Equal(t, expected, actual)
	}
}

func TestAutoStart(t *testing.T) {
	// start server
	rec := isumid.New(&settings.Setting{
		OutputDir:     OUTPUT_DIR,
		RecordOnStart: false,
		AutoStart: &settings.AutoSwitch{
			TriggerEndpoint: "/trigger",
			AfterSec:        1,
		},
	})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(utils.SampleHandler)))
	srv := &http.Server{Addr: ":8888", Handler: mux}
	go utils.StartServer(srv)
	defer utils.StopServer(srv)

	{
		previousExpected := len(utils.FetchList(t))

		// send superfluous request
		requestBody := "Hello World"
		res, err := http.Post("http://localhost:8888/trigger/foo", "text/plain", bytes.NewBufferString(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Fatal("status code is not 200")
		}

		// fetch result
		actual := len(utils.FetchList(t))
		expected := previousExpected
		assert.Equal(t, expected, actual)
	}

	{
		expected := len(utils.FetchList(t))

		// send request that is trigger to turn off recording after 1 sec
		sendSampleRequestTrigger(t)

		sendSampleRequest(t)

		time.Sleep(1 * time.Second)

		sendSampleRequest(t)
		expected++

		// fetch result
		actual := len(utils.FetchList(t))
		assert.Equal(t, expected, actual)
	}
}

func TestAutoStop(t *testing.T) {
	// start server
	rec := isumid.New(&settings.Setting{
		OutputDir:     OUTPUT_DIR,
		RecordOnStart: true,
		AutoStop: &settings.AutoSwitch{
			TriggerEndpoint: "/trigger",
			AfterSec:        1,
		},
	})
	mux := http.NewServeMux()
	mux.Handle("/", rec.Middleware(http.HandlerFunc(utils.SampleHandler)))
	srv := &http.Server{Addr: ":8888", Handler: mux}
	go utils.StartServer(srv)
	defer utils.StopServer(srv)

	{
		previousExpected := len(utils.FetchList(t))

		// send superfluous request
		requestBody := "Hello World"
		res, err := http.Post("http://localhost:8888/trigger/foo", "text/plain", bytes.NewBufferString(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Fatal("status code is not 200")
		}

		// fetch result
		actual := len(utils.FetchList(t))
		expected := previousExpected + 1
		assert.Equal(t, expected, actual)
	}

	{
		expected := len(utils.FetchList(t))

		// send request that is trigger to turn off recording after 1 sec
		sendSampleRequestTrigger(t)
		expected++

		sendSampleRequest(t)
		expected++

		time.Sleep(1 * time.Second)

		sendSampleRequest(t)

		// fetch result
		actual := len(utils.FetchList(t))
		assert.Equal(t, expected, actual)
	}
}

func sendSampleRequest(t *testing.T) {
	requestBody := "Hello World"
	res, err := http.Post("http://localhost:8888/", "text/plain", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}
}

func sendSampleRequestTrigger(t *testing.T) {
	requestBody := "Hello World"
	res, err := http.Post("http://localhost:8888/trigger", "text/plain", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("status code is not 200")
	}
}
