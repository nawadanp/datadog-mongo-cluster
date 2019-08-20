package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type metric struct {
	Metric string   `json:"metric"`
	Type   string   `json:"type"`
	Points []points `json:"points"`
	Tags   []string `json:"tags"`
}

type series struct {
	Series []*metric `json:"series"`
}

type points [2]float64

func pushToDatadog(key string, series series) error {
	// Build the query
	url := dataDogAPIURL + "/series?api_key=" + key
	jsonStr, _ := json.Marshal(series)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Run the query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e := fmt.Errorf("HTTP Request Error : %d", err)
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		e := fmt.Errorf("Datadog API Error : %d", resp.StatusCode)
		return e
	}
	log.Println("Request sent")
	return nil
}
