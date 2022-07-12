//go:generate go run ../../../cmd/shoutrrr-gen
package matrix

import (
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/conf"
)

func (config *Config) UpdateLegacyURL(legacyURL *url.URL) *url.URL {
	updatedURL := *legacyURL
	query := legacyURL.Query()

	for _, key := range []string{"rooms", "room"} {
		rooms, _ := conf.ParseListValue(query.Get(key), ",")
		if len(rooms) < 1 {
			continue
		}
		for r, room := range rooms {
			// If room does not begin with a '#' let's prepend it
			if room[0] != '#' && room[0] != '!' {
				rooms[r] = "#" + room
			}
		}

		query.Set(key, conf.FormatListValue(rooms, ","))
	}

	updatedURL.RawQuery = query.Encode()
	return &updatedURL
}
