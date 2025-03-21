package main

import (
	_ "my-fullstack-app/backend/docs" // Import generated swagger docs
	"my-fullstack-app/backend/internal/api"
	"my-fullstack-app/backend/internal/blockchain"
	"my-fullstack-app/backend/internal/logger"
	"my-fullstack-app/backend/internal/market"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	// Initialize logger
	logger.Init(logger.InfoLevel, true)

	// Initialize router
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Initialize API handlers
	// Note: This is kept for backward compatibility and health checks
	if err := api.InitEthClient(); err != nil {
		logger.Warn().Msgf("Failed to initialize legacy Ethereum client: %v", err)
	}

	// Initialize blockchain handlers
	blockchainHandler, err := blockchain.NewHandler()
	if err != nil {
		logger.Warn().Msgf("Failed to initialize blockchain handler: %v", err)
	}

	// Initialize market data handlers
	marketHandler, err := market.NewHandler()
	if err != nil {
		logger.Warn().Msgf("Failed to initialize market data handler: %v", err)
	}

	// Register API routes
	apiRouter.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")

	// Register blockchain routes if handler was initialized
	if blockchainHandler != nil {
		apiRouter.HandleFunc("/eth/block", blockchainHandler.BlockNumberHandler).Methods("GET")
		apiRouter.HandleFunc("/eth/balance", blockchainHandler.GetBalanceHandler).Methods("GET")
		apiRouter.HandleFunc("/eth/store-balance", blockchainHandler.StoreBalanceHandler).Methods("GET")
		apiRouter.HandleFunc("/eth/get-token-balances", blockchainHandler.GetTokenBalancesHandler).Methods("GET")

		// You can add the ERC20 token handlers here
		// apiRouter.HandleFunc("/eth/token-balance", blockchainHandler.GetTokenBalanceHandler).Methods("GET")
	}

	if marketHandler != nil {
		apiRouter.HandleFunc("/market/price", marketHandler.GetCurrentPriceHandler).Methods("GET")
		apiRouter.HandleFunc("/market/historical", marketHandler.GetHistoricalPriceHandler).Methods("GET")
		logger.Info().Msg("Registered market data endpoints")
	}

	// Swagger documentation endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Set up CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000", // Local React dev server
			"http://127.0.0.1:3000", // Alternative local access
			"http://frontend:3000",  // Container name if accessed within Docker network
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})

	// Apply CORS middleware to our router
	corsHandler := corsMiddleware.Handler(r)

	// Set the main router with CORS handling
	http.Handle("/", corsHandler)

	// Start the server
	logger.Info().Msg("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal().Msgf("Could not start server: %s", err)
	}
}
