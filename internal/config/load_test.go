package config

import (
	"fmt"
	"maps"
	"os"
	"testing"

	"github.com/IridiumNan/link-manager/internal/models"
)

var expectedConf = &models.LinkConfig{
	DirMap:  map[string]string{},
	FileMap: map[string]string{},
}

var testDirs = []string{
	"config_test",
	".",
}

func TestLoadConfig(t *testing.T) {
	// NOTE: As this function will call Exit function in process if no configuration file exist
	// so we should test internal function instead
	//
	// config, err := LoadConfig()
	// if err != nil {
	// 	t.Errorf("error while init configuration: %s", err.Error())
	// }

	successPath := initConfig(testDirs)

	if successPath == models.NonePath {
		t.Errorf("error while creating a new configuration file, no path success")
	}

	defer func() {
		err := os.Remove(successPath)
		if err != nil {
			t.Logf("error when removing test config file, remove it manually -> %s", successPath)
		}
	}()

	_, err := LoadConfigFromDirs(testDirs)
	if err == ErrInit {
		return
	}
}

func compareConfig(expected *models.LinkConfig, got *models.LinkConfig) error {
	if !maps.Equal(expected.DirMap, got.DirMap) {
		return fmt.Errorf("dir map compare error, expected: %v, got: %v", expected, got)
	}
	if !maps.Equal(expected.FileMap, got.FileMap) {
		return fmt.Errorf("file map compare error, expected: %v, got: %v", expected, got)
	}

	return nil
}
