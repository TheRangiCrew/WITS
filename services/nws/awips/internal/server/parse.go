package server

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/handler"
)

func ParseText(filename string, dbConfig db.DBConfig) {
	data, err := os.ReadFile(filename)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	text := string(data)

	config := ServerConfig{
		DB: dbConfig,
	}

	server, err := New(config)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	h, err := handler.New(server.DB, server.Data.UGC)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info(fmt.Sprintf("Parsing %s", filename))

	h.Handle(text, time.Now())

}
