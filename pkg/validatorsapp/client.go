package validatorsapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Client struct {
		httpClient *http.Client
		apiKey     string
	}
	ValidatorAppInfo struct {
		Network                      string    `json:"network"`
		Account                      string    `json:"account"`
		Name                         string    `json:"name"`
		WwwURL                       string    `json:"www_url"`
		Details                      string    `json:"details"`
		AvatarURL                    string    `json:"avatar_url"`
		CreatedAt                    time.Time `json:"created_at"`
		UpdatedAt                    time.Time `json:"updated_at"`
		TotalScore                   int64     `json:"total_score"`
		RootDistanceScore            int64     `json:"root_distance_score"`
		VoteDistanceScore            int64     `json:"vote_distance_score"`
		SkippedSlotScore             int64     `json:"skipped_slot_score"`
		SoftwareVersion              string    `json:"software_version"`
		SoftwareVersionScore         int64     `json:"software_version_score"`
		StakeConcentrationScore      int64     `json:"stake_concentration_score"`
		DataCenterConcentrationScore int64     `json:"data_center_concentration_score"`
		PublishedInformationScore    int64     `json:"published_information_score"`
		SecurityReportScore          int64     `json:"security_report_score"`
		ActiveStake                  int64     `json:"active_stake"`
		Commission                   int64     `json:"commission"`
		Delinquent                   bool      `json:"delinquent"`
		DataCenterKey                string    `json:"data_center_key"`
		DataCenterHost               string    `json:"data_center_host"`
		AutonomousSystemNumber       int64     `json:"autonomous_system_number"`
		VoteAccount                  string    `json:"vote_account"`
		EpochCredits                 int64     `json:"epoch_credits"`
		SkippedSlots                 int64     `json:"skipped_slots"`
		SkippedSlotPercent           string    `json:"skipped_slot_percent"`
		PingTime                     string    `json:"ping_time"`
		URL                          string    `json:"url"`
	}
)

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: time.Second * 10},
		apiKey:     apiKey,
	}
}

func (c *Client) GetValidatorInfo(network string, account string) (info ValidatorAppInfo, err error) {
	url := fmt.Sprintf("https://www.validators.app/api/v1/validators/%s/%s.json", network, account)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return info, fmt.Errorf("http.NewRequest: %s", err.Error())
	}
	req.Header.Add("Token", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return info, fmt.Errorf("http.Do: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return info, fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return info, nil
}
