package canvas

import (
	"time"

	"github.com/JCGrant/twitch-paints/pixels"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	appName = "Twitch Paints"
	fps     = 60
)

type canvas struct {
	pixels        chan pixels.Pixel
	windowWidth   int
	windowHeight  int
	initialPixels []pixels.Pixel
}

func (c canvas) run() {
	cfg := pixelgl.WindowConfig{
		Title:  appName,
		Bounds: pixel.R(0, 0, float64(c.windowWidth), float64(c.windowHeight)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.White)
	imd := imdraw.New(nil)
	for _, p := range c.initialPixels {
		drawPixel(p, win, imd)
	}
	imd.Draw(win)

	for !win.Closed() {
		select {
		case p := <-c.pixels:
			drawPixel(p, win, imd)
		case <-time.After(time.Duration(1000.0) * time.Millisecond / fps):
		}
		imd.Draw(win)
		win.Update()
	}
}

// Run runs the canvas
func Run(pixels chan pixels.Pixel, windowWidth int, windowHeight int, initialPixels []pixels.Pixel) {
	c := canvas{pixels, windowWidth, windowHeight, initialPixels}
	pixelgl.Run(c.run)
}

func drawPixel(p pixels.Pixel, win *pixelgl.Window, imd *imdraw.IMDraw) {
	imd.Color = p.Color
	imd.Push(pixel.V(float64(p.X), float64(p.Y)), pixel.V(float64(p.X+1), float64(p.Y+1)))
	imd.Rectangle(0)
}
