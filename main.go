package main

import (
	"bytes"
	"encoding/json"
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
	originEnpointIDs []string
	channelIDs       []string
}

type ClipRequest struct {
	// StartTime in UTC
	StartTime time.Time `json:"startTime"`

	// EndTime in UTC
	EndTime          time.Time `json:"endTime"`
	ID               string    `json:"id"`
	BucketName       string    `json:"bucketName"`
	ManifestKey      string    `json:"manifestKey"`
	RoleArn          string    `json:"roleArn"`
	OriginEndpointID string    `json:"originEndpointId"`
}

type GetJobRequest struct {
	ChannelID string `json:"id"`
}

type Harvestjob struct {
	Arn              string        `json:"Arn"`
	Channelid        string        `json:"ChannelId"`
	Createdat        time.Time     `json:"CreatedAt"`
	Endtime          time.Time     `json:"EndTime"`
	ID               string        `json:"Id"`
	Originendpointid string        `json:"OriginEndpointId"`
	S3Destination    S3Destination `json:"S3Destination"`
	Starttime        time.Time     `json:"StartTime"`
	Status           string        `json:"Status"`
}

type S3Destination struct {
	Bucketname  string `json:"BucketName"`
	Manifestkey string `json:"ManifestKey"`
	Rolearn     string `json:"RoleArn"`
}

type GetJobResponse struct {
	Result struct {
		Harvestjobs []*Harvestjob `json:"HarvestJobs"`
		Nexttoken   string        `json:"NextToken"`
	} `json:"result"`
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

func (lc *LogicConnection) SetOriginEndpointIDs(ids []string) {
	lc.originEnpointIDs = ids
}

func (lc *LogicConnection) SetChannelIDs(ids []string) {
	lc.channelIDs = ids
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
	resp, err := lc.Post(lc.apiEndpoint+"/mediapackage/harvestjob/start", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		payload, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create clip: %s", string(payload))
	}
	return nil
}

func (lc *LogicConnection) GetJobs(channelID string) ([]*Harvestjob, error) {
	req := GetJobRequest{
		ChannelID: channelID,
	}
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := lc.Post(lc.apiEndpoint+"/mediapackage/harvestjob/list", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		payload, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get clip: %s", string(payload))
	}
	var jobsResponse *GetJobResponse
	responsePayload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responsePayload, &jobsResponse)
	if err != nil {
		return nil, err
	}
	return jobsResponse.Result.Harvestjobs, nil
}

func (lc *LogicConnection) GetClips() ([]*s3.Object, error) {
	return lc.s3.ListObjects(lc.bucketOutputName)
}
