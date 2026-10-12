package models

import (
	"log/slog"
	"os"
	"path"
)

// SetLogger config the slog output into [LogFile]
// Load the log file path from the configuration file
// Remeber call the Close function for log file closure
// It return the log file with name, you can call [Name] function to get path
func SetLogger() NameCloser {
	dir := logDir()

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		slog.Error("error while creating log file dir", "err", err)
		return nil
	}

	filePath := path.Join(dir, LogFile)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		slog.Error("error while creating log file", "err", err)
	}

	h := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: logLevel,
	})

	slog.SetDefault(slog.New(h))

	return f
}

// logDir return target log file path
func logDir() string {
	home := UserHome()

	if home == NonePath {
		return "."
	}
	return path.Join(UserHome(), ".local", "share", AppName)
}
