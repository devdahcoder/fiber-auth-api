package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
	"time"
)

type Database struct {
	psqlDb *sql.DB
}

var (
	once sync.Once
	db   *Database
)

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(config DbConfig) (*Database, error) {
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

	return &Database{psqlDb: db}, nil
}

func InitializeDb() *Database {

	once.Do(func() {
		config := DbConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "youruser",
			Password: "yourpassword",
			DBName:   "yourdbname",
			SSLMode:  "disable",
		}

		var err error
		db, err = NewDatabase(config)
		if err != nil {
			log.Fatalf("Could not connect to database: %v", err)
		}
	})

	return db
}

func (db *Database) Close() error {
	if db.psqlDb != nil {
		return db.psqlDb.Close()
	}
	return nil
}
