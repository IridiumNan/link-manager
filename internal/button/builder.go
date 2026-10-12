package button

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/IridiumNan/link-manager/internal/models"
)

type ButtonBuilder struct {
	srcDir string

	dstDir string

	nameMap map[string]string
}

// expandHomeDir replace the first ~ with user home dir
func expandHomeDir(dir string) string {
	if dir == "~" || strings.HasPrefix(dir, "~/") {
		home := models.UserHome()

		return strings.Replace(dir, "~", home, 1)
	}

	return dir
}

// NewBuilder create a new [LinkButton] builder
// if you want to set the name map, use [ButtonBuilder.SetNameMap] before [ButtonBuilder.Build]
func NewBuilder(srcDir string, dstDir string) *ButtonBuilder {
	// expand homedir ~ with $HOME

	return &ButtonBuilder{
		srcDir:  expandHomeDir(srcDir),
		dstDir:  expandHomeDir(dstDir),
		nameMap: map[string]string{},
	}
}

// SetNameMap set the name map which remap the original file nam to another
func (b *ButtonBuilder) SetNameMap(nameMap map[string]string) {
	b.nameMap = nameMap
}

// scan build a link button based on the fileName
func (b *ButtonBuilder) scan(fileName models.Name) *LinkButton {
	name := fileName.Name()

	srcBase := path.Base(name)
	src := path.Join(b.srcDir, srcBase)
	var dst string

	// support dir/baseName mapping
	if mapName, ok := b.nameMap[name]; ok {
		dst = path.Join(b.dstDir, mapName)
	} else {

		dstBase := path.Base(name)
		dst = path.Join(b.dstDir, dstBase)
	}

	button := NewLinkButton(src, dst)

	// NOTE: use Lstat instead of Stat
	// this function will ont follow sym link
	// ensure link can be identified
	f, err := os.Lstat(dst)
	if os.IsNotExist(err) {
		button.SetStatus(models.StatusOff)
	} else if IsSymlink(f) {
		button.SetStatus(models.StatusOn)
	}

	// if exist and isn't symlink, stay the default status [models.StatusNone]
	return button
}

// Build build all buttons based on srcDir and dstDir (nameMap)
// build by sequence
func (b *ButtonBuilder) Build() ([]models.NameButtonOperator, error) {
	entries, err := os.ReadDir(b.srcDir)

	if err != nil && len(entries) == 0 {
		// this means all read failed, just panic and report error
		slog.Error("error while reading dir", "dir", b.srcDir, "err", err)

		return nil, fmt.Errorf("failed to read any thing, %s", err.Error())
	}

	if len(entries) == 0 {
		slog.Warn("dir is empty, quit...", "dir", b.srcDir)

		fmt.Printf("empty src dir: %s, quit program now\n", b.srcDir)
		os.Exit(0)
	}

	r := make([]models.NameButtonOperator, 0, len(entries))
	// return all available links
	for _, entry := range entries {
		r = append(r, b.scan(entry))
	}

	return r, nil
}
