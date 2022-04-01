// Package postgres create database connection
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	"boostersNews/pkg/helpers"
)

var ErrNothingUpdate = errors.New("nothing updated")

type IRunner interface {
	Query(ctx context.Context, sql string, fn func(rows pgx.Rows) error, args ...interface{}) error
	QueryRow(ctx context.Context, sql string, fn func(row pgx.Row) error, args ...interface{}) error
	Exec(ctx context.Context, sql string, fn func(ct pgconn.CommandTag) error, args ...interface{}) error
	SendBatch(ctx context.Context, batch *pgx.Batch) (int64, error)
}

const maxTryCount = 3

// NewConnection returns a new database connection with the schema applied, if not already
// applied.
func NewConnection(ctx context.Context, dbURL string, logger *zap.Logger) (*pgxpool.Pool, error) {
	ll := zapadapter.NewLogger(logger.With(zap.String("module", "pgx")))

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Logger = ll
	poolConfig.ConnConfig.LogLevel = pgx.LogLevelWarn
	// poolConfig.ConnConfig.LogLevel = pgx.LogLevelDebug

	var pool *pgxpool.Pool

	for i := 0; i <= maxTryCount; i++ {
		pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
		if err != nil {
			logger.Error("Unable to connection to database on StartUp. Wait second and try again",
				zap.Int("try", i),
				zap.Stack("stack"),
				zap.String("host", dbURL),
				zap.Error(err))
			if i == maxTryCount {
				panic("very slow, exits")
			}

			time.Sleep(time.Second)
		} else {
			i = maxTryCount
			break
		}
	}

	logger.Debug("connected to DB", zap.String("dbURL", dbURL))

	return pool, nil
}

type Runner struct {
	logger   *zap.Logger
	database *pgxpool.Pool
}

func NewRunner(pool *pgxpool.Pool, logger *zap.Logger) IRunner {
	return &Runner{
		database: pool,
		logger:   logger.With(zap.String("module", "postgres")),
	}
}

// Query executes the provided query.
func (r *Runner) Query(ctx context.Context, sql string, fn func(rows pgx.Rows) error, args ...interface{}) error {
	startTime := time.Now()
	defer helpers.PostgresRequestDuration.With(prometheus.Labels{"type": "query"}).Observe(time.Since(startTime).Seconds())

	conn, err := r.database.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err = fn(rows)
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// QueryRow executes the provided query.
func (r *Runner) QueryRow(ctx context.Context, sql string, fn func(row pgx.Row) error, args ...interface{}) error {
	startTime := time.Now()
	defer helpers.PostgresRequestDuration.With(prometheus.Labels{"type": "query"}).Observe(time.Since(startTime).Seconds())

	conn, err := r.database.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	row := conn.QueryRow(ctx, sql, args...)
	err = fn(row)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) Exec(ctx context.Context, sql string, fn func(ct pgconn.CommandTag) error, args ...interface{}) error {
	startTime := time.Now()
	defer helpers.PostgresRequestDuration.With(prometheus.Labels{"type": "query"}).Observe(time.Since(startTime).Seconds())

	conn, err := r.database.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	ct, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	err = fn(ct)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return ErrNothingUpdate
	}

	return nil
}

func (r *Runner) SendBatch(ctx context.Context, batch *pgx.Batch) (int64, error) {
	startTime := time.Now()
	defer helpers.PostgresRequestDuration.With(prometheus.Labels{"type": "sendBatch"}).Observe(time.Since(startTime).Seconds())

	var (
		rowsAffected int64
		ct           pgconn.CommandTag
	)

	conn, err := r.database.Acquire(ctx)
	if err != nil {
		return rowsAffected, err
	}

	defer conn.Release()

	br := conn.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		ct, err = br.Exec()
		if err != nil {
			return rowsAffected, err
		}
		rowsAffected += ct.RowsAffected()
	}

	err = br.Close()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}
