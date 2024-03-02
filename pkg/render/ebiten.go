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
	bounds := screen.Bounds()
	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), shader, &ebiten.DrawRectShaderOptions{
		Images: [4]*ebiten.Image{img},
	})
	e.DebugPrint(screen)
}

func (e *Ebiten) DebugPrint(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`
tps:%f
LY:%#08b
LCDC:%#08b
STAT&3:%#08b
STAT:%#08b`,
		ebiten.ActualTPS(),
		e.ram.ReadU8(graphics.ADDR_LY),
		e.ram.ReadU8(graphics.ADDR_LCDC),
		e.display.STAT()&3,
		e.display.STAT(),
	))
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

func p(i uint8) string { return fmt.Sprintf("%#02x", i) }
func b(i uint8) string { return fmt.Sprintf("%#08b", i) }

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
		slog.With(
			"STAT", p(e.display.STAT()),
			"LCDC", p(e.display.LCDC()),
		).Info("display")
		slog.With(
			"IE", b(e.ram.ReadU8(cpu.ADDR_IE)),
			"IF", b(e.ram.ReadU8(cpu.ADDR_IF)),
		).Info("display")
		slog.Info("Halted?", "halt", e.cpu.Halt)
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

	for cycs := uint(0); cycs < 69905; cycs += e.cpu.Update() {
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

//go:embed shaders/crt.kage
var shaderSrc []byte

func (e *Ebiten) Start() {
	var err error
	ebiten.SetWindowSize(graphics.SCREEN_W*4, graphics.SCREEN_H*4)
	ebiten.SetWindowTitle("GB Emulator")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	shader, err = ebiten.NewShader(shaderSrc)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}
