package core_logger_repository_clickhouse

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ClickHouseCore struct {
	conn    clickhouse.Conn
	level   zap.AtomicLevel
	service string
	fields  map[string]any
}

func NewClickHouseCore(ctg config, level zap.AtomicLevel) (*ClickHouseCore, error) {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{ctg.host},
		Auth: clickhouse.Auth{
			Database: ctg.database,
			Username: ctg.user,
			Password: ctg.password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("open clickhouse: %w", err)
	}
	ctx, cansel := context.WithTimeout(context.Background(), ctg.timeout)
	defer cansel()
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping clickhouse: %w", err)
	}

	return &ClickHouseCore{
		conn:    conn,
		level:   level,
		service: ctg.service,
		fields:  make(map[string]any),
	}, nil
}

func (c *ClickHouseCore) Enabled(level zapcore.Level) bool { return level >= c.level.Level() }
func (c *ClickHouseCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *c
	if clone.fields == nil {
		clone.fields = make(map[string]any)
	}
	for _, f := range fields {
		clone.fields[f.Key] = fieldToString(f)
	}
	return &clone
}
func (c *ClickHouseCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}
func (c *ClickHouseCore) Sync() error { return nil }

func (c *ClickHouseCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	var userID string
	var message string

	// Проверяем c.fields (из With())
	for k, v := range c.fields {
		if k == "user_id" && v != nil {
			userID = fmt.Sprintf("%v", v)
		}
		if k == "message" && v != nil {
			message = fmt.Sprintf("%v", v)
		}
	}

	for _, f := range fields {
		// 👇 ПРОВЕРЯЕМ, ЧТО userID ЕЩЁ НЕ УСТАНОВЛЕН
		if f.Key == "user_id" && userID == "" {
			userID = fieldToString(f)
		}
		if f.Key == "message" && message == "" {
			message = fieldToString(f)
		}
	}

	if message == "" {
		message = entry.Message
	}

	query := `
    INSERT INTO app_logs (
        timestamp, level, service, user_id, message
    ) VALUES (?, ?, ?, ?, ?)
    `

	return c.conn.Exec(context.Background(), query,
		entry.Time,
		entry.Level.String(),
		c.service,
		userID,
		message,
	)
}

func fieldToString(f zapcore.Field) string {
	switch f.Type {
	case zapcore.StringType:
		return f.String
	case zapcore.StringerType:
		if f.Interface != nil {
			return fmt.Sprintf("%v", f.Interface)
		}
	case zapcore.ReflectType:
		if f.Interface != nil {
			return fmt.Sprintf("%v", f.Interface)
		}
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Uint64Type, zapcore.Uint32Type:
		return fmt.Sprintf("%d", f.Integer)
	case zapcore.Float64Type, zapcore.Float32Type:
		return fmt.Sprintf("%v", f.Interface)
	case zapcore.BoolType:
		return strconv.FormatBool(f.Integer != 0)
	case zapcore.TimeType:
		if f.Interface != nil {
			return fmt.Sprintf("%v", f.Interface)
		}
	}
	if f.Interface != nil {
		return fmt.Sprintf("%v", f.Interface)
	}
	return f.String
}
