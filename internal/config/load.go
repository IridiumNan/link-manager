package config

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/link-manager/internal/models"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// this file load configuration from specific dir
// then return raw configuration struct

var c models.LinkConfig

var ErrInit = errors.New("init program")

var defaultSearchDirs = []string{
	".",
	configDir(),
}

// LoadConfig load configuration from default seach dirs
func LoadConfig() (*models.LinkConfig, error) {
	return LoadConfigFromDirs(defaultSearchDirs)
}

// LoadConfigFromDirs try to load configuratoin from dirs
//
// if all dirs have tried but not found, it will initialize a new one then exit
// because the default configuration is empty so program can't run as expected
//
// if program initialize for this time, it return [ErrInit], just exit 0
func LoadConfigFromDirs(dirs []string) (*models.LinkConfig, error) {
	if len(dirs) == 0 {
		return nil, errors.New("no configuration dir provided")
	}
	viper.SetOptions(viper.KeyDelimiter("::"))

	viper.SetConfigName("config")
	viper.SetConfigType("json")

	for _, dir := range dirs {
		viper.AddConfigPath(dir)
	}
	err := viper.ReadInConfig()
	if err != nil {
		slog.Warn("error while unmarshal, try to write a new default configuration", "err", err)

		successPath := initConfig(dirs)
		if successPath == models.NonePath {
			log.Fatalf("error while creating a new configuration file, err: %s", err.Error())
		}

		slog.Info("creating a configuration file", "path", successPath)

		// because the default configuration file is empty, exit directly
		// print configuration help manual then
		fmt.Println(initHelp)

		color.Green("creating an empty config file %s", successPath)
		color.Yellow("edit empty config the file %s before you run this program again", successPath)
		return nil, ErrInit
	}

	slog.Info("load configuration success, try to unmarshal")

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshal configuration, err: %s", err.Error())
	}

	if len(c.DirMap) == 0 {
		return nil, errors.New(color.RedString("empty configuration content, edit it before you run this program"))
	}

	slog.Info("unmarshal configuration success")

	return &c, nil
}

// initConfig try to initialized all dirs until success
// it return success init file path, it all failed, return [models.NonePath]
func initConfig(dirs []string) string {
	for _, dir := range dirs {
		if successPath := tryInitConfig(dir); successPath != models.NonePath {
			return successPath
		}
	}

	return models.NonePath
}

// tryInitConfig try to init configuration file on specific dir
// if success, return this configuration file path
// else [models.NonePath]
func tryInitConfig(dir string) string {
	// drop invalid dir path
	if dir == models.NonePath {
		return models.NonePath
	}
	filePath := path.Join(dir, configFullName())

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		slog.Error("error while make dir for configuration file", "err", err)
		return models.NonePath
	}

	err = os.WriteFile(filePath, initJSONConf, 0o644)
	if err != nil {
		slog.Error("error while writting configuration file to path", "file_path", filePath)
		return models.NonePath
	}

	return filePath
}
