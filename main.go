package main

import (
	"bufio"
	"context"
	"errors"
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
	"go.mau.fi/whatsmeow/types/events"
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

func connectSafe(client *whatsmeow.Client) error {
	err := client.Connect()
	if err != nil && !errors.Is(err, whatsmeow.ErrAlreadyConnected) {
		return err
	}
	return nil
}

func pairWithCode(ctx context.Context, client *whatsmeow.Client, phone string, paired, dropped chan struct{}) bool {
	for attempt := 1; attempt <= 10; attempt++ {
		if err := connectSafe(client); err != nil {
			fmt.Println("Error conectando, reintentando:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		time.Sleep(3 * time.Second)

		select {
		case <-dropped:
		default:
		}

		code, err := client.PairPhone(ctx, phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			fmt.Println("Error generando código, reintentando:", err)
			client.Disconnect()
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("Código de vinculación:", code)
		fmt.Println("WhatsApp > Dispositivos vinculados > Vincular con el número de teléfono")
		fmt.Println("Tienes unos 2 minutos. Si se corta la conexión se genera un código nuevo solo.")

		select {
		case <-paired:
			return true
		case <-dropped:
			fmt.Println("Se cortó la conexión antes de vincular, generando código nuevo...")
			time.Sleep(3 * time.Second)
		case <-time.After(150 * time.Second):
			fmt.Println("El código expiró, generando uno nuevo...")
			client.Disconnect()
			time.Sleep(2 * time.Second)
		}
	}
	return false
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

	paired := make(chan struct{}, 1)
	dropped := make(chan struct{}, 1)
	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.PairSuccess:
			select {
			case paired <- struct{}{}:
			default:
			}
		case *events.Disconnected:
			select {
			case dropped <- struct{}{}:
			default:
			}
		}
	})

	if client.Store.ID == nil {
		phone := askPhone()
		if !pairWithCode(ctx, client, phone, paired, dropped) {
			fmt.Println("No se pudo vincular. Revisa tu conexión y vuelve a intentarlo.")
			return
		}
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
