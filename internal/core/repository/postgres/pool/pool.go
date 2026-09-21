package core_postgres_pool

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
	OpTimeout() time.Duration
}

type ConnectionPool struct {
	*pgxpool.Pool
	timeout time.Duration
}

func NewConnectionPool(config config, ctx context.Context) (*ConnectionPool, error) {
	connectionString := os.Getenv("POSTGRES_URL")
	if connectionString == "" {
		connectionString = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Database)
	}
	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to pars config: %w", err)
	}
	timectx, cansel := context.WithTimeout(ctx, config.Timeout)
	defer cansel()
	pool, err := pgxpool.NewWithConfig(timectx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	if err := pool.Ping(timectx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &ConnectionPool{
		Pool:    pool,
		timeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}
