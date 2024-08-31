package configuration

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func InitialConfig() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	err := godotenv.Load()
	if err != nil {
		logrus.Warnln("Error loading .env file, use os environment")
	}

	logLevelEnv := os.Getenv("LOG_LEVEL")

	// Default to Info level if not set or invalid
	logLevel, err := logrus.ParseLevel(strings.ToLower(logLevelEnv))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

}

func InitPostgres() *sql.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?sslmode=%s", os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DATABASE"), os.Getenv("POSTGRES_SSLMODE"))
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		logrus.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return db
}
