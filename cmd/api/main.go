package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"bank-app/internal/config"
	"bank-app/internal/handler"
	"bank-app/internal/repository"
	"bank-app/internal/service"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize db: %v", err)
	}
	defer db.Close()

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, cfg)
	handlers := handler.NewHandlers(services, logger)

	router := mux.NewRouter()

	// Публичные маршруты
	router.HandleFunc("/api/v1/register", handlers.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/login", handlers.Login).Methods(http.MethodPost)

	// Защищенные маршруты
	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(handlers.AuthMiddleware)

	// Счета
	protected.HandleFunc("/accounts", handlers.CreateAccount).Methods(http.MethodPost)
	protected.HandleFunc("/accounts", handlers.GetAccounts).Methods(http.MethodGet)
	protected.HandleFunc("/accounts/{id}", handlers.GetAccount).Methods(http.MethodGet)

	// Карты
	protected.HandleFunc("/cards", handlers.CreateCard).Methods(http.MethodPost)
	protected.HandleFunc("/cards", handlers.GetCards).Methods(http.MethodGet)
	protected.HandleFunc("/cards/{id}", handlers.GetCard).Methods(http.MethodGet)

	// Переводы
	protected.HandleFunc("/transfers", handlers.CreateTransfer).Methods(http.MethodPost)

	// Кредиты
	protected.HandleFunc("/credits", handlers.CreateCredit).Methods(http.MethodPost)
	protected.HandleFunc("/credits/{id}/schedule", handlers.GetCreditSchedule).Methods(http.MethodGet)

	// Аналитика
	protected.HandleFunc("/analytics", handlers.GetAnalytics).Methods(http.MethodGet)
	protected.HandleFunc("/accounts/{id}/predict", handlers.PredictBalance).Methods(http.MethodGet)

	logger.Infof("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
