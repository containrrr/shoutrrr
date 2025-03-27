package util

import (
	"fmt"

	"github.com/nicholas-fedor/shoutrrr/internal/meta"
)

// DocsURL returns a full documentation URL for the current version of Shoutrrr with the path appended.
// If the path contains a leading slash, it is stripped.
func DocsURL(path string) string {
	// strip leading slash if present
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	return fmt.Sprintf("https://nicholas-fedor.github.io/shoutrrr/%s/%s", meta.DocsVersion, path)
}
