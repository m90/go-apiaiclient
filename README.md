# go-apiaiclient

[![Build Status](https://travis-ci.org/m90/go-apiaiclient.svg?branch=master)](https://travis-ci.org/m90/go-apiaiclient)
[![godoc](https://godoc.org/github.com/m90/go-apiaiclient?status.svg)](http://godoc.org/github.com/m90/go-apiaiclient)

> request messages from api.ai like what

## DEPRECATION NOTICE

This library is built against the v1 API of api.ai / Dialogflow, which is in the process of being deprecated. When building something new, use the [offical Go package for the V2 API](https://godoc.org/cloud.google.com/go/dialogflow/apiv2) instead. **This package will not be maintained in any way anymore.**

## Installation

Use `go get`:

```sh
$ go get github.com/m90/go-apiaiclient
```

## Usage

Instantiate a new client using `New(token string, language string)` and call `Request(message string, sessionID string, contexts json.Marshaler)`:

```go
client := apiaiclient.New("my_token", "en")
response, err := client.Request("Good morning Mr. Magpie", "session_id", &apiaiclient.ContextCollection{/*...*/})
```

## Tests

Run the tests:

```sh
$ make
```

### License
MIT © [Frederik Ring](http://www.frederikring.com)
