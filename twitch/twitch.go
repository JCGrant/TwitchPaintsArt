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
// Will output messages from the Twitch channel on to the out chan.
// Will send messages from the in chan to the Twitch channel.
// This function blocks.
func Run(config Config, in chan string, out chan string) {
	irccon := irc.IRC(config.Nick, config.Nick)
	irccon.Password = config.Pass
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(config.Channel) })
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		out <- e.Message()
	})
	err := irccon.Connect(twitchIRC)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for msg := range in {
			irccon.Privmsg(config.Channel, msg)
		}
	}()
	irccon.Loop()
}
