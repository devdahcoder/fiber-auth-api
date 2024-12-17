package database

import (
	"context"
	"database/sql"
	"fiber-auth-api/internal/logger"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sync"
	"time"
)

type PsqlDatabase struct {
	psqlDb *sql.DB
}

var (
	once sync.Once
	db   *PsqlDatabase
)

type PsqlDsnConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPsqlDsnConfig() PsqlDsnConfig {
	return PsqlDsnConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

func NewPsqlDatabase(config PsqlDsnConfig) (*PsqlDatabase, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &PsqlDatabase{psqlDb: db}, nil
}

func GetPsqlDatabase() *PsqlDatabase {

	once.Do(func() {
		dsnConfig := NewPsqlDsnConfig()

		var err error
		db, err = NewPsqlDatabase(dsnConfig)
		if err != nil {
			logger.Error("Could not connect to database: %v", err)
		}
	})

	return db
}

func (db *PsqlDatabase) GetPsqlDB() *sql.DB {
	return db.psqlDb
}

func (db *PsqlDatabase) ClosePsqlDb() error {
	if db.psqlDb != nil {
		return db.psqlDb.Close()
	}
	return nil
}
