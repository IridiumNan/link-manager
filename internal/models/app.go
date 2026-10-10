package models

import (
	"log"
	"os"
)

const (
	AppName = "link-manager"

	ConfigName = "config"

	ConfigFormat = "json"

	NonePath = ""
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
