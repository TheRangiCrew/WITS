package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LogRecord struct {
	Time  time.Time  `json:"time"`
	Level slog.Level `json:"level"`
	Msg   string     `json:"msg"`
}

type Logger struct {
	logger  *slog.Logger
	db      *pgxpool.Pool
	Product string      `json:"product,omitempty"`
	AWIPS   string      `json:"awips,omitempty"`
	WMO     string      `json:"wmo,omitempty"`
	Text    string      `json:"text"`
	Records []LogRecord `json:"-"`
}

func New(db *pgxpool.Pool, level slog.Level) Logger {

	opts := &slog.HandlerOptions{
		Level: level,
	}

	logger := Logger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, opts)),
		db:     db,
	}

	return logger

}

func (logger *Logger) Enabled(level slog.Level) bool {
	return logger.logger.Enabled(context.TODO(), level)
}

func (logger *Logger) Debug(msg string) {
	logger.addRecord(msg, slog.LevelDebug)

	logger.logger.Debug(msg)
}

func (logger *Logger) Info(msg string) {
	logger.addRecord(msg, slog.LevelInfo)

	logger.logger.Info(msg)
}

func (logger *Logger) Warn(msg string, args ...any) {
	logger.addRecord(msg, slog.LevelWarn)

	logger.logger.Warn(msg, args...)
}

func (logger *Logger) Error(msg string, args ...any) {
	logger.addRecord(msg, slog.LevelError)

	logger.logger.Error(msg, args...)
}

func (logger *Logger) Save() error {

	logs := []db.Log{}

	for _, record := range logger.Records {
		if !logger.logger.Enabled(context.TODO(), record.Level) {
			continue
		}

		logs = append(logs, db.Log{
			Time:    &record.Time,
			Level:   record.Level.String(),
			Product: logger.Product,
			AWIPS:   logger.AWIPS,
			WMO:     logger.WMO,
			Text:    logger.Text,
			Message: record.Msg,
		})
	}
	// _, err := surrealdb.Insert[db.Log](logger.db, "log", logs)

	return nil
}

func (logger *Logger) SetProduct(id string) {
	logger.Product = id

	logger.logger = logger.logger.With("product", id)
}

func (logger *Logger) SetAWIPS(data string) {
	logger.AWIPS = data

	logger.logger = logger.logger.With("awips", data)
}

func (logger *Logger) SetWMO(data string) {
	logger.WMO = data

	logger.logger = logger.logger.With("wmo", data)
}

func (logger *Logger) SetText(data string) {
	logger.Text = data

	logger.logger = logger.logger.With("text", data)
}

func (logger *Logger) addRecord(msg string, level slog.Level) {
	logger.Records = append(logger.Records, LogRecord{
		Time:  time.Now(),
		Level: level,
		Msg:   msg,
	})
}
