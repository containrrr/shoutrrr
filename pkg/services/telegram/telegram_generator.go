package telegram

import (
	f "github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/generator"
	"os/signal"
	"syscall"

	"fmt"
	"os"
	"strconv"
)

// Generator is the telegram-specific URL generator
type Generator struct {
	ud            *generator.UserDialog
	client        *Client
	chats         []string
	chatNames     []string
	chatTypes     []string
	done          bool
	owner         *User
	statusMessage int64
	botName       string
}

// Generate a telegram Shoutrrr configuration from a user dialog
func (g *Generator) Generate(_ types.Service, props map[string]string, _ []string) (types.ServiceConfig, error) {
	var config Config

	g.ud = generator.NewUserDialog(os.Stdin, os.Stdout, props)
	ud := g.ud

	ud.Writeln("To start we need your bot token. If you haven't created a bot yet, you can use this link:")
	ud.Writeln("  %v", f.ColorizeLink("https://t.me/botfather?start"))
	ud.Writeln("")

	token := ud.QueryString("Enter your bot token:", generator.ValidateFormat(IsTokenValid), "token")

	ud.Writeln("Fetching bot info...")
	// ud.Writeln("Session token: %v", g.sessionToken)

	g.client = &Client{token: token}
	botInfo, err := g.client.GetBotInfo()
	if err != nil {
		return &Config{}, err
	}

	g.botName = botInfo.Username
	ud.Writeln("")
	ud.Writeln("Okay! %v will listen for any messages in PMs and group chats it is invited to.",
		f.ColorizeString("@", g.botName, ":"))

	g.done = false
	lastUpdate := 0

	signals := make(chan os.Signal, 1)

	// Subscribe to system signals
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for !g.done {

		ud.Writeln("Waiting for messages to arrive...")

		updates, err := g.client.GetUpdates(lastUpdate, 10, 120, nil)
		if err != nil {
			panic(err)
		}

		for _, update := range updates {
			lastUpdate = update.UpdateID + 1

			message := update.Message
			if update.ChannelPost != nil {
				message = update.ChannelPost
			}

			if message != nil {
				chat := message.Chat

				source := message.Chat.Username
				if message.From != nil {
					source = message.From.Username
				}
				ud.Writeln("Got Message '%v' from @%v in %v chat %v",
					f.ColorizeString(message.Text),
					f.ColorizeProp(source),
					f.ColorizeEnum(chat.Type),
					f.ColorizeNumber(chat.ID))
				ud.Writeln(g.addChat(chat))
			} else {
				ud.Writeln("Got unknown Update. Ignored!")
			}
		}

		ud.Writeln("")

		g.done = !ud.QueryBool(fmt.Sprintf("Got %v chat ID(s) so far. Want to add some more?",
			f.ColorizeNumber(len(g.chats))), "")
	}

	ud.Writeln("")
	ud.Writeln("Cleaning up the bot session...")

	// Notify API that we got the updates
	if _, err = g.client.GetUpdates(lastUpdate, 0, 0, nil); err != nil {
		g.ud.Writeln("Failed to mark last updates as received: %v", f.ColorizeError(err))
	}

	if len(g.chats) < 1 {
		return nil, fmt.Errorf("no chats were selected")
	}

	ud.Writeln("Selected chats:")

	for i, id := range g.chats {
		name := g.chatNames[i]
		chatType := g.chatTypes[i]
		ud.Writeln("  %v (%v) %v", f.ColorizeNumber(id), f.ColorizeEnum(chatType), f.ColorizeString(name))
	}

	ud.Writeln("")

	config = Config{
		Notification: true,
		Token:        token,
		Chats:        g.chats,
	}

	return &config, nil
}

func (g *Generator) addChat(chat *chat) (result string) {
	id := strconv.FormatInt(chat.ID, 10)
	name := chat.Name()

	for _, c := range g.chats {
		if c == id {
			return fmt.Sprintf("chat %v is already selected!", f.ColorizeString(name))
		}
	}
	g.chats = append(g.chats, id)
	g.chatNames = append(g.chatNames, name)
	g.chatTypes = append(g.chatTypes, chat.Type)

	return fmt.Sprintf("Added new chat %v!", f.ColorizeString(name))
}
