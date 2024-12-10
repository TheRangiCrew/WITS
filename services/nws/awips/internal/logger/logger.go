package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type LogRecord struct {
	Time  time.Time  `json:"time"`
	Level slog.Level `json:"level"`
	Msg   string     `json:"msg"`
}

type Logger struct {
	logger  *slog.Logger
	db      *surrealdb.DB
	Product models.RecordID `json:"product,omitempty"`
	AWIPS   string          `json:"awips,omitempty"`
	WMO     string          `json:"wmo,omitempty"`
	Text    string          `json:"text"`
	Records []LogRecord     `json:"-"`
}

func New(db *surrealdb.DB, level slog.Level) Logger {

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

func (logger *Logger) Warn(msg string) {
	logger.addRecord(msg, slog.LevelWarn)

	logger.logger.Warn(msg)
}

func (logger *Logger) Error(msg string) {
	logger.addRecord(msg, slog.LevelError)

	logger.logger.Error(msg)
}

func (logger *Logger) Save() error {

	logs := []db.Log{}

	for _, record := range logger.Records {
		if !logger.logger.Enabled(context.TODO(), record.Level) {
			continue
		}

		logs = append(logs, db.Log{
			Time:    &models.CustomDateTime{Time: record.Time},
			Level:   record.Level.String(),
			Product: &logger.Product,
			AWIPS:   logger.AWIPS,
			WMO:     logger.WMO,
			Text:    logger.Text,
			Message: record.Msg,
		})
	}
	_, err := surrealdb.Insert[db.Log](logger.db, "log", logs)

	return err
}

func (logger *Logger) SetProduct(product models.RecordID) {
	logger.Product = product

	logger.logger = logger.logger.With("product", product.ID)
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
