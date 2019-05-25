package main

import (
	"flag"
	"log"

	"github.com/JCGrant/twitch-paints/canvas"
	"github.com/JCGrant/twitch-paints/pixels"
	"github.com/JCGrant/twitch-paints/twitch"
)

var configPath = flag.String("config", "config.json", "path to config.json")

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
	pixels := make(chan pixels.Pixel)
	go twitch.Run(twitchConfig, twitchMessages)
	go func() {
		for msg := range twitchMessages {
			log.Println(msg)
			if p, err := parseMessage(msg); err == nil {
				log.Println("adding pixel")
				pixels <- p
			}
		}
	}()
	canvas.Run(pixels)
}

func parseMessage(msg string) (pixels.Pixel, error) {
	return pixels.FromString(msg)
}
