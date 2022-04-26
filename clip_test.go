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
	// assetName := "request_4"
	assetName := "request_multi_8_1"
	clip, err := client.GetClipByAssetName(assetName)
	assert.NoError(t, err)
	assert.NotEmpty(t, clip)

	data, err := clip.GetData()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	err = ioutil.WriteFile(fmt.Sprintf("/tmp/%s.mp4", assetName), data, 0644)
	assert.NoError(t, err)
}
