package logicclipping

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetClips(t *testing.T) {
	client := GetDemoConnection()
	clips, err := client.GetClips()
	assert.NoError(t, err)
	assert.NotEmpty(t, clips)
	for _, clip := range clips {
		fmt.Printf("Key: %s\n", *clip.Key)
	}
}

func TestGetClipByAssetName(t *testing.T) {
	client := GetDemoConnection()
	// assetName := "request_4"
	assetName := "request_multi_9_6"
	clip, err := client.GetClipByAssetName(assetName)
	assert.NoError(t, err)
	assert.NotEmpty(t, clip)

	data, err := clip.Data()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	dir := "/tmp/clips"
	err = os.MkdirAll(dir, os.ModePerm)
	require.NoError(t, err)

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.mp4", dir, assetName), data, 0644)
	assert.NoError(t, err)
}

func TestDownloadClip(t *testing.T) {
	client := GetDemoConnection()
	// assetName := "request_4"
	assetName := "request_multi_8"
	clip, err := client.GetClipByAssetName(assetName)
	assert.NoError(t, err)
	assert.NotEmpty(t, clip)

	err = clip.Download("/tmp/clips/test")
	assert.NoError(t, err)
}

func TestFileName(t *testing.T) {
	clip := &Clip{
		Object: s3.Object{
			Key: aws.String("clips/someKey"),
		},
	}
	assert.Equal(t, "someKey", clip.FileName())
}
