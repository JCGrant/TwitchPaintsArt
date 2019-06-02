package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/JCGrant/twitch-paints/pixels"
	"github.com/JCGrant/twitch-paints/twitch"
)

var configPath = flag.String("config", "config.json", "path to config.json")
var imagePath = flag.String("i", "img.png", "path to image")
var leftOffset = flag.Int("left", 0, "how far from left")
var bottomOffset = flag.Int("bottom", 0, "how far from bottom")
var startPixelIndex = flag.Int("start", 0, "at what pixel to start painting")
var delay = flag.Int("delay", 2500, "the amount of delay between sending Twitch commands")
var dryRun = flag.Bool("dry", false, "will run the painter, but will not send to Twitch")

func main() {
	flag.Parse()
	imagePixels, _, imageHeight, err := pixels.ImageFilePixels(*imagePath)
	if err != nil {
		log.Fatalln(err)
	}
	twitchConfig, err := twitch.ReadConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}
	twitchCommands := make(chan string)
	go twitch.Run(twitchConfig, twitchCommands, nil)
	for i := *startPixelIndex; i < len(imagePixels); i++ {
		time.Sleep(time.Duration(*delay) * time.Millisecond)
		p := imagePixels[i]
		p.X += *leftOffset
		// Vertical coords of canvas are flipped, so must flip image pixels
		p.Y = *bottomOffset + imageHeight - 1 - p.Y
		fmt.Printf("%d %#v\n", i, p)
		if !*dryRun {
			twitchCommands <- pixelToCommand(p)
		}
	}
}

func pixelToCommand(p pixels.Pixel) string {
	return fmt.Sprintf("!%d %d #%.2x%.2x%.2x", p.X, p.Y, p.Color.R, p.Color.G, p.Color.B)
}
