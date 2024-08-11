package api

import (
	"encoding/json"
	"net/http"
)

type header struct {
	Key   string
	Value string
}

var jsonHeader = header{
	"Accept",
	"application/json",
}

func makeGetRequest[T any](url string, headers []header, response *T) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}

	return nil
}
