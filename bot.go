package main

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type Bot struct {
	Client   *whatsmeow.Client
	Config   Config
	Commands map[string]*Command
}

type Ctx struct {
	Bot    *Bot
	Event  *events.Message
	Chat   types.JID
	Sender types.JID
	Text   string
	Args   []string
}

func NewBot(client *whatsmeow.Client, cfg Config) *Bot {
	b := &Bot{
		Client:   client,
		Config:   cfg,
		Commands: make(map[string]*Command),
	}
	registerCommands(b)
	return b
}

func (b *Bot) Register(cmd *Command) {
	b.Commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		b.Commands[alias] = cmd
	}
}

func (b *Bot) HandleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go b.onMessage(v)
	case *events.PairSuccess:
		fmt.Println("Vinculado correctamente:", v.ID.String())
	case *events.Connected:
		fmt.Println("Bot conectado")
	case *events.Disconnected:
		fmt.Println("Bot desconectado")
	case *events.LoggedOut:
		fmt.Println("Sesión cerrada desde el teléfono, borra la base de sesión y vuelve a vincular")
	}
}

func extractText(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}
	if t := msg.GetConversation(); t != "" {
		return t
	}
	if t := msg.GetExtendedTextMessage().GetText(); t != "" {
		return t
	}
	if t := msg.GetImageMessage().GetCaption(); t != "" {
		return t
	}
	if t := msg.GetVideoMessage().GetCaption(); t != "" {
		return t
	}
	return ""
}

func (b *Bot) onMessage(evt *events.Message) {
	if evt.Info.IsFromMe {
		return
	}

	text := strings.TrimSpace(extractText(evt.Message))
	if text == "" || !strings.HasPrefix(text, b.Config.Prefix) {
		return
	}

	fields := strings.Fields(strings.TrimPrefix(text, b.Config.Prefix))
	if len(fields) == 0 {
		return
	}

	name := strings.ToLower(fields[0])
	cmd, ok := b.Commands[name]
	if !ok {
		return
	}

	c := &Ctx{
		Bot:    b,
		Event:  evt,
		Chat:   evt.Info.Chat,
		Sender: evt.Info.Sender,
		Text:   strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(text, b.Config.Prefix), fields[0])),
		Args:   fields[1:],
	}

	if cmd.GroupOnly && !evt.Info.IsGroup {
		c.Reply("Este comando solo funciona en grupos.")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error en comando", name, ":", r)
			c.Reply("Ocurrió un error ejecutando el comando.")
		}
	}()

	cmd.Run(c)
}

func (c *Ctx) Reply(text string) {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(c.Event.Info.ID),
				Participant:   proto.String(c.Sender.String()),
				QuotedMessage: c.Event.Message,
			},
		},
	}
	if _, err := c.Bot.Client.SendMessage(context.Background(), c.Chat, msg); err != nil {
		fmt.Println("Error enviando mensaje:", err)
	}
}

func (c *Ctx) Send(text string) {
	msg := &waE2E.Message{Conversation: proto.String(text)}
	if _, err := c.Bot.Client.SendMessage(context.Background(), c.Chat, msg); err != nil {
		fmt.Println("Error enviando mensaje:", err)
	}
}
