package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("could not parse response body: %v", err)
	}

	expected := "OK"
	if response.Message != expected {
		t.Errorf("handler returned unexpected message: got %v want %v", response.Message, expected)
	}
}
