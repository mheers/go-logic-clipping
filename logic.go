package logicclipping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type LogicConnection struct {
	client           http.Client
	apiKey           string
	apiEndpoint      string
	roleArn          string
	bucketInputName  string
	bucketOutputName string
	s3               *AwsS3ClientAPI
	originEnpointIDs []string
	channelIDs       []string
}

type MultiClipRequest struct {
	AssetName string
	// StartTime in UTC
	StartTime time.Time `json:"startTime"`

	// EndTime in UTC
	EndTime           time.Time `json:"endTime"`
	BucketName        string    `json:"bucketName"`
	RoleArn           string    `json:"roleArn"`
	OriginEndpointIDs []string  `json:"originEndpointIds"`
}

func (mcr *MultiClipRequest) ToClipRequests() []ClipRequest {
	var requests []ClipRequest
	originEndpointIDs := mcr.OriginEndpointIDs
	sort.Strings(originEndpointIDs)
	for x, id := range originEndpointIDs {
		assetName := fmt.Sprintf("%s_%d", mcr.AssetName, x)
		requests = append(requests, ClipRequest{
			ID:               assetName,
			StartTime:        mcr.StartTime,
			EndTime:          mcr.EndTime,
			BucketName:       mcr.BucketName,
			RoleArn:          mcr.RoleArn,
			ManifestKey:      GetManifestKey(assetName),
			OriginEndpointID: id,
		})
	}
	return requests
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

type CreateClipResponse struct {
	Result struct {
		Harvestjob
	} `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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

func NewLogicConnection(apiKey, apiEndpoint, bucketInputName, bucketOutputName, roleArn string) (*LogicConnection, error) {
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
		roleArn:          roleArn,
		bucketInputName:  bucketInputName,
		bucketOutputName: bucketOutputName,
		client:           client,
		s3:               s3,
	}, nil
}

func (lc *LogicConnection) SetOriginEndpointIDs(ids []string) {
	lc.originEnpointIDs = ids
}

func (lc *LogicConnection) GetOriginEndpointIDs() []string {
	return lc.originEnpointIDs
}

func (lc *LogicConnection) SetChannelIDs(ids []string) {
	lc.channelIDs = ids
}

func (lc *LogicConnection) GetChannelIDs() []string {
	return lc.channelIDs
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

func (lc *LogicConnection) CreateClip(clipRequest ClipRequest) (*CreateClipResponse, error) {
	if clipRequest.BucketName == "" {
		clipRequest.BucketName = lc.bucketInputName
	}
	if clipRequest.RoleArn == "" {
		clipRequest.RoleArn = lc.roleArn
	}
	payload, err := json.Marshal(clipRequest)
	if err != nil {
		return nil, err
	}
	resp, err := lc.Post(lc.apiEndpoint+"/mediapackage/harvestjob/start", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// response from json to CreateClipResponse
	var createClipResponse CreateClipResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, &createClipResponse)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse = ErrorResponse{}
		err = json.Unmarshal(respBody, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create clip: %s", errorResponse.Error)
	}

	return &createClipResponse, nil
}

func (lc *LogicConnection) CreateMultiClip(multiClipRequest *MultiClipRequest) ([]*CreateClipResponse, error) {
	clipRequests := multiClipRequest.ToClipRequests()
	clipResponses := []*CreateClipResponse{}
	for _, cr := range clipRequests {
		clipResponse, err := lc.CreateClip(cr)
		if err != nil {
			return nil, err
		}
		clipResponses = append(clipResponses, clipResponse)
	}
	return clipResponses, nil
}

func (lc *LogicConnection) CreateMultiClipDelayedUntilEndTime(multiClipRequest *MultiClipRequest) ([]*CreateClipResponse, error) {
	// check if endtime of multi clip is in the future
	if multiClipRequest.EndTime.After(time.Now()) {
		// wait the time until the endtime
		// time.Sleep(multiClipRequest.EndTime.Sub(time.Now()))
		time.Sleep(time.Until(multiClipRequest.EndTime))

		// wait 2sec more
		time.Sleep(time.Second * 2)
	}
	// create clip requests
	return lc.CreateMultiClip(multiClipRequest)
}

func (lc *LogicConnection) GetJobs(channelIDs []string) ([]*Harvestjob, error) {
	var jobs []*Harvestjob
	for _, id := range channelIDs {
		resp, err := lc.GetJobsForChannel(id)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, resp...)
	}
	return jobs, nil
}

func (lc *LogicConnection) GetJobsForChannel(channelID string) ([]*Harvestjob, error) {
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
