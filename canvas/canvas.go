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
func Run(pixels chan pixels.Pixel) {
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
			surface.FillRect(&sdl.Rect{X: p.X, Y: windowHeight - p.Y, W: 1, H: 1}, p.Color)
			window.UpdateSurface()
		case <-time.After(time.Duration(1000.0) * time.Millisecond / fps):
		}
	}
}
