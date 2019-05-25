package pixels

import (
	"fmt"
	"regexp"
	"strconv"
)

// Pixel represents a single pixel on a screen
type Pixel struct {
	X     int32
	Y     int32
	Color uint32
}

var r = regexp.MustCompile(`(\d+)\s*(\d+)\s*(\w+)`)

var colors = map[string]uint32{
	"red": 0xffff0000,
}

// FromString parses a string and returns a Pixel
func FromString(s string) (Pixel, error) {
	subMatches := r.FindStringSubmatch(s)
	if len(subMatches) != 4 {
		return Pixel{}, fmt.Errorf("did not match pattern")
	}
	xStr, yStr, colorStr := subMatches[1], subMatches[2], subMatches[3]
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return Pixel{}, err
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return Pixel{}, err
	}
	color, exists := colors[colorStr]
	if !exists {
		return Pixel{}, fmt.Errorf("not a valid color")
	}
	return Pixel{int32(x), int32(y), color}, nil
}
