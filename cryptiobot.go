package main

import (
	"github.com/thoj/go-ircevent"
	"github.com/BurntSushi/toml"
	"fmt"
	"bufio"
	"os"
	"strings"
	"time"
	"net/http"
	"encoding/json"
	"strconv"
)
// Global configuration object.
var config Config

// Global IRC connection object.
var ircObj *irc.Connection

type Config struct {
	UserName string
	Nick string
	ServerAddress string
	Channels []string
	UpdateIntervalSeconds int
	Currency string
	CryptoSymbols []string
}

type Auth struct {
	AuthToken string
}

type Ticker struct {
	Ticker CryptoData
}

type CryptoData struct {
	Base string
	Price string
	Volume string
	Change string
}

func onNotice(event *irc.Event) {
	fmt.Println(event)
	fmt.Println(event.Raw)
}

func onPrivMessage(event *irc.Event) {
	fmt.Println(event)
	fmt.Println(event.Raw)
}

func startLoop() {
	ticker := time.NewTicker(time.Second * time.Duration(config.UpdateIntervalSeconds))

	// Asynchronously query for crypto data.
	go func() {
		queryCryptos()
		for range ticker.C {
			queryCryptos()
		}
	}()

}

func queryCryptos() {
	const API_URI string = "https://api.cryptonator.com/api/ticker/"
	for _, symbol := range config.CryptoSymbols {
		resp, err := http.Get(fmt.Sprintf("%s%s-%s", API_URI, symbol, config.Currency))
		if err != nil {
			fmt.Printf("Error while fetching %s data: %s.\n", symbol, err)
			return
		}

		var data * Ticker = new(Ticker)
		err = json.NewDecoder(resp.Body).Decode(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(*data)

		for _, ch := range config.Channels {
			price, _ := strconv.ParseFloat(data.Ticker.Price, 32)
			change, _ := strconv.ParseFloat(data.Ticker.Change, 32)
			output := fmt.Sprintf("%s is currently priced at %.2f. %s has changed by %.2f%% in the last hour.\n",
				data.Ticker.Base, price, data.Ticker.Base, change)

			fmt.Println(output)
			ircObj.Noticef(ch, output)
		}
	}
}

func commandHandle(cmd string) {
	// Commands must start with '!'
	if len(cmd) <= 0 || cmd[0] != '!' {
		return
	}

	// Trim '!'.
	cmd = cmd[1:]

	// Split space-separated input into string array.
	input := strings.Fields(cmd)

	cmd = input[0]
	switch cmd {
	case "w":
		if len(input) != 3 {
			fmt.Printf("Error: Command '%s' needs 3 args!\n", cmd)
			return
		}

		ircObj.SendRawf("PRIVMSG #jtv :/w %s %s", input[1], input[2])
		break
	}
}

func main() {
	var auth Auth
	if _, err := toml.DecodeFile("auth.toml", &auth); err != nil {
		fmt.Printf("Error while parsing auth TOML: %s.\n", err)
		return
	}

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Printf("Error while parsing configuration TOML: %s.\n", err)
		return
	}

	ircObj = irc.IRC(config.Nick, config.UserName)
	ircObj.Password = auth.AuthToken

	ircObj.AddCallback("NOTICE", onNotice)
	ircObj.AddCallback("PRIVMSG", onPrivMessage)

	if err := ircObj.Connect(config.ServerAddress); err != nil {
		fmt.Printf("Error connecting to IRC server: %s.\n", err)
		return
	}

	for _, ch := range config.Channels {
		ircObj.Join(ch)
	}

	// Begin information dump loop.
	startLoop()

	reader := bufio.NewReader(os.Stdin)
	for {
		if text, err := reader.ReadString('\n'); err == nil {
			// Trim line break char.
			trimmedText := strings.TrimRight(text, "\n")

			// Check for quit command.
			if trimmedText == "q" || trimmedText == "quit" {
				return
			}

			// Attempt to handle command.
			commandHandle(trimmedText)
		} else {
			fmt.Printf("Error reading input: %s.\n", err)
		}
	}
}