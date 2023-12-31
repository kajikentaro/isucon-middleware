package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUlidFromPath(t *testing.T) {
	expected := "SampleUlid"
	actual, err := getUlidFromPath("/isumid/SampleEndpoint/SampleUlid")
	if err != "" {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)

	expectedError := "invalid URL: /isumid/SampleEndpoint/SampleUlid/, should be /isumid/SampleEndpoint/[ulid]"
	_, actualError := getUlidFromPath("/isumid/SampleEndpoint/SampleUlid/")
	assert.Equal(t, expectedError, actualError)

	expectedError = "invalid URL: /isumid/SampleEndpoint/, should be /isumid/SampleEndpoint/[ulid]"
	_, actualError = getUlidFromPath("/isumid/SampleEndpoint/")
	assert.Equal(t, expectedError, actualError)

	expectedError = "invalid URL: /isumid/"
	_, actualError = getUlidFromPath("/isumid/")
	assert.Equal(t, expectedError, actualError)
}

func TestGetFilePathFromUrlPath(t *testing.T) {
	expected := "Sample/Path"
	actual, err := getFilePathFromUrlPath("/isumid/Sample/Path")
	if err != "" {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)

	expected = "SamplePath/"
	actual, err = getFilePathFromUrlPath("/isumid/SamplePath/")
	if err != "" {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)

	expectedError := "invalid URL: /isumid/, should be /isumid/[file path]"
	_, actualError := getFilePathFromUrlPath("/isumid/")
	assert.Equal(t, expectedError, actualError)

	expectedError = "invalid URL: /fuga"
	_, actualError = getFilePathFromUrlPath("/fuga")
	assert.Equal(t, expectedError, actualError)
}
