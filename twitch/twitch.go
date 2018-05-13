package twitch

import (
	"encoding/json"
	"log"
	"os"

	irc "github.com/thoj/go-ircevent"
)

const (
	twitchIRC = "irc.twitch.tv:6667"
)

// Config contains the data needed to connect to a Twitch IRC channel
type Config struct {
	Nick    string `json:"nick"`
	Pass    string `json:"pass"`
	Channel string `json:"channel"`
}

// ReadConfig will create a Config from a path to a config.json
func ReadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

// Run starts an IRC client connected to a Twitch channel.
// This function blocks.
func Run(config Config, messages chan string) {
	irccon := irc.IRC(config.Nick, config.Nick)
	irccon.Password = config.Pass
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(config.Channel) })
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		msg := e.Message()
		messages <- msg
	})
	err := irccon.Connect(twitchIRC)
	if err != nil {
		log.Fatalln(err)
	}
	irccon.Loop()
}
