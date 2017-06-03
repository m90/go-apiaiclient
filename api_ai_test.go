package apiaiclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		message       string
		sessionID     string
		contexts      ContextCollection
		expectedError bool
		token         string
		lang          string
	}{
		{"Hello Gophers", "foo-bar", ContextCollection{}, false, "secret", "en"},
		{"Hello Gophers", "foo-bar", ContextCollection{}, true, "secret", "de"},
	}
	for _, test := range tests {
		previousEndpoint := endpoint
		previousVersion := version
		defer func() {
			endpoint = previousEndpoint
			version = previousVersion
		}()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != fmt.Sprintf("Bearer %v", test.token) {
				t.Errorf("Expected proper Authorization header, got %v", auth)
			}
			version := r.URL.Query().Get("v")
			if version != "testversion" {
				t.Errorf("Expected proper version parameter, got %v", version)
			}
			if test.expectedError {
				http.Error(w, "Expected Failure", http.StatusInternalServerError)
			}
			w.Write([]byte("{}"))
		}))
		endpoint = ts.URL
		version = "testversion"
		client := New(test.token, test.lang)
		_, err := client.Request(test.message, test.sessionID, &test.contexts)
		if err != nil && !test.expectedError {
			t.Errorf("Expected error of %v, got %v", test.expectedError, err)
		} else if test.expectedError && err == nil {
			t.Errorf("Expected error of %v, got %v", test.expectedError, err)
		}
	}
}
