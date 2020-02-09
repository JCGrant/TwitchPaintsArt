package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/JCGrant/twitch-paints/pixels"
)

// DB represents a storage of pixels
type DB struct {
	pixels [][]*pixels.Pixel
}

// New creates a DB
func New(width, height int) *DB {
	db := &DB{}
	for i := 0; i < height; i++ {
		row := make([]*pixels.Pixel, width)
		db.pixels = append(db.pixels, row)
	}
	return db
}

// AddPixel adds pixels to the DB
func (db *DB) AddPixel(p pixels.Pixel) {
	db.pixels[p.Y][p.X] = &p
}

// Pixels returns a list of pixels stored in the DB
func (db *DB) Pixels() []pixels.Pixel {
	var ps []pixels.Pixel
	for _, row := range db.pixels {
		for _, p := range row {
			if p == nil {
				continue
			}
			ps = append(ps, *p)
		}
	}
	return ps
}

// LoadPixels loads the pixels from a file and adds them to the DB
func (db *DB) LoadPixels(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file '%s' failed: %s", path, err)
	}
	var ps []pixels.Pixel
	err = json.Unmarshal(bs, &ps)
	if err != nil {
		return fmt.Errorf("unmarshaling database failed: %s", err)
	}
	for _, p := range ps {
		db.AddPixel(p)
	}
	return nil
}

// SavePixels saves the DB's pixels to a file
func (db *DB) SavePixels(path string) error {
	bs, err := json.Marshal(db.Pixels())
	if err != nil {
		return fmt.Errorf("marshalling database failed: %s", err)
	}
	err = ioutil.WriteFile(path, bs, 0644)
	if err != nil {
		return fmt.Errorf("writing file '%s' failed: %s", path, err)
	}
	return nil
}

// Run starts the DB, listening for any new pixels to store
func Run(pixels chan pixels.Pixel, savePath string, db *DB) error {
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case p := <-pixels:
			db.AddPixel(p)
		case <-ticker.C:
			log.Println("saving...")
			db.SavePixels(savePath)
		}
	}
}
