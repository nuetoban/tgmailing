package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akamensky/argparse"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/nuetoban/tgmailing/input/envsrc"
	"github.com/nuetoban/tgmailing/input/jsonsrc"
	"github.com/nuetoban/tgmailing/input/linesfile"
	"github.com/nuetoban/tgmailing/input/sqlsrc"
	"github.com/nuetoban/tgmailing/sending"
	"github.com/nuetoban/tgmailing/view"
)

const (
	JSONSRC      = "JSONFILE"
	LINESFILESRC = "LINESFILE"
	ENVSRC       = "ENV"
	PGSQLSRC     = "PGSQL"
)

func main() {
	parser := argparse.NewParser("tgmailing", "Sends post to Telegram chats via Bots")

	N := ""

	// Define args
	// Input sources
	adSrc := parser.Selector(N, "ad-src", []string{JSONSRC, PGSQLSRC}, &argparse.Options{Default: JSONSRC, Help: "Ad source"})
	botsSrc := parser.Selector(N, "bots-src", []string{ENVSRC, LINESFILESRC, JSONSRC, PGSQLSRC}, &argparse.Options{Default: LINESFILESRC, Help: "Bots source"})
	chatSrc := parser.Selector(N, "chats-src", []string{LINESFILESRC, PGSQLSRC}, &argparse.Options{Default: LINESFILESRC, Help: "Chats source"})

	// Files
	file := parser.String(N, "ad-file", &argparse.Options{Help: "Path to Ad file"})
	botsFile := parser.String(N, "bots-file", &argparse.Options{Help: "Path to bots file"})
	chatsFile := parser.String(N, "chats-file", &argparse.Options{Help: "Path to chats file"})

	// Queries
	botsQuery := parser.String(N, "bots-query", &argparse.Options{Help: "SQL query to fetch Bots"})
	chatsQuery := parser.String(N, "chats-query", &argparse.Options{Help: "SQL query to fetch Chats"})
	postQuery := parser.String(N, "ad-query", &argparse.Options{Help: "SQL query to fetch Ad post"})

	// Env prefixes
	botsEnvPrefix := parser.String(N, "bots-db-env-prefix", &argparse.Options{Help: "Prefix for DB env credentials for Bots", Default: "SENDER_"})
	chatsEnvPrefix := parser.String(N, "chats-db-env-prefix", &argparse.Options{Help: "Prefix for DB env credentials for Chats", Default: "SENDER_"})
	postEnvPrefix := parser.String(N, "ad-db-env-prefix", &argparse.Options{Help: "Prefix for DB env credentials for Ad post", Default: "SENDER_"})

	// Flags/chats
	noServer := parser.Flag(N, "no-server", &argparse.Options{Help: "Disable metrics server"})
	metricsPort := parser.Int("m", "metrics-port", &argparse.Options{Help: "Metrics server port", Default: 9090})
	enableStartNotification := parser.Flag("n", "start-notification", &argparse.Options{Help: "Send message to chat on start"})
	enableFinishNotification := parser.Flag("f", "finish-notification", &argparse.Options{Help: "Send message to chat on finish"})
	enableFinishEachNotification := parser.Flag(N, "each-finish-notification", &argparse.Options{Help: "Send message to chat on finish for each bot"})
	notificationChat := parser.Int(N, "notification-chat", &argparse.Options{Help: "Chat to send notifications"})
	serviceChat := parser.Int("s", "service-chat", &argparse.Options{Help: "Chat to send files", Required: true})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	var (
		bots sending.BotsProvider
		chat sending.ChatsProvider
		ad   sending.AdProvider
	)

	// Set ad post source
	switch *adSrc {
	case JSONSRC:
		if file == nil || *file == "" {
			fmt.Println("The argument --ad-file should be used with this Ad source")
			return
		}
		af, err := os.Open(*file)
		if err != nil {
			log.Fatalln(err)
		}
		ad, err = jsonsrc.New(af)
		if err != nil {
			log.Fatalln(err)
		}
	case PGSQLSRC:
		if postQuery == nil || *postQuery == "" {
			fmt.Println("The argument --ad-query should be used with this Ad source")
			return
		}

		ad, err = sqlsrc.NewPostgreSQL(sqlsrc.Config{PostSelectQuery: *postQuery}, false, *postEnvPrefix)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Set chats source
	switch *chatSrc {
	case LINESFILESRC:
		if chatsFile == nil || *chatsFile == "" {
			fmt.Println("The argument --chats-file should be used with this Chats source")
			return
		}
		cf, err := os.Open(*chatsFile)
		if err != nil {
			log.Fatalln(err)
		}
		chat, err = linesfile.New(cf)
		if err != nil {
			log.Fatalln(err)
		}
	case PGSQLSRC:
		if chatsQuery == nil || *chatsQuery == "" {
			fmt.Println("The argument --chats-query should be used with this Bots source")
			return
		}

		chat, err = sqlsrc.NewPostgreSQL(sqlsrc.Config{ChatsSelectQuery: *chatsQuery}, false, *chatsEnvPrefix)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Set bots source
	switch *botsSrc {
	case ENVSRC:
		bots, err = envsrc.New("SENDER_TOKEN")
		if err != nil {
			log.Fatalln(err)
		}
	case LINESFILESRC:
		if botsFile == nil || *botsFile == "" {
			fmt.Println("The argument --bots-file should be used with this Bots source")
			return
		}
		bf, err := os.Open(*botsFile)
		if err != nil {
			log.Fatalln(err)
		}
		bots, err = linesfile.New(bf)
		if err != nil {
			log.Fatalln(err)
		}
	case PGSQLSRC:
		if botsQuery == nil || *botsQuery == "" {
			fmt.Println("The argument --bots-query should be used with this Bots source")
			return
		}

		bots, err = sqlsrc.NewPostgreSQL(sqlsrc.Config{BotsSelectQuery: *botsQuery}, false, *botsEnvPrefix)
		if err != nil {
			log.Fatalln(err)
		}
	case JSONSRC:
		if botsFile == nil || *botsFile == "" {
			fmt.Println("The argument --bots-file should be used with this Bots source")
			return
		}
		af, err := os.Open(*botsFile)
		if err != nil {
			log.Fatalln(err)
		}
		bots, err = jsonsrc.New(af)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// New sender
	swb := sending.NewSenderWithBots(bots, chat, ad)

	if (*enableStartNotification || *enableFinishNotification || *enableFinishEachNotification) && *notificationChat == 0 {
		fmt.Println("The argument --notification-chat is required if you ask to notify")
		return
	}

	if *enableStartNotification {
		swb.SetAfterStartHook(sending.NotifyDevChatOnStart(int64(*notificationChat)))
	}

	if *enableFinishNotification {
		swb.SetAfterFinishHook(sending.NotifyDevChatOnFinish(int64(*notificationChat)))
	}

	if *enableFinishEachNotification {
		swb.SetAfterFinishEachHook(sending.NotifyDevChatOnFinishEach(int64(*notificationChat)))
	}

	swb.SetServiceChat(int64(*serviceChat))

	if !*noServer {
		// Add metrics endpoints
		http.HandleFunc("/metrics/json", view.NewJSON(swb).StatisticsHTTP)
		http.HandleFunc("/metrics", view.NewPrometheus(swb).StatisticsHTTP)

		// Run server
		go http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
	}

	// Run program
	err = swb.Run()
	if err != nil {
		log.Fatalf("error during creation sender with bots: %v", err)
	}
}
