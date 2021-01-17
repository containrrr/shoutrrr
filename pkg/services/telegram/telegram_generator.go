package telegram

import (
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util/generator"

	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const sessionTokenLength = 32

type Generator struct {
	ud            *generator.UserDialog
	client        *Client
	chats         []string
	chatNames     []string
	done          bool
	chatMessages  map[string]int64
	owner         *User
	sessionToken  string
	statusMessage int64
	botName       string
}

func (g *Generator) Generate(service types.Service, props map[string]string, args []string) (types.ServiceConfig, error) {
	var config Config

	sessionTokenBytes := make([]byte, sessionTokenLength)
	if _, err := rand.Read(sessionTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate session token: %v", err)
	}
	g.sessionToken = hex.EncodeToString(sessionTokenBytes)

	g.ud = generator.NewUserDialog(os.Stdin, os.Stdout, props)
	ud := g.ud

	ud.Writeln("To start we need your bot token. If you haven't created a bot yet, you can use this link:")
	ud.Writeln("  https://t.me/botfather?start")
	ud.Writeln("")

	token := ud.QueryString("Enter your bot token: ", generator.ValidateFormat(IsTokenValid), "token")

	ud.Writeln("Fetching bot info...")
	ud.Writeln("Session token: %v", g.sessionToken)

	g.client = &Client{token: token}
	botInfo, err := g.client.GetBotInfo()
	if err != nil {
		return &Config{}, err
	}

	err = g.client.SetCommands(map[string]string{
		"start":   "Start selecting chats for send list",
		"include": "Includes the chat in the Shoutrrr send list",
		"remove":  "Removes the chat from the Shoutrrr send list",
		"done":    "Stops listening for new chats",
	})
	if err != nil {
		return &Config{}, err
	}

	// fmt.Printf("%#v", botInfo)
	g.botName = botInfo.Username
	ud.Writeln("")
	ud.Writeln("Okay! Use the following link to add the bot to the groups you want to include:")
	//ud.Writeln("  https://t.me/%v?startgroup=%v", botInfo.Username, g.sessionToken)
	ud.Writeln("  https://t.me/%v?start=%v", botInfo.Username, g.sessionToken)
	ud.Writeln("")
	ud.Writeln("...or you can send the /start command to @%v:", g.botName)
	ud.Writeln("  /start@%v %v", g.botName, g.sessionToken)
	ud.Writeln("")
	//ud.Writeln("You can also use the slash commands:")
	//ud.Writeln("  /include%v", botName)
	//ud.Writeln("  /remove%v", botName)
	//ud.Writeln("")
	ud.Writeln("When you are done, send the /done command to the bot:")
	ud.Writeln("  /done@%v", g.botName)
	ud.Writeln("")

	g.chatMessages = make(map[string]int64, 0)
	g.done = false
	lastUpdate := 0

	for !g.done {

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

			if update.CallbackQuery != nil {
				g.HandleCallback(update.CallbackQuery)
			} else if update.InlineQuery != nil {
				ud.Writeln("Got inline query from @%v, which is not supported", message.From.Username)
			} else if message != nil {
				// fmt.Printf("MSG#%v [%v:%v] %v: %v\n", message.MessageID, message.Chat.Type, message.Chat.ID, message.From.Username, message.Text)
				// chatId := strconv.FormatInt(message.Chat.ID, 10)
				cmd, params, err := g.client.ParseCommand(message, g.botName, message.Chat.Type == "private")
				if err != nil {
					source := message.Chat.Username
					if message.From != nil {
						source = message.From.Username
					}
					ud.Writeln("Got invalid command '%v' from @%v: %v", message.Text, source, err)
				}

				// fmt.Printf("CMD: %#v PARAMS: %#v\n", cmd, params)
				if cmd != "" {
					g.Reply(message, g.HandleCommand(cmd, params, message.Chat, message.From))
				}
			} else {
				fmt.Printf("Unknown payload:\n%#v", update)
			}

		}

		// time.Sleep(time.Second * 10)
	}

	// Notify API that we got the updates
	if _, err = g.client.GetUpdates(lastUpdate, 0, 0, nil); err != nil {
		g.ud.Writeln("Failed to mark last updates as received")
	}

	if err := g.client.SetCommands(map[string]string{}); err != nil {
		g.ud.Writeln("Failed to reset  bot commands")
	}

	config = Config{
		Notification: true,
		Token:        token,
		Channels:     g.chats,
	}

	return &config, nil
}

func (g *Generator) AddChat(chat *Chat) (result string) {
	id := strconv.FormatInt(chat.ID, 10)
	name := chat.Name()

	for _, c := range g.chats {
		if c == id {
			return fmt.Sprintf("Chat '%v' is already selected!", name)
		}
	}
	g.chats = append(g.chats, id)
	g.chatNames = append(g.chatNames, name)

	g.UpdateMessage()

	return fmt.Sprintf("Added new chat '%v'!", name)
}

func (g *Generator) DelChat(chat *Chat) (result string) {
	id := strconv.FormatInt(chat.ID, 10)

	for i, chatId := range g.chats {
		if chatId == id {
			g.chats = append(g.chats[:i], g.chats[i+1:]...)
			g.chatNames = append(g.chatNames[:i], g.chatNames[i+1:]...)
			return fmt.Sprintf("Removed chat '%v'!", chat.Name())
		}
	}

	g.UpdateMessage()
	return fmt.Sprintf("Chat '%v' not selected!", chat.Name())
}

func (g *Generator) Reply(original *Message, text string) {
	// text = strings.ReplaceAll(text, "!", "\\!")
	if _, err := g.client.Reply(original, text); err != nil {
		g.ud.Writeln("Failed to send reply: %v", err)
	}
}

func (g *Generator) HandleCallback(cq *CallbackQuery) {
	params := strings.Split(cq.Data, " ")
	if len(params) < 1 {
		return
	}
	cmd := params[0]
	params = params[1:]

	result := g.HandleCommand(cmd, params, cq.Message.Chat, cq.From)
	if err := g.client.AnswerCallbackQuery(&CallbackQueryAnswer{
		CallbackQueryID: cq.ID,
		Text:            result,
		ShowAlert:       false,
	}); err != nil {
		g.ud.Writeln("Failed to answer callback query: %v", err)
	}
}

func (g *Generator) HandleCommand(cmd string, params []string, chat *Chat, sender *User) string {
	if cmd != "start" {
		if g.owner == nil {
			return "No selection has been started yet"
		} else if sender != nil && sender.ID != g.owner.ID {
			return fmt.Sprintf("Nah, ask @%v", g.owner.Username)
		}
	}

	switch cmd {
	case "start":
		if len(params) < 1 {
			return "Hi!"
		} else {
			if g.owner == nil {
				if params[0] == g.sessionToken {
					return g.SetOwner(sender)
				}
				return "Not currently selecting, perhaps you forgot to start the session?"
			} else {
				return g.HandleCommand(params[0], params[1:], chat, sender)
			}
		}
	case "include", "add":
		return g.AddChat(chat)
	case "exclude", "remove":
		return g.DelChat(chat)
	case "done":
		g.done = true
		g.UpdateMessage()
		return "Awesome! The configuration URL will be written to the console output"
	default:
		return fmt.Sprintf("Invalid command `%v`", cmd)
	}
}

func (g *Generator) UpdateMessage() {

	s := "s"
	if len(g.chats) == 1 {
		s = ""
	}
	text := fmt.Sprintf("🐹📣 %v chat%v selected\\!", len(g.chats), s)

	keys := make([][]InlineKey, 0, len(g.chats)+1)

	if g.done {
		text += "\n\nDone\\! 🎉"
	} else {
		for c, chat := range g.chats {
			keys = append(keys, []InlineKey{{
				Text:         "Remove " + g.chatNames[c],
				CallbackData: "remove " + chat,
			}})
		}
		keys = append(keys, []InlineKey{
			{Text: "Add a channel", SwitchInlineQuery: "/include"},
			// {Text: "Add 2", Url: fmt.Sprintf("tg://resolve/%v?startgroup=include", g.botName)},
			{Text: "Add a group", Url: fmt.Sprintf("https://t.me/%v?startgroup=include", g.botName)},
			{Text: "Done", CallbackData: "done"},
		})
	}

	ownerChat := strconv.FormatInt(g.owner.ID, 10)
	if g.statusMessage == 0 {

		message, err := g.client.SendMessage(&SendMessagePayload{
			Text:        text,
			ID:          ownerChat,
			ParseMode:   ParseModes.MarkdownV2.String(),
			ReplyMarkup: &ReplyMarkup{InlineKeyboard: keys},
		})
		if err != nil {
			g.ud.Writeln("Failed to send message: %v", err)
		} else {
			g.statusMessage = message.MessageID
		}

	} else {
		replyMarkup, err := json.Marshal(&ReplyMarkup{InlineKeyboard: keys})
		if err == nil {
			err = g.client.UpdateMessage(&UpdateMessagePayload{
				Text:        text,
				ChatID:      ownerChat,
				MessageID:   g.statusMessage,
				ParseMode:   ParseModes.MarkdownV2.String(),
				ReplyMarkup: string(replyMarkup),
			})
		}
		if err != nil {
			g.ud.Writeln("Failed to update messages: %v", err)
		}
	}
}

func (g *Generator) SetOwner(owner *User) string {
	g.owner = owner

	g.UpdateMessage()

	return fmt.Sprintf("Started session for @%v", owner.Username)
}
