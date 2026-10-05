package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUlidFromPath(t *testing.T) {
	expected := "SampleUlid"
	actual, err := getUlidFromPath("/SampleEndpoint/SampleUlid")
	if err != "" {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)

	expectedError := "invalid URL: /SampleEndpoint/SampleUlid/, should be /SampleEndpoint/[ulid]"
	_, actualError := getUlidFromPath("/SampleEndpoint/SampleUlid/")
	assert.Equal(t, expectedError, actualError)

	expectedError = "invalid URL: /SampleEndpoint/, should be /SampleEndpoint/[ulid]"
	_, actualError = getUlidFromPath("/SampleEndpoint/")
	assert.Equal(t, expectedError, actualError)

	expectedError = "invalid URL: /"
	_, actualError = getUlidFromPath("/")
	assert.Equal(t, expectedError, actualError)
}

func TestAcceptsGzip(t *testing.T) {
	cases := map[string]bool{
		"":                        false,
		"gzip":                    true,
		"gzip, deflate, br, zstd": true,
		"deflate, gzip;q=0.8":     true,
		"identity":                false,
		"gzip;q=0":                false,
		"br, gzip; q=0.000":       false,
		"x-gzip":                  false,
	}
	for header, expected := range cases {
		r, _ := http.NewRequest("GET", "/", nil)
		if header != "" {
			r.Header.Set("Accept-Encoding", header)
		}
		assert.Equal(t, expected, acceptsGzip(r), header)
	}
}
