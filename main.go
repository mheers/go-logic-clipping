package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type LogicConnection struct {
	client           http.Client
	apiKey           string
	apiEndpoint      string
	bucketName       string
	roleArn          string
	bucketOutputName string
	s3               *AwsS3ClientAPI
}

type ClipRequest struct {
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	ID               string    `json:"id"`
	BucketName       string    `json:"bucketName"`
	ManifestKey      string    `json:"manifestKey"`
	RoleArn          string    `json:"roleArn"`
	OriginEndpointID string    `json:"originEndpointId"`
}

func NewLogicConnection(apiKey, apiEndpoint, bucketName, roleArn, bucketOutputName string) (*LogicConnection, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	s3, err := NewAwsS3ClientAPI("", "")
	if err != nil {
		return nil, err
	}
	return &LogicConnection{
		apiKey:           apiKey,
		apiEndpoint:      apiEndpoint,
		bucketName:       bucketName,
		roleArn:          roleArn,
		bucketOutputName: bucketOutputName,
		client:           client,
		s3:               s3,
	}, nil
}

func (lc *LogicConnection) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-api-key", lc.apiKey)
	return lc.client.Do(req)
}

func (lc *LogicConnection) Post(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err

	}
	return lc.Do(req)
}

func (lc *LogicConnection) CreateClip(clipRequest ClipRequest) error {
	if clipRequest.BucketName == "" {
		clipRequest.BucketName = lc.bucketName
	}
	if clipRequest.RoleArn == "" {
		clipRequest.RoleArn = lc.roleArn
	}
	payload, err := json.Marshal(clipRequest)
	if err != nil {
		return err
	}
	resp, err := lc.Post(lc.apiEndpoint+"/mediapackage/start", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to create clip")
	}
	return nil
}

func (lc *LogicConnection) GetJobs() ([]*ClipRequest, error) {
	resp, err := lc.Post(lc.apiEndpoint+"/mediapackage/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// reads the body and returns an error
		payload, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("failed to get clip: %s", string(payload))
	}
	var clipRequests []*ClipRequest
	err = json.NewDecoder(resp.Body).Decode(&clipRequests)
	if err != nil {
		return nil, err
	}
	return clipRequests, nil
}

func (lc *LogicConnection) GetClips() ([]*s3.Object, error) {
	return lc.s3.ListObjects(lc.bucketOutputName)
}
