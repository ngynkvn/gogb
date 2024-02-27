package render

import (
	"fmt"
	"gogb/pkg/cpu"
	"gogb/pkg/graphics"
	"gogb/pkg/mem"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Ebiten struct {
	cpu      *cpu.CPU
	display  *graphics.Display
	ram      *mem.RAM
	lastDraw time.Time
}

func (e *Ebiten) Draw(screen *ebiten.Image) {
	img := ebiten.NewImageFromImage(e.display.Frame)
	screen.DrawImage(img, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS:%f\nTPS:%f", ebiten.ActualFPS(), ebiten.ActualTPS()))
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
		time.Now(),
	}
}

const TICK_RATE = 1048576 * 4

func (e *Ebiten) Start() {
	ebiten.SetTPS(TICK_RATE)
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}
