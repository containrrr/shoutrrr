package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func FindGitRootFromCwd() (string, error) {
	if cwd, err := os.Getwd(); err != nil {
		return "", err
	} else {
		return FindGitRoot(cwd)
	}
}

func FindGitRoot(initialPath string) (string, error) {
	path := strings.TrimSuffix(initialPath, "/")

	for {

		info, err := os.Stat(path + "/.git")
		if err == nil {
			if info.IsDir() {
				return path, nil
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}

		if li := strings.LastIndex(path, "/"); li > 0 {
			path = path[:li]
		} else {
			return "", fmt.Errorf("no .git folder in path ancestry tree")
		}
	}

}
