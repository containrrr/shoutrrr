package router

import (
	"github.com/containrrr/shoutrrr/pkg/services/bark"
	"github.com/containrrr/shoutrrr/pkg/services/discord"
	"github.com/containrrr/shoutrrr/pkg/services/generic"
	"github.com/containrrr/shoutrrr/pkg/services/googlechat"
	"github.com/containrrr/shoutrrr/pkg/services/gotify"
	"github.com/containrrr/shoutrrr/pkg/services/ifttt"
	"github.com/containrrr/shoutrrr/pkg/services/join"
	"github.com/containrrr/shoutrrr/pkg/services/logger"
	"github.com/containrrr/shoutrrr/pkg/services/matrix"
	"github.com/containrrr/shoutrrr/pkg/services/mattermost"
	"github.com/containrrr/shoutrrr/pkg/services/opsgenie"
	"github.com/containrrr/shoutrrr/pkg/services/pushbullet"
	"github.com/containrrr/shoutrrr/pkg/services/pushover"
	"github.com/containrrr/shoutrrr/pkg/services/rocketchat"
	"github.com/containrrr/shoutrrr/pkg/services/slack"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
	"github.com/containrrr/shoutrrr/pkg/services/teams"
	"github.com/containrrr/shoutrrr/pkg/services/telegram"
	"github.com/containrrr/shoutrrr/pkg/services/zulip"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

var serviceMap = map[string]func() t.Service{
	"bark":       func() t.Service { return &bark.Service{} },
	"discord":    func() t.Service { return &discord.Service{} },
	"generic":    func() t.Service { return &generic.Service{} },
	"gotify":     func() t.Service { return &gotify.Service{} },
	"googlechat": func() t.Service { return &googlechat.Service{} },
	"hangouts":   func() t.Service { return &googlechat.Service{} },
	"ifttt":      func() t.Service { return &ifttt.Service{} },
	"join":       func() t.Service { return &join.Service{} },
	"logger":     func() t.Service { return &logger.Service{} },
	"matrix":     func() t.Service { return &matrix.Service{} },
	"mattermost": func() t.Service { return &mattermost.Service{} },
	"opsgenie":   func() t.Service { return &opsgenie.Service{} },
	"pushbullet": func() t.Service { return &pushbullet.Service{} },
	"pushover":   func() t.Service { return &pushover.Service{} },
	"rocketchat": func() t.Service { return &rocketchat.Service{} },
	"slack":      func() t.Service { return &slack.Service{} },
	"smtp":       func() t.Service { return &smtp.Service{} },
	"teams":      func() t.Service { return &teams.Service{} },
	"telegram":   func() t.Service { return &telegram.Service{} },
	"zulip":      func() t.Service { return &zulip.Service{} },
}
