package button

import (
	"io/fs"
	"os"
)

type Mode interface {
	Mode() fs.FileMode
}

// IsSymlink check if a file is a symlink
func IsSymlink(f Mode) bool {
	return (f.Mode() & os.ModeSymlink) != 0
}
