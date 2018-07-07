package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type HoaxResponse struct {
	Data []struct {
		Score float64 `json:"score"`
		Text  string  `json:"text"`
		Link  string  `json:"link"`
	} `json:"data"`
	Error string `json:"error"`
}

func CalculateHoax(query string) (HoaxResponse, error) {
	response := HoaxResponse{}
	data := struct {
		Title string `json:"title"`
	}{
		Title: query,
	}

	ep := "http://localhost:5000/hoax-detector"

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return response, err
	}

	req, err := http.NewRequest("POST", ep, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, err
}
