package helper

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

func InitialConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	logLevelEnv := os.Getenv("LOG_LEVEL")

	// Default to Info level if not set or invalid
	logLevel, err := logrus.ParseLevel(strings.ToLower(logLevelEnv))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

}
