package render

import (
	"gogb/pkg/cpu"
	"gogb/pkg/graphics"
	"gogb/pkg/mem"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ebiten struct {
	cpu      *cpu.CPU
	display  *graphics.Display
	ram      *mem.RAM
	buffer   []byte
	lastDraw time.Time
}

func (e *Ebiten) Draw(screen *ebiten.Image) {
	if time.Since(e.lastDraw) < time.Millisecond*17 {
		return
	}
	for yi, ys := range e.display.Frame {
		for xi, xs := range ys {
			location := yi*4*graphics.SCREEN_W + xi*4
			copy(e.buffer[location:], xs[:])
		}
	}
	screen.WritePixels(e.buffer)
	e.lastDraw = time.Now()
}

func (e *Ebiten) Update() error {
	e.cpu.Update()
	return nil
}

func (e *Ebiten) Layout(int, int) (int, int) {
	return int(graphics.SCREEN_W), int(graphics.SCREEN_H)
}

func NewEbiten(cpu *cpu.CPU, display *graphics.Display, ram *mem.RAM) *Ebiten {
	return &Ebiten{
		cpu,
		display,
		ram,
		make([]byte, graphics.SCREEN_W*graphics.SCREEN_H*4),
		time.Now(),
	}
}

const TICK_RATE = 1048576

func (e *Ebiten) Start() {
	ebiten.SetTPS(TICK_RATE)
	ebiten.SetScreenClearedEveryFrame(false)
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}
