package render

import (
	_ "embed"
	"fmt"
	"gogb/pkg/cpu"
	"gogb/pkg/graphics"
	"gogb/pkg/mem"
	"log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ebiten struct {
	cpu     *cpu.CPU
	display *graphics.Display
	ram     *mem.RAM
}

var shader *ebiten.Shader

func (e *Ebiten) Draw(screen *ebiten.Image) {
	img := ebiten.NewImageFromImage(e.display.Frame)
	screen.DrawImage(img, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS:%f\nTPS:%f\n", ebiten.ActualFPS(), ebiten.ActualTPS()))
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

var KeyMap = map[ebiten.Key]uint8{
	ebiten.KeyW: UP,
	ebiten.KeyA: LEFT,
	ebiten.KeyS: DOWN,
	ebiten.KeyD: RIGHT,
	ebiten.KeyJ: A,
	ebiten.KeyK: B,

	ebiten.KeyEnter: SELECT,
	ebiten.KeySpace: START,
}

// TODO: accurate throttling
// TODO: fix janky input
func (e *Ebiten) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		slog.With(
			"cpu", e.cpu.Dump(),
			"0xFF47", fmt.Sprintf("%#02x", e.ram.ReadU8(0xFF47)),
			"0xFF48", fmt.Sprintf("%#02x", e.ram.ReadU8(0xFF48)),
			"0xFF49", fmt.Sprintf("%#02x", e.ram.ReadU8(0xFF49)),
		).Info("debug key pressed")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	for k, v := range KeyMap {
		if ebiten.IsKeyPressed(k) {
			e.cpu.KeyPressed(v)
		} else {
			e.cpu.KeyReleased(v)
		}
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
	}
}

// -go:embed shaders/crt.kage
// var shaderSrc []byte

func (e *Ebiten) Start() {
	// var err error
	// shader, err = ebiten.NewShader(shaderSrc)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}
