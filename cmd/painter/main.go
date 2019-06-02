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
		time.Sleep(2500 * time.Millisecond)
		p := imagePixels[i]
		p.X += *leftOffset
		// Vertical coords of canvas are flipped, so must flip image pixels
		p.Y = *bottomOffset + imageHeight - 1 - p.Y
		fmt.Printf("%d %#v\n", i, p)
		twitchCommands <- pixelToCommand(p)
	}
}

func pixelToCommand(p pixels.Pixel) string {
	return fmt.Sprintf("!%d %d #%.2x%.2x%.2x", p.X, p.Y, p.Color.R, p.Color.G, p.Color.B)
}
