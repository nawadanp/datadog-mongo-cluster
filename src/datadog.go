package main

import (
	"bytes"
	"encoding/json"
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
	Series []metric `json:"series"`
}

type points [2]float64

func pushToDatadog(key string, series series) {
	// Build the query
	url := dataDogAPIURL + "/series?api_key=" + key
	jsonStr, _ := json.Marshal(series)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Run the query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		log.Fatalf("API Error : %d", resp.StatusCode)
	}
	log.Println("Request sent")
}
