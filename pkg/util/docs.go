package util

import (
	"fmt"

	"github.com/containrrr/shoutrrr/internal/meta"
)

func DocsURL(path string) string {
	return fmt.Sprintf("https://containrrr.dev/shoutrrr/%s/%s", meta.DocsVersion, path)
}
