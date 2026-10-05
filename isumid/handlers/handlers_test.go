package handlers

import (
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
