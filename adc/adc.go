package adc

import (
	"encoding/json"
	"net/http"
)

var metadataurl = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"

type ADCResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func GetADCToken() (string, error) {
	r, err := http.NewRequest("POST", metadataurl, nil)
	if err != nil {
		return "", err
	}
	r.Header.Add("Metadata-Flavor", "Google")
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	adc := ADCResult{}
	if err := json.NewDecoder(resp.Body).Decode(&adc); err != nil {
		return "", err
	}
	return adc.AccessToken, nil
}
