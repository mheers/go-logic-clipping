package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	client := GetDemoConnection()
	s, err := client.GetJobs()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

func TestCreateClip(t *testing.T) {
	client := GetDemoConnection()
	assetName := "test_asset"
	manifestKey := GetManifestKey(assetName)
	clipRequest := ClipRequest{
		StartTime:        time.Now(),
		EndTime:          time.Now().Add(time.Minute * 5),
		ID:               assetName,
		ManifestKey:      manifestKey,
		OriginEndpointID: "test",
	}
	err := client.CreateClip(clipRequest)
	assert.NoError(t, err)
}

func GetDemoConnection() *LogicConnection {
	apiKey := os.Getenv("API_KEY")
	apiEndpoint := os.Getenv("API_ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")
	roleArn := os.Getenv("ROLE_ARN")
	client := NewLogicConnection(apiKey, apiEndpoint, bucketName, roleArn)
	return client
}
