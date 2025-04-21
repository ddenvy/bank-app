package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	JWTSecret     string
	SMTPConfig    SMTPConfig
	CBRConfig     CBRConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type CBRConfig struct {
	BaseURL    string
	SOAPAction string
}

func Load() (*Config, error) {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/bank_app?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		SMTPConfig: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.example.com"),
			Port:     587,
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
		CBRConfig: CBRConfig{
			BaseURL:    "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
			SOAPAction: "http://web.cbr.ru/KeyRate",
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
