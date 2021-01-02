package router

import (
	"github.com/containrrr/shoutrrr/pkg/services/discord"
	"github.com/containrrr/shoutrrr/pkg/services/gotify"
	"github.com/containrrr/shoutrrr/pkg/services/hangouts"
	"github.com/containrrr/shoutrrr/pkg/services/ifttt"
	"github.com/containrrr/shoutrrr/pkg/services/join"
	"github.com/containrrr/shoutrrr/pkg/services/logger"
	"github.com/containrrr/shoutrrr/pkg/services/mattermost"
	"github.com/containrrr/shoutrrr/pkg/services/pushbullet"
	"github.com/containrrr/shoutrrr/pkg/services/pushover"
	"github.com/containrrr/shoutrrr/pkg/services/rocketchat"
	"github.com/containrrr/shoutrrr/pkg/services/slack"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
	"github.com/containrrr/shoutrrr/pkg/services/teams"
	"github.com/containrrr/shoutrrr/pkg/services/telegram"
	"github.com/containrrr/shoutrrr/pkg/services/xmpp"
	"github.com/containrrr/shoutrrr/pkg/services/zulip"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

var serviceMap = map[string]func() t.Service{
	"discord":    func() t.Service { return &discord.Service{} },
	"pushover":   func() t.Service { return &pushover.Service{} },
	"slack":      func() t.Service { return &slack.Service{} },
	"teams":      func() t.Service { return &teams.Service{} },
	"telegram":   func() t.Service { return &telegram.Service{} },
	"smtp":       func() t.Service { return &smtp.Service{} },
	"ifttt":      func() t.Service { return &ifttt.Service{} },
	"gotify":     func() t.Service { return &gotify.Service{} },
	"logger":     func() t.Service { return &logger.Service{} },
	"xmpp":       func() t.Service { return &xmpp.Service{} },
	"pushbullet": func() t.Service { return &pushbullet.Service{} },
	"mattermost": func() t.Service { return &mattermost.Service{} },
	"hangouts":   func() t.Service { return &hangouts.Service{} },
	"zulip":      func() t.Service { return &zulip.Service{} },
	"join":       func() t.Service { return &join.Service{} },
	"rocketchat": func() t.Service { return &rocketchat.Service{} },
}
