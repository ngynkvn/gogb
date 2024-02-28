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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func printHex(i uint8) {
	fmt.Printf("%#02x\n", i)
}

const TICK_RATE = 1048576 // hz

// TODO: accurate throttling
func (e *Ebiten) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		println(e.cpu.Dump())
		printHex(e.ram.ReadU8(0xFF47))
		printHex(e.ram.ReadU8(0xFF48))
		printHex(e.ram.ReadU8(0xFF49))
	}
	for cycs := uint(0); cycs < 17476/4; cycs += e.cpu.Update() {
	}
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

func (e *Ebiten) Start() {
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}
