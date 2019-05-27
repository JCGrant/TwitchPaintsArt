package main

import (
	"flag"
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

	twitchMessages := make(chan string)
	canvasPixels := make(chan pixels.Pixel)
	backupPixels := make(chan pixels.Pixel)

	dbPath := "pixels.json"
	db := database.New(windowWidth, windowHeight)
	err = db.LoadPixels(dbPath)
	if err != nil {
		db.SavePixels(dbPath)
		err = db.LoadPixels(dbPath)
		if err != nil {
			log.Fatalln("loading databse failed: ", err)
		}
	}

	go twitch.Run(twitchConfig, twitchMessages)

	// Get twitch message, parse them to pixels, and pass pixels to other channels
	go func() {
		for msg := range twitchMessages {
			log.Println(msg)
			if p, err := parseMessage(msg); err == nil {
				log.Println("adding pixel")
				canvasPixels <- p
				backupPixels <- p
			}
		}
	}()

	go database.Run(backupPixels, dbPath, db)

	canvas.Run(canvasPixels, int32(windowWidth), int32(windowHeight), db.Pixels())
}

func parseMessage(msg string) (pixels.Pixel, error) {
	return pixels.FromString(msg)
}

func mustGetIntEnvVar(name string) int {
	valueStr := os.Getenv(name)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalln(err)
	}
	return value
}
