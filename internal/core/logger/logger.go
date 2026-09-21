package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logge struct {
	*zap.Logger //(без названия поля) встроил zap.logger в свой логер (все методы zap.logger доступны теперь моему логеру)
	file        *os.File
}

func FromContext(ctx context.Context) *Logge {
	log, ok := ctx.Value("log").(*Logge)
	if !ok {
		panic("no logger in context")
	}
	return log
}

func NewLogger(config config) (*Logge, error) {
	var errs []error
	//Получаем zap.AtomicLevel
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}
	var logFile *os.File
	// создаём отдельную папку для логов
	// logs-относительный путь до деректории
	// 0755 - права доступа к деректории
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		errs = append(errs, fmt.Errorf("mkdir log folder: %w", err))

	} else {
		// получаем время когда был вызван логер
		timestamp := time.Now().UTC().Format("2006-01-02T15.04.05.000000")
		// добавляем время создания в названия log файла чтобы понимать когда он был создан
		logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", timestamp))
		// создаём log файл
		// logFilePath - имя файла
		// os.O_CREATE - флаг означает что если данного файла нет его надо создать
		// os.O_WRONLY - означает что файл только для записи
		// 0644 - права доступа
		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logFile = nil
			errs = append(errs, fmt.Errorf("open log file: %w", err))
		}
	}
	// сoздаём конфиг для логера
	zapConfig := zap.NewProductionEncoderConfig()
	// меняем то как будет выглябеть время при логировании
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15.04.05.000000")
	// создаём подЪядра логера
	encoder := zapcore.NewConsoleEncoder(zapConfig)

	//СОБИРАЕМ ВСЕ CORES
	multicore, errsCore := NewMultiCore(encoder, lvl, logFile)
	errs = append(errs, errsCore...)
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "[bootstrap] logger init: %v\n", err)
	}

	//инициализируем логер
	logger := zap.New(
		multicore,
		zap.AddCaller(),                       // параметр логирует место из которого произведён сам лог
		zap.AddStacktrace(zapcore.ErrorLevel), //параметр говорит что стек вызовов надо показывать для уровней Error и выше
	)
	return &Logge{
		Logger: logger,
		file:   logFile,
	}, nil
}

func (l *Logge) Close() {
	if l.file == nil {
		return
	}
	if err := l.file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close application logger: %v", err)
	}

}

// переопределяем метод With для нашего логера чтобы он возвращал не *zap.Logger а ссылку на наш логер *Logge
func (l *Logge) With(field ...zap.Field) *Logge {
	return &Logge{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}
