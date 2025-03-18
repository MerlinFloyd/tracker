package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"my-fullstack-app/backend/internal/api"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	// Initialize Ethereum client
	_ = api.InitEthClient()

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/api/health", api.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/api/eth/block", api.BlockNumberHandler).Methods("GET")
	r.HandleFunc("/api/eth/balance", api.GetBalanceHandler).Methods("GET")

	return r
}

func TestServerEndpoints(t *testing.T) {
	r := setupRouter()

	// Test the /api/hello endpoint
	req, err := http.NewRequest("GET", "/api/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/api/hello handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test the /api/health endpoint
	req, err = http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/api/health handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if message, ok := response["message"].(string); !ok || message != "OK" {
		t.Errorf("Expected message to be 'OK', got %v", response["message"])
	}
}
