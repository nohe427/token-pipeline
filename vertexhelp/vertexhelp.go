package vertexhelp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

var DEFAULT_MODEL = "text-bison"
var DEFAULT_LOCATION = "us-central1"
var countTokenUrl = "https://{{.LOCATION}}-aiplatform.googleapis.com/v1/projects/{{.PROJECT_ID}}/locations/{{.LOCATION}}/publishers/google/models/{{.MODEL}}:countTokens"

type CountTokenParams struct {
	LOCATION   string
	PROJECT_ID string
	MODEL      string
}

type CountTokenRequest struct {
	Instances []Prompt `json:"instances"`
}

type Prompt struct {
	Prompt string `json:"prompt"`
}

type CountTokenResponse struct {
	TotalTokens             int `json:"totalTokens"`
	TotalBillableCharacters int `json:"totalBillableCharacters"`
}

func NewCountTokenParams(location string, project_id string, model string) *CountTokenParams {
	if model == "" {
		model = DEFAULT_MODEL
	}
	if location == "" {
		location = DEFAULT_LOCATION
	}
	return &CountTokenParams{
		LOCATION:   location,
		PROJECT_ID: project_id,
		MODEL:      model,
	}
}

func formatCountTokenUrl(params *CountTokenParams) (string, error) {
	tmpl, err := template.New("countTokenUrl").Parse(countTokenUrl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, params)
	return buf.String(), nil
}

func RequestTokenCount(params *CountTokenParams, req *CountTokenRequest, token string) (int, error) {
	url, err := formatCountTokenUrl(params)
	if err != nil {
		return 0, err
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}
	jsonReqBuf := bytes.NewBuffer(jsonReq)
	r, err := http.NewRequest("POST", url, jsonReqBuf)
	if err != nil {
		return 0, err
	}
	fmtToken := fmt.Sprintf("Bearer %s", token)
	r.Header.Add("Authorization", fmtToken)
	r.Header.Add("Content-Type", "application/json")
	c := http.DefaultClient
	resp, err := c.Do(r)
	if err != nil {
		return 0, err
	}

	respJson := CountTokenResponse{}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&respJson)

	return respJson.TotalTokens, nil
}
