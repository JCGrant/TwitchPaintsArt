package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/JCGrant/twitch-paints/canvas"
	"github.com/JCGrant/twitch-paints/database"
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
	windowWidth := mustGetIntEnvVar("WIDTH")
	windowHeight := mustGetIntEnvVar("HEIGHT")

	twitchMessages := make(chan twitch.Message)
	canvasPixels := make(chan pixels.Pixel)
	backupPixels := make(chan pixels.Pixel)

	dbPath := "pixels.json"
	db := database.New(windowWidth, windowHeight)
	err = db.LoadPixels(dbPath)
	if err != nil {
		log.Printf("loading pixels failed: %s\n", err)
		log.Printf("recreating %s\n", dbPath)
		db.SavePixels(dbPath)
		err = db.LoadPixels(dbPath)
		if err != nil {
			log.Fatalln("loading databse failed: ", err)
		}
	}

	go twitch.Run(twitchConfig, nil, twitchMessages)

	// Get twitch message, parse them to pixels, and pass pixels to other channels
	go func() {
		for msg := range twitchMessages {
			if p, err := parseMessage(msg.Text); err == nil && isValidPixel(p, windowWidth, windowHeight) {
				log.Println(fmt.Sprintf("%s: %s", msg.Nickname, msg.Text))
				canvasPixels <- p
				backupPixels <- p
			}
		}
	}()

	go database.Run(backupPixels, dbPath, db)

	canvas.Run(canvasPixels, windowWidth, windowHeight, db.Pixels())
}

func parseMessage(msg string) (pixels.Pixel, error) {
	return pixels.FromString(msg)
}

func isValidPixel(p pixels.Pixel, windowWidth, windowHeight int) bool {
	return p.X >= 0 && p.X < windowWidth && p.Y >= 0 && p.Y < windowHeight
}

func mustGetIntEnvVar(name string) int {
	valueStr := os.Getenv(name)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalln(err)
	}
	return value
}
