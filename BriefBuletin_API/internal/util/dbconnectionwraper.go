package util

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var _logger = logrus.New()

type DBConnectionWrapper struct {
	dbPool   *pgxpool.Pool
	dbConnStr string
}

// Configuration model (can be extended as needed)
type dbConfig struct {
	DBHost         string `json:"dbhost"`
	DBName         string `json:"dbname"`
	DBUserID       string `json:"dbuid"`
	Password       string `json:"dbpassword"`
	Timeout        int    `json:"timeout"`
	Port           int    `json:"dbPort"`
	ConnRetryCount int    `json:"connRetryCount"`
}

// Initialize the pool from JSON config
func NewDBConnectionWrapper(configBytes []byte) *DBConnectionWrapper {
	var config dbConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		_logger.Errorf("Invalid DB config: %v", err)
		return nil
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUserID,
		config.Password,
		config.DBHost,
		config.Port,
		config.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		_logger.Errorf("Error creating DB pool: %v", err)
		return nil
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		_logger.Errorf("Error pinging DB: %v", err)
		return nil
	}

	_logger.Info("DB pool initialized successfully")

	return &DBConnectionWrapper{
		dbPool:    pool,
		dbConnStr: connStr,
	}
}

// Get pooled connection (for use in queries)
func (dcw *DBConnectionWrapper) GetPool() *pgxpool.Pool {
	return dcw.dbPool
}

// Graceful close
func (dcw *DBConnectionWrapper) Close() {
	if dcw.dbPool != nil {
		dcw.dbPool.Close()
		_logger.Info("DB pool closed")
	}
}
