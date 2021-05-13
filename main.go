// The MIT License (MIT)
//
// Copyright (c) 2015-2016 Martin Lindhe
// Copyright (c) 2016      Hajime Hoshi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// +build example

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// World represents the game state.
type World struct {
	iterations []uint8
	width      int
	height     int
	realMin    float64
	realMax    float64
	imagMin    float64
	imagMax    float64
	palette    [maxIteration + 1]byte
}

// NewWorld creates a new world.
func NewWorld(width, height int) *World {
	w := &World{
		iterations: make([]uint8, width*height),
		width:      width,
		height:     height,
		realMin:    -2,
		realMax:    1,
		imagMin:    -1,
		imagMax:    1,
	}

	for i := range w.palette {
		w.palette[i] = byte(math.Sqrt(float64(i)/float64(len(w.palette))) * 0xff)
		if debug {
			log.Println(w.palette[i])
		}
	}
	if debug {
		log.Printf("Size of palette is %d\n", len(w.palette))
	}
	w.Update()

	return w
}

// Update game state by one tick.
func (w *World) Update() {
	width := w.width
	height := w.height
	realMin := w.realMin
	realMax := w.realMax
	imagMin := w.imagMin
	imagMax := w.imagMax
	next := make([]uint8, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := complex(realMin+(float64(x)/float64(width))*(realMax-realMin), imagMin+(float64(y)/float64(height))*(imagMax-imagMin))
			next[y*width+x] = mandelbrot(c)
		}
	}
	w.iterations = next
}

func mandelbrot(c complex128) uint8 {
	z := complex(0, 0)
	var n uint8 = 0
	for complexModulusSquared(z) <= 4 && n < maxIteration {
		z = z*z + c
		n++
	}
	return n
}

func complexModulusSquared(c complex128) uint {
	return uint(real(c)*real(c) + imag(c)*imag(c))
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	for i := 0; i < len(w.iterations); i++ {
		colorOffset := w.palette[w.iterations[i]]
		pix[4*i] = 255 - colorOffset // This could be better ??
		pix[4*i+1] = 255 - colorOffset
		pix[4*i+2] = 255 - colorOffset
		pix[4*i+3] = 0xff
	}
}

const (
	screenWidth  = 640
	screenHeight = 480
	maxIteration = 30
	debug        = true
	maxTPS       = 120
)

type Game struct {
	world  *World
	pixels []byte
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.world.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)

	if debug {
		msg := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
		ebitenutil.DebugPrint(screen, msg)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		world: NewWorld(screenWidth, screenHeight),
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("goFrac")

	if debug {
		ebiten.SetMaxTPS(maxTPS)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
