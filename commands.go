package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Command struct {
	Name        string
	Aliases     []string
	Description string
	GroupOnly   bool
	Run         func(c *Ctx)
}

func registerCommands(b *Bot) {
	b.Register(&Command{
		Name:        "ping",
		Aliases:     []string{"p"},
		Description: "Mide la latencia del bot",
		Run: func(c *Ctx) {
			sent := c.Event.Info.Timestamp
			c.Reply(fmt.Sprintf("Pong! %dms", time.Since(sent).Milliseconds()))
		},
	})

	b.Register(&Command{
		Name:        "menu",
		Aliases:     []string{"help", "ayuda"},
		Description: "Muestra los comandos disponibles",
		Run: func(c *Ctx) {
			seen := map[*Command]bool{}
			var list []*Command
			for _, cmd := range c.Bot.Commands {
				if !seen[cmd] {
					seen[cmd] = true
					list = append(list, cmd)
				}
			}
			sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })

			var sb strings.Builder
			sb.WriteString("*" + c.Bot.Config.BotName + "*\n\n")
			for _, cmd := range list {
				sb.WriteString(c.Bot.Config.Prefix + cmd.Name + " - " + cmd.Description + "\n")
			}
			c.Reply(sb.String())
		},
	})

	b.Register(&Command{
		Name:        "echo",
		Description: "Repite el texto que escribas",
		Run: func(c *Ctx) {
			if c.Text == "" {
				c.Reply("Uso: " + c.Bot.Config.Prefix + "echo <texto>")
				return
			}
			c.Reply(c.Text)
		},
	})

	b.Register(&Command{
		Name:        "id",
		Description: "Muestra el JID del chat y del remitente",
		Run: func(c *Ctx) {
			c.Reply("Chat: " + c.Chat.String() + "\nUsuario: " + c.Sender.String())
		},
	})

	b.Register(&Command{
		Name:        "grupo",
		Description: "Información básica del grupo",
		GroupOnly:   true,
		Run: func(c *Ctx) {
			info, err := c.Bot.Client.GetGroupInfo(context.Background(), c.Chat)
			if err != nil {
				c.Reply("No pude obtener la información del grupo.")
				return
			}
			c.Reply(fmt.Sprintf("Nombre: %s\nParticipantes: %d", info.Name, len(info.Participants)))
		},
	})
}
