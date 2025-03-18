package main

import (
	"log"
	_ "my-fullstack-app/backend/docs" // Import generated swagger docs
	"my-fullstack-app/backend/internal/api"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Ethereum Balance Tracker API
// @version         1.0
// @description     A service to track and store Ethereum wallet balances
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.yourcompany.com/support
// @contact.email  support@yourcompany.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth
func main() {
	// Initialize Ethereum client
	if err := api.InitEthClient(); err != nil {
		log.Printf("Warning: Failed to initialize Ethereum client: %v", err)
	}

	r := mux.NewRouter()

	// Define your routes here
	r.HandleFunc("/api/health", api.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/api/eth/block", api.BlockNumberHandler).Methods("GET")
	r.HandleFunc("/api/eth/balance", api.GetBalanceHandler).Methods("GET")
	r.HandleFunc("/api/eth/store-balance", api.StoreBalanceHandler).Methods("GET")

	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	http.Handle("/", r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
