package main

import (
	"github.com/thoj/go-ircevent"
	"github.com/BurntSushi/toml"
	"fmt"
	"bufio"
	"os"
	"strings"
)

var ircObj *irc.Connection

type Auth struct {
	AuthToken string
}

func onNotice(event *irc.Event) {
	fmt.Println(event)
	fmt.Println(event.Raw)
}

func onPrivMessage(event *irc.Event) {
	fmt.Println(event)
	fmt.Println(event.Raw)
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

	ircObj = irc.IRC("Cryptiobot", "Cryptiobot")
	ircObj.Password = auth.AuthToken
	ircObj.AddCallback("NOTICE", onNotice)
	ircObj.AddCallback("PRIVMSG", onPrivMessage)

	ircObj.RunCallbacks(&irc.Event{})


	if err := ircObj.Connect("irc.chat.twitch.tv:6667"); err != nil {
		fmt.Printf("Error connecting to Twitch IRC server: %s.\n", err)
		return
	}

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
