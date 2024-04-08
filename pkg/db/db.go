package db

import (
	"context"
	"fmt"
	"github.com/babafemi99/WR/internal/config"
	"github.com/babafemi99/WR/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func connectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.TODO(), cfg.DataBaseUrl)
	if err != nil {
		return nil, err
	}
	return pool, pool.Ping(context.TODO())
}

func New(cfg *config.Config) *DB {
	pool, err := connectDB(cfg)
	if err != nil {
		logger.Log.Error(fmt.Errorf("[DB]: unable to connect: %v", err.Error()).Error())
	}

	d := &DB{
		pool,
	}
	return d
}
