package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	client := GetDemoConnection()
	s, err := client.GetJobs(client.channelIDs[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

func TestCreateClip(t *testing.T) {
	client := GetDemoConnection()
	assetName := "test_asset"
	manifestKey := GetManifestKey(assetName)
	clipRequest := ClipRequest{
		StartTime:        time.Now().UTC().Add(time.Minute * -10),
		EndTime:          time.Now().UTC().Add(time.Minute * -5),
		ID:               assetName,
		ManifestKey:      manifestKey,
		OriginEndpointID: client.originEnpointIDs[0],
	}
	err := client.CreateClip(clipRequest)
	assert.NoError(t, err)
}

func TestGetClips(t *testing.T) {
	client := GetDemoConnection()
	s, err := client.GetClips()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

func TestGetDemoConnection(t *testing.T) {
	client := GetDemoConnection()
	assert.NotNil(t, client)
}

func GetDemoConnection() *LogicConnection {
	apiKey := os.Getenv("API_KEY")
	apiEndpoint := os.Getenv("API_ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")
	roleArn := os.Getenv("ROLE_ARN")
	bucketOutputName := os.Getenv("BUCKET_OUTPUT_NAME")
	originEndpointIDs := getOriginEndpointIDsFromEnv()
	channelIDs := getChannelIDsFromEnv()
	client, err := NewLogicConnection(apiKey, apiEndpoint, bucketName, roleArn, bucketOutputName)
	if err != nil {
		panic(err)
	}
	client.SetOriginEndpointIDs(originEndpointIDs)
	client.SetChannelIDs(channelIDs)
	return client
}

func getOriginEndpointIDsFromEnv() []string {
	id01 := os.Getenv("ORIGIN_ENDPOINT_ID_01")
	id02 := os.Getenv("ORIGIN_ENDPOINT_ID_02")
	id03 := os.Getenv("ORIGIN_ENDPOINT_ID_03")
	return []string{id01, id02, id03}
}

func getChannelIDsFromEnv() []string {
	id01 := os.Getenv("CHANNEL_ID_01")
	id02 := os.Getenv("CHANNEL_ID_02")
	id03 := os.Getenv("CHANNEL_ID_03")
	return []string{id01, id02, id03}
}
