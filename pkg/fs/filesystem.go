package fs

import (
	"fmt"
	"os"
)

func EnsureDirExists(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat directory: %v", err)
	}
	return nil
}
