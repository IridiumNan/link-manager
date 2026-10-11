package models

import (
	"log"
	"log/slog"
	"os"
	"path"
)

const (
	AppName = "link-manager"

	ConfigName = "config"

	ConfigFormat = "json"

	NonePath = ""

	LogFile = "log.txt"

	logLevel = slog.LevelDebug
)

var userHome string = NonePath

// UserHome return path to user home dir
// if err, it will panic
func UserHome() string {
	if userHome != NonePath {
		return userHome
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error while getting homedir, exiting, err: %s", err.Error())
	}

	return userHome
}

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

	h := slog.NewTextHandler(f, nil)

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
