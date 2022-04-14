package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type LogicConnection struct {
	client      http.Client
	apiKey      string
	apiEndpoint string
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

func NewLogicConnection(apiKey string, apiEndpoint string) *LogicConnection {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	return &LogicConnection{
		apiKey:      apiKey,
		apiEndpoint: apiEndpoint,
		client:      client,
	}
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
		return nil, errors.New("failed to get clip")
	}
	var clipRequests []*ClipRequest
	err = json.NewDecoder(resp.Body).Decode(&clipRequests)
	if err != nil {
		return nil, err
	}
	return clipRequests, nil
}
