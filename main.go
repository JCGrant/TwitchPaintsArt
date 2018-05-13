package main

import (
	"flag"
	"log"

	"github.com/JCGrant/twitch-paints/twitch"
)

var configPath = flag.String("config", "", "path to config.json")

func main() {
	flag.Parse()
	if *configPath == "" {
		log.Fatalln("no path to config.json supplied")
	}
	twitchConfig, err := twitch.ReadConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}
	twitchMessages := make(chan string)
	go twitch.Run(twitchConfig, twitchMessages)
	go func() {
		for msg := range twitchMessages {
			log.Println(msg)
		}
	}()
	// canvas.Run()
}
