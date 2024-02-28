package render

import (
	"fmt"
	"gogb/pkg/cpu"
	"gogb/pkg/graphics"
	"gogb/pkg/mem"
	"log"
	"log/slog"
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

const (
	RIGHT = 0
	LEFT  = 1
	UP    = 2
	DOWN  = 3

	A      = 4
	B      = 5
	START  = 6
	SELECT = 7
)

// TODO: accurate throttling
// TODO: fix janky input
func (e *Ebiten) Update() error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		slog.With(
			"cpu", e.cpu.Dump(),
			"0xFF47", fmt.Sprintf("%#02x", e.ram.ReadU8(0xFF47)),
			"0xFF48", fmt.Sprintf("%#02x", e.ram.ReadU8(0xFF48)),
			"0xFF49", fmt.Sprintf("%#02x", e.ram.ReadU8(0xFF49)),
		).Info("debug key pressed")
	case inpututil.IsKeyJustPressed(ebiten.KeyW):
		e.cpu.KeyPressed(UP)
	case inpututil.IsKeyJustPressed(ebiten.KeyA):
		e.cpu.KeyPressed(LEFT)
	case inpututil.IsKeyJustPressed(ebiten.KeyS):
		e.cpu.KeyPressed(DOWN)
	case inpututil.IsKeyJustPressed(ebiten.KeyD):
		e.cpu.KeyPressed(RIGHT)
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		e.cpu.KeyPressed(START)
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		e.cpu.KeyPressed(SELECT)
	case inpututil.IsKeyJustPressed(ebiten.KeyJ):
		e.cpu.KeyPressed(A)
	case inpututil.IsKeyJustPressed(ebiten.KeyK):
		e.cpu.KeyPressed(B)
	case inpututil.IsKeyJustReleased(ebiten.KeyW):
		e.cpu.KeyReleased(UP)
	case inpututil.IsKeyJustReleased(ebiten.KeyA):
		e.cpu.KeyReleased(LEFT)
	case inpututil.IsKeyJustReleased(ebiten.KeyS):
		e.cpu.KeyReleased(DOWN)
	case inpututil.IsKeyJustReleased(ebiten.KeyD):
		e.cpu.KeyReleased(RIGHT)
	case inpututil.IsKeyJustReleased(ebiten.KeySpace):
		e.cpu.KeyReleased(START)
	case inpututil.IsKeyJustReleased(ebiten.KeyEnter):
		e.cpu.KeyReleased(SELECT)
	case inpututil.IsKeyJustReleased(ebiten.KeyJ):
		e.cpu.KeyReleased(A)
	case inpututil.IsKeyJustReleased(ebiten.KeyK):
		e.cpu.KeyReleased(B)
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
