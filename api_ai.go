package apiaiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var endpoint = "https://api.api.ai/v1/query"
var version = "20150910"

// New returns a new client
func New(token string) Requester {
	instance := &client{
		token:      token,
		httpClient: &http.Client{},
	}
	return instance
}

// Requester is the interface describing a client instance
type Requester interface {
	Request(string, string, *ContextCollection) (*Response, error)
}

type client struct {
	token      string
	httpClient *http.Client
}

type payloadGetter interface {
	GetPayload() interface{}
}

func (b *client) makeAPIRequest(payload RequestPayload) ([]byte, error) {

	requestBody, requestBodyErr := json.Marshal(payload)
	if requestBodyErr != nil {
		return nil, requestBodyErr
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(requestBody))

	q := req.URL.Query()
	q.Set("v", version)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", b.token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	call, callErr := b.httpClient.Do(req)
	if callErr != nil {
		return nil, callErr
	}

	return ioutil.ReadAll(call.Body)
}

// Request calls the backend using the given message and contexts
func (b *client) Request(message, sessionID string, contexts *ContextCollection) (*Response, error) {
	payload := RequestPayload{
		Query:     message,
		Contexts:  *contexts,
		SessionID: sessionID,
	}

	responseBody, responseErr := b.makeAPIRequest(payload)
	if responseErr != nil {
		return nil, responseErr
	}
	response := &Response{}
	unmarshalErr := json.Unmarshal(responseBody, response)

	return response, unmarshalErr
}
