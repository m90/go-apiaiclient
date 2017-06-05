# go-apiaiclient

[![Build Status](https://travis-ci.org/m90/go-apiaiclient.svg?branch=master)](https://travis-ci.org/m90/go-apiaiclient)
[![godoc](https://godoc.org/github.com/m90/go-apiaiclient?status.svg)](http://godoc.org/github.com/m90/go-apiaiclient)

> request messages from api.ai like what

## Installation

Use `go get`:

```sh
$ go get github.com/m90/go-apiaiclient
```

## Usage

Instantiate a new client using `New(token string, language string)` and call `Request(message string, sessionID string, contexts *[]map[string]interface{})`:

```go
client := apiaiclient.New("my_token", "en")
contexts := &[]map[string]interface{}{
	map[string]interface{}{
		"name": "some-name",
		"lifespan": 12,
		"parameters": map[string]interface{}{
			"some-parameter": true,
		},
	},
}
response, err := client.Request("Good morning Mr. Magpie", "session_id", contexts)
```

## Tests

Run the tests:

```sh
$ make
```

### License
MIT © [Frederik Ring](http://www.frederikring.com)
