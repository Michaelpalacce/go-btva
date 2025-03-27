package file

import "os"

func DeleteIfExists(filename string) error {
	if Exists(filename) {
		return os.Remove(filename)
	}

	return nil
}
