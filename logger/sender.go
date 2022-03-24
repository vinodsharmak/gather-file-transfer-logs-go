package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type sender struct {
	accessToken   string
	url           string
	machinePairID string
	machineID     string
}

type request struct {
	Instance   string `json:"instance"`
	LogContent string `json:"log_content"`
	MachineID  string `json:"machine_id,omitempty"`
}

type response struct {
	Details   string `json:"details"`
	Instance  string `json:"instance"`
	MachineID string `json:"machine_id"`
	Detail    string `json:"detail"`
	Error     string `json:"error"`
}

func newSender(accessToken, url, machinePairID string, machineID string) *sender {
	return &sender{
		accessToken:   accessToken,
		url:           url,
		machinePairID: machinePairID,
		machineID:     machineID,
	}
}

func (s *sender) send(requestBody []byte) error {
	url := fmt.Sprintf("%s/api/v1/file_transfer/machine_pair/%s/logs/", s.url, s.machinePairID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}

	req.Header.Add("Authorization-Token", s.accessToken)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %v", err)
	}
	defer resp.Body.Close()

	var data response

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("response: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		Logger.Infof("Success response: %s", data.Details)
		return nil

	case http.StatusBadRequest:
		if data.Instance != "" {
			return fmt.Errorf("error-400: %v", data.Instance)
		}
		if data.MachineID != "" {
			return fmt.Errorf("error-400: %v", data.MachineID)
		}
		return fmt.Errorf("Error: 400")

	case http.StatusNotFound:
		return fmt.Errorf("error-404: %v", data.Detail)

	case http.StatusUnauthorized:
		return fmt.Errorf("error-401: %v", data.Error)

	}
	return fmt.Errorf("Error Response: %v", resp.StatusCode)
}
