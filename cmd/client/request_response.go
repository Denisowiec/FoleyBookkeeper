package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func sendRequest[T any](input T, method, url, jwt string) (*http.Response, error) {
	// This function sends a request to the server and returns the response
	dat, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(dat)
	client := &http.Client{}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	if jwt != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	}

	return client.Do(req)
}

func sendEmptyRequest(method, url, jwt string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	if jwt != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	}

	return client.Do(req)
}

func processResponse[T any](resp *http.Response, dat *T) error {
	// This is a quick function that encapsulates the steps for processing server's response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, dat)
	if err != nil {
		return err
	}
	return nil
}
