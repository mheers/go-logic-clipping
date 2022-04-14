package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	apiEndpoint := os.Getenv("API_ENDPOINT")
	client := NewLogicConnection(apiKey, apiEndpoint)
	s, err := client.GetJobs()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}
