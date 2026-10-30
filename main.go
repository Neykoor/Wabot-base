package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var nonDigits = regexp.MustCompile(`\D`)

func askPhone() string {
	if env := os.Getenv("BOT_PHONE"); env != "" {
		return nonDigits.ReplaceAllString(env, "")
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Número del bot con código de país (ej: 521234567890): ")
		line, _ := reader.ReadString('\n')
		phone := nonDigits.ReplaceAllString(strings.TrimSpace(line), "")
		if len(phone) >= 8 {
			return phone
		}
		fmt.Println("Número inválido, intenta de nuevo.")
	}
}

func main() {
	ctx := context.Background()
	cfg := LoadConfig()

	dbLog := waLog.Stdout("DB", cfg.LogLevel, true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:"+cfg.SessionDB+"?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("Client", cfg.LogLevel, true))
	bot := NewBot(client, cfg)
	client.AddEventHandler(bot.HandleEvent)

	if client.Store.ID == nil {
		if err := client.Connect(); err != nil {
			panic(err)
		}
		phone := askPhone()
		time.Sleep(3 * time.Second)
		code, err := client.PairPhone(ctx, phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			panic(err)
		}
		fmt.Println("Código de vinculación:", code)
		fmt.Println("WhatsApp > Dispositivos vinculados > Vincular con el número de teléfono")
	} else {
		if err := client.Connect(); err != nil {
			panic(err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	client.Disconnect()
}
