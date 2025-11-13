package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type errorResponseType struct {
	Error string `json:"error"`
}

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

func processErrorResponse(resp *http.Response) error {
	errorBody := errorResponseType{}
	err := processResponse(resp, &errorBody)
	if err != nil {
		return err
	}
	return fmt.Errorf(errorBody.Error)
}

// Helper functions for simple GET requests

func getThing[InputType any, ReturnType any](cfg *config, urlSuffix string, requestBody InputType, output ReturnType) (ReturnType, error) {
	// Generic function since simple GET requests work similarly for all types in the REST api
	url := fmt.Sprintf("%s%s", cfg.serverAddress, urlSuffix)
	var nothing ReturnType

	resp1, err := sendRequest(requestBody, "GET", url, cfg.jwt)
	if err != nil {
		return nothing, err
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return nothing, processErrorResponse(resp1)
	}

	err = processResponse(resp1, &output)
	if err != nil {
		return nothing, err
	}

	return output, nil
}

func getThingByID[ReturnType any](cfg *config, urlSuffix, id string, output ReturnType) (ReturnType, error) {
	url := fmt.Sprintf("%s%s/%s", cfg.serverAddress, urlSuffix, id)
	var nothing ReturnType

	resp, err := sendEmptyRequest("GET", url, cfg.jwt)
	if err != nil {
		return nothing, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nothing, processErrorResponse(resp)
	}

	err = processResponse(resp, &output)
	if err != nil {
		return nothing, err
	}

	return output, nil
}
