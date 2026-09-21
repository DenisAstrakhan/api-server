package core_logger

import (
	"errors"
	"fmt"
	"os"

	core_logger_repository_clickhouse "github.com/DenisAstrakhan/api-server/internal/core/logger/repository/clickhouse"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MultiCore пишет сразу во все выходы
type MultiCore struct {
	console    zapcore.Core
	file       zapcore.Core
	repository zapcore.Core
}

func (m *MultiCore) Enabled(level zapcore.Level) bool {
	return true
}

func (m *MultiCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *m
	// Передаём поля в дочерние core — они сами управляют своими полями
	clone.console = clone.console.With(fields)
	clone.file = clone.file.With(fields)
	if clone.repository != nil {
		clone.repository = clone.repository.With(fields)
	}
	return &clone
}

func (m *MultiCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if m.Enabled(entry.Level) {
		return ce.AddCore(entry, m)
	}
	return ce
}

func (m *MultiCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Поля уже накоплены в дочерних core через With(), просто передаём fields от текущего вызова
	if err := m.console.Write(entry, fields); err != nil {
		fmt.Printf("Console write error: %v\n", err)
	}
	if err := m.file.Write(entry, fields); err != nil {
		fmt.Printf("File write error: %v\n", err)
	}
	if m.repository != nil {
		if err := m.repository.Write(entry, fields); err != nil {
			fmt.Printf("ClickHouse write error: %v\n", err)
		}
	}
	return nil
}

func (m *MultiCore) Sync() error {
	if err := m.console.Sync(); err != nil {
		return err
	}
	if err := m.file.Sync(); err != nil {
		return err
	}
	if m.repository != nil {
		return m.repository.Sync()
	}
	return nil
}

func NewMultiCore(encoder zapcore.Encoder, lvl zap.AtomicLevel, logFile *os.File) (*MultiCore, []error) {
	var repository zapcore.Core
	var file zapcore.Core
	var errs []error
	switch os.Getenv("LOGGER_REPOSITORY_TYPE") {
	case "clickhouse":
		ctg, err := core_logger_repository_clickhouse.NewConfig()
		if err != nil {
			errs = append(errs, fmt.Errorf("clickhouse config: %w", err))
			repository = zap.NewNop().Core()
			break
		}
		repository, err = core_logger_repository_clickhouse.NewClickHouseCore(ctg, lvl)
		if err != nil {
			errs = append(errs, fmt.Errorf("clickhouse core: %w", err))
			repository = zap.NewNop().Core()
		}
	default:
		repository = zap.NewNop().Core()
	}
	if logFile == nil {
		errs = append(errs, errors.New("log file is nil"))
		file = zap.NewNop().Core()
	} else {
		file = zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl)
	}

	multi := &MultiCore{
		console:    zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl),
		file:       file,
		repository: repository,
	}
	return multi, errs
}
