package config

import (
	"fmt"
	"os"
	"path"

	_ "embed"

	"github.com/IridiumNan/link-manager/internal/models"
)

//go:embed default.json
var initJSONConf []byte

// configFullName return the default configuration file full name [models.ConfigName].[models.ConfigFormat]
func configFullName() string {
	return fmt.Sprintf("%s.%s", models.ConfigName, models.ConfigFormat)
}

func xdgConfigFilePath() string {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return models.NonePath
	}
	return path.Join(confDir, models.AppName, configFullName())
}

// configDir return a dir with XDG standard
// if failed, return [models.NonePath]
func configDir() string {
	xdgDir, err := os.UserConfigDir()
	if err != nil {
		return models.NonePath
	}

	return path.Join(xdgDir, models.AppName)
}
