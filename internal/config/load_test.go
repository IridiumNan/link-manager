package config

import (
	"fmt"
	"maps"
	"os"
	"testing"

	"github.com/IridiumNan/link-manager/internal/models"
)

var expectedConf = &models.LinkConfig{
	DirMap: map[string]string{
		"~":       "~/.local",
		"~/utils": "~/.local/bin",
	},
	FileMap: map[string]string{
		"demo.sh":    "demo",
		"install.py": "install",
		"demo1.sh":   "demo1",
		"demo2.sh":   "demo2",
		"demo3.sh":   "demo3",
	},
}

var testConfigStr = `
{
    "dir_map": {
        "~": "~/.local",
		"~/utils": "~/.local/bin"
    },
    "file_map": {
        "demo.sh": "demo",
        "demo1.sh": "demo1",
        "demo2.sh": "demo2",
        "demo3.sh": "demo3",
		"install.py": "install"
    }
}
`

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

	err = testOnSuccessPath(successPath)
	if err != nil {
		t.Errorf("error while test self-written configuration, err: %s", err.Error())
	}
}

func testOnSuccessPath(filePath string) error {
	err := os.WriteFile(filePath, []byte(testConfigStr), 0o644)
	if err != nil {
		return fmt.Errorf("error while write the test configuration string into file, err: %s", err.Error())
	}

	c, err := LoadConfigFromDirs(testDirs)
	if err != nil {
		return fmt.Errorf("error while loading config from written file, err: %s", err.Error())
	}

	err = compareConfig(expectedConf, c)
	if err != nil {
		return err
	}

	return nil
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
