package logicclipping

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
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
	clip, err := client.GetClipByAssetName("request_4")
	assert.NoError(t, err)
	assert.NotEmpty(t, clip)

	data, err := clip.GetData()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	err = ioutil.WriteFile("/tmp/request_4.mp4", data, 0644)
	assert.NoError(t, err)
}
