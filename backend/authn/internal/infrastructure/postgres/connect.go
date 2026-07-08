package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/faber-numeris/beholder/authn/internal/infrastructure/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	DBInstance *sqlx.DB
	Pool       AbstractPool
	mu         sync.Mutex
)

func Connect(ctx context.Context, cfg config.IDatabaseConfig) (AbstractPool, error) {
	mu.Lock()
	defer mu.Unlock()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost(),
		cfg.DBPort(),
		cfg.DBUser(),
		cfg.DBPassword(),
		cfg.DBName(),
		cfg.DBSSLMode(),
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		pool.Close()
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	DBInstance = db
	Pool = pool

	return pool, nil
}

func GetDB() *sqlx.DB {
	mu.Lock()
	defer mu.Unlock()
	return DBInstance
}

func GetPool() AbstractPool {
	mu.Lock()
	defer mu.Unlock()
	return Pool
}
