package config

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type HygraphClient struct {
	endpoint string
}



func CreateHygraphClient(endpint string) *HygraphClient {
	return &HygraphClient{
		endpoint: endpint,
	}
}

func (h *HygraphClient) SendRequest(req string) []byte {
	data := map[string]string{
		"query": req,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("failed to parse request as json")
		log.Println(err.Error())
		return nil
	}

	jsonBodyIo := strings.NewReader(string(jsonData))

	request, err := http.NewRequest("POST", h.endpoint, jsonBodyIo)

	if err != nil {
		log.Println("failed to send request to URL")
		log.Println(err.Error())
		return nil
	}

	request.Header.Set("content-type", "application/json")
	client := &http.Client{Timeout: time.Second * 60}
	response, err := client.Do(request)

	if err != nil {
		log.Println("failed to send request to URL")
		log.Println(err.Error())
		return nil
	}

	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("failed to read response body")
		log.Println(err.Error())
		return nil
	}

	return responseData
}
