package models

import (
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
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

var SelectedColor = color.BgRGB(110, 110, 110)

// UserHome return path to user home dir
// if err, it will panic
// just use it without checking
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
