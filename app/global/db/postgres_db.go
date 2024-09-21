package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type PostgresRawCredential struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSL      string
}

type PostgresParam struct {
	PostgresURL string
}

func (s *PostgresRawCredential) NewPostgresParam() PostgresParam {
	return PostgresParam{
		PostgresURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", s.Username, s.Password, s.Host, s.Port, s.Database, s.SSL),
	}
}

func PostgresGetEnvVariable() *PostgresRawCredential {
	return &PostgresRawCredential{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DATABASE"),
		SSL:      os.Getenv("POSTGRES_SSLMODE"),
	}
}

func NewPostgresClient(dbURL PostgresParam) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbURL.PostgresURL)
	if err != nil {
		logrus.Fatal("error while trying connect to database, err:", err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		logrus.Fatal("error while ping to database, err:", err)
	}
	return conn
}
