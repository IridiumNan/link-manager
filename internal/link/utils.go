package link

import "os"

// IsSymlink check if a file is a symlink
func IsSymlink(mode os.FileMode) bool {
	return (mode & os.ModeSymlink) != 0
}
