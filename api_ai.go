package apiaiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	endpoint = "https://api.api.ai/v1/query"
	version  = "20150910"
)

// New returns a new client
func New(token, lang string) Requester {
	return &client{token, lang}
}

// Requester is the interface describing a client instance
type Requester interface {
	Request(message, sessionID string, contexts json.Marshaler) (*Response, error)
}

type client struct {
	token string
	lang  string
}

func (c *client) makeAPIRequest(payload RequestPayload) (io.ReadCloser, error) {
	requestBody, requestBodyErr := json.Marshal(payload)
	if requestBodyErr != nil {
		return nil, requestBodyErr
	}

	req, reqErr := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(requestBody),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	q := req.URL.Query()
	q.Set("v", version)

	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"api.ai returned %v %v",
			res.StatusCode,
			http.StatusText(res.StatusCode),
		)
	}

	return res.Body, nil
}

// Request calls api.ai using the given message and contexts
func (c *client) Request(message, sessionID string, contexts json.Marshaler) (*Response, error) {
	reader, err := c.makeAPIRequest(RequestPayload{
		Query:     message,
		Contexts:  contexts,
		SessionID: sessionID,
		Lang:      c.lang,
	})
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	response := &Response{}
	decodeErr := json.NewDecoder(reader).Decode(&response)
	return response, decodeErr
}
