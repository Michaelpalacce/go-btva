package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Copy copies a file from src to dst
// returns the number of bytes copied and an err if any
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	return CopyFile(source, dst)
}

// CopyFile copies a fs.File ptr to a dest
func CopyFile(src fs.File, dst string) (int64, error) {
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	return io.Copy(destination, src)
}
