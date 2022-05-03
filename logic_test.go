package logicclipping

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	client := GetDemoConnection()
	jobs, err := client.GetJobs(client.GetChannelIDs())
	assert.NoError(t, err)
	assert.NotEmpty(t, jobs)

	for _, job := range jobs {
		jobStartTime := job.Starttime.Format(time.RFC3339)
		jobEndTime := job.Endtime.Format(time.RFC3339)
		jobCreatedTime := job.Createdat.Format(time.RFC3339)
		jobBucket := job.S3Destination.Bucketname
		fmt.Printf("Job: %s, Created: %s, StartTime: %s, EndTime: %s, Bucket: %s, Status: %s\n", job.ID, jobCreatedTime, jobStartTime, jobEndTime, jobBucket, job.Status)
	}
}

func TestCreateClip(t *testing.T) {
	client := GetDemoConnection()
	assetName := "request_7"
	manifestKey := GetManifestKey(assetName)
	clipRequest := ClipRequest{
		StartTime:        time.Now().UTC().Add(time.Minute * -2),
		EndTime:          time.Now().UTC().Add(time.Minute * -1),
		ID:               assetName,
		ManifestKey:      manifestKey,
		OriginEndpointID: client.originEnpointIDs[0],
	}
	clipResponse, err := client.CreateClip(clipRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, clipResponse.Result.Arn)
}

func TestCreateMultiClip(t *testing.T) {
	client := GetDemoConnection()
	assetName := "request_multi_9"
	multiClipRequest := MultiClipRequest{
		StartTime: time.Now().UTC().Add(time.Minute * -2),
		EndTime:   time.Now().UTC().Add(time.Minute * -1),
		AssetName: assetName,
		OriginEndpointIDs: []string{
			client.originEnpointIDs[0],
			client.originEnpointIDs[1],
			client.originEnpointIDs[2],
			client.originEnpointIDs[3],
			client.originEnpointIDs[4],
			client.originEnpointIDs[5],
		},
	}
	clipResponses, err := client.CreateMultiClip(multiClipRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, clipResponses)
	assert.Len(t, clipResponses, 6)
}

func TestGetDemoConnection(t *testing.T) {
	client := GetDemoConnection()
	assert.NotNil(t, client)
}

func GetDemoConnection() *LogicConnection {
	apiKey := os.Getenv("API_KEY")
	apiEndpoint := os.Getenv("API_ENDPOINT")
	bucketInputName := os.Getenv("BUCKET_INPUT_NAME")
	bucketOutputName := os.Getenv("BUCKET_OUTPUT_NAME")
	roleArn := os.Getenv("ROLE_ARN")
	originEndpointIDs := getOriginEndpointIDsFromEnv()
	channelIDs := getChannelIDsFromEnv()
	client, err := NewLogicConnection(apiKey, apiEndpoint, bucketInputName, bucketOutputName, roleArn)
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
	id04 := os.Getenv("ORIGIN_ENDPOINT_ID_04")
	id05 := os.Getenv("ORIGIN_ENDPOINT_ID_05")
	id06 := os.Getenv("ORIGIN_ENDPOINT_ID_06")
	return []string{id01, id02, id03, id04, id05, id06}
}

func getChannelIDsFromEnv() []string {
	id01 := os.Getenv("CHANNEL_ID_01")
	id02 := os.Getenv("CHANNEL_ID_02")
	id03 := os.Getenv("CHANNEL_ID_03")
	id04 := os.Getenv("CHANNEL_ID_04")
	id05 := os.Getenv("CHANNEL_ID_05")
	id06 := os.Getenv("CHANNEL_ID_06")
	return []string{id01, id02, id03, id04, id05, id06}
}
