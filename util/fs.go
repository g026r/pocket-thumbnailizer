package util

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func MakeDir(path string) error {
	el := strings.Split(path, "/")
	p := ""

	for _, e := range el {
		p = fmt.Sprintf("%s/%s", p, e)
		s, err := os.Stat(p)
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(p, os.ModePerm); err != nil { // Don't check for ErrExist, as we can't confirm it's a directory.
				return fmt.Errorf("error creating output dir: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("error opening output dir: %w", err)
		} else {
			if !s.IsDir() {
				return fmt.Errorf("output file is not a directory")
			}
		}
	}

	return nil
}
