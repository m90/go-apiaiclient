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
func New(token, lang string) Requester {
	instance := &client{
		token:      token,
		httpClient: &http.Client{},
		lang:       lang,
	}
	return instance
}

// Requester is the interface describing a client instance
type Requester interface {
	Request(message string, sessionID string, contexts *ContextCollection) (*Response, error)
}

type client struct {
	token      string
	httpClient *http.Client
	lang       string
}

type payloadGetter interface {
	GetPayload() interface{}
}

func (c *client) makeAPIRequest(payload RequestPayload) ([]byte, error) {

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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, resErr := c.httpClient.Do(req)
	if resErr != nil {
		return nil, resErr
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"api.ai returned %v %v",
			res.StatusCode,
			http.StatusText(res.StatusCode),
		)
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// Request calls the backend using the given message and contexts
func (c *client) Request(message, sessionID string, contexts *ContextCollection) (*Response, error) {
	payload := RequestPayload{
		Query:     message,
		Contexts:  *contexts,
		SessionID: sessionID,
		Lang:      c.lang,
	}

	responseBody, responseErr := c.makeAPIRequest(payload)
	if responseErr != nil {
		return nil, responseErr
	}
	response := &Response{}
	unmarshalErr := json.Unmarshal(responseBody, response)
	return response, unmarshalErr
}
