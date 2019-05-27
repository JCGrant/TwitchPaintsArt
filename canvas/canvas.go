package canvas

import (
	"time"

	"github.com/JCGrant/twitch-paints/pixels"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 1000
	windowHeight = 1000
	appName      = "Twitch Paints"
	fps          = 60
)

// Run will instantiate the SDL2 window and run the rendering loop
func Run(pixels chan pixels.Pixel, windowWidth int32, windowHeight int32, initialPixels []pixels.Pixel) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(appName, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0xffffffff)
	for _, p := range initialPixels {
		drawPixel(surface, p.X, windowHeight-p.Y, p.Color)
	}
	window.UpdateSurface()

Loop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break Loop
			}
		}

		select {
		case p := <-pixels:
			drawPixel(surface, p.X, windowHeight-p.Y, p.Color)
			window.UpdateSurface()
		case <-time.After(time.Duration(1000.0) * time.Millisecond / fps):
		}
	}
}

func drawPixel(surface *sdl.Surface, x int32, y int32, color uint32) {
	surface.FillRect(&sdl.Rect{X: x, Y: y, W: 1, H: 1}, 0xff000000+color)
}
