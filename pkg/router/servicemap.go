package router

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/services/bark"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/discord"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/generic"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/googlechat"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/gotify"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/ifttt"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/join"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/logger"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/matrix"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/mattermost"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/ntfy"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/opsgenie"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/pushbullet"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/pushover"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/rocketchat"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/slack"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/smtp"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/teams"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/telegram"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/zulip"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

var serviceMap = map[string]func() types.Service{
	"bark":       func() types.Service { return &bark.Service{} },
	"discord":    func() types.Service { return &discord.Service{} },
	"generic":    func() types.Service { return &generic.Service{} },
	"gotify":     func() types.Service { return &gotify.Service{} },
	"googlechat": func() types.Service { return &googlechat.Service{} },
	"hangouts":   func() types.Service { return &googlechat.Service{} },
	"ifttt":      func() types.Service { return &ifttt.Service{} },
	"join":       func() types.Service { return &join.Service{} },
	"logger":     func() types.Service { return &logger.Service{} },
	"matrix":     func() types.Service { return &matrix.Service{} },
	"mattermost": func() types.Service { return &mattermost.Service{} },
	"ntfy":       func() types.Service { return &ntfy.Service{} },
	"opsgenie":   func() types.Service { return &opsgenie.Service{} },
	"pushbullet": func() types.Service { return &pushbullet.Service{} },
	"pushover":   func() types.Service { return &pushover.Service{} },
	"rocketchat": func() types.Service { return &rocketchat.Service{} },
	"slack":      func() types.Service { return &slack.Service{} },
	"smtp":       func() types.Service { return &smtp.Service{} },
	"teams":      func() types.Service { return &teams.Service{} },
	"telegram":   func() types.Service { return &telegram.Service{} },
	"zulip":      func() types.Service { return &zulip.Service{} },
}
