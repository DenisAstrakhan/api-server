package core_pgx_pool

import (
	"context"
	"fmt"
	"os"
	"time"

	core_postgres_pool "github.com/DenisAstrakhan/api-server/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	timeout time.Duration
}

func NewPool(config config, ctx context.Context) (*Pool, error) {
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
	return &Pool{
		Pool:    pool,
		timeout: config.Timeout,
	}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.timeout
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxRows{rows}, nil
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	return pgxRow{p.Pool.QueryRow(ctx, sql, args...)}
}

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (core_postgres_pool.CommandTag, error) {
	tag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxCommandTag{tag}, nil
}
