package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

var logFile *os.File

func Init() {
	var err error
	filename := fmt.Sprintf("nankion-%s.log", time.Now().Format("02-01-2006"))
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	logger := slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

// Close closes the log file. Call at application shutdown (optional).
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
