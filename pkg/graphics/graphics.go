package graphics

import (
	"gogb/pkg/bits"
	"gogb/pkg/mem"
	"image"
	"image/color"
	"image/png"
	"log"
	"log/slog"
	"os"
)

const (
	ADDR_LCDC    = 0xFF40
	ADDR_STAT    = 0xFF41
	ADDR_LY      = 0xFF44
	ADDR_LYC     = 0xFF45
	ADDR_DMA     = 0xFF46
	ADDR_SCROLLY = 0xFF42
	ADDR_SCROLLX = 0xFF43
	ADDR_WINDOWY = 0xFF4A
	ADDR_WINDOWX = 0xFF4B
)

const (
	LCDMode2Bound = 456 - 80
	LCDMode3Bound = LCDMode2Bound - 172
)

const (
	SCREEN_W = 160
	SCREEN_H = 144
)

type Surface [SCREEN_H][SCREEN_W]color.RGBA

type Display struct {
	ram        *mem.RAM
	Frame      image.Image
	screenData *image.RGBA
	// "Dots" are a time unit that represents one pixel push to the screen.
	// TODO: currently we render dots by scanline, change to per dot for more accuracy?
	Dots int
}

func NewDisplay(ram *mem.RAM) *Display {
	screenData := image.NewRGBA(image.Rect(0, 0, SCREEN_W, SCREEN_H))
	return &Display{
		ram: ram,
		// TODO: proper interface
		Frame:      screenData,
		screenData: screenData,
		Dots:       456,
	}
}

func (d *Display) STAT() uint8 {
	return d.ram.ReadU8(ADDR_STAT)
}

func (d *Display) LCDC() uint8 {
	return d.ram.ReadU8(ADDR_LCDC)
}

func (d *Display) LCDEnabled() bool {
	return bits.Test(d.LCDC(), 7)
}

// Test Bit 0
func (d *Display) BGWindowEnabled() bool {
	return bits.Test(d.LCDC(), 0)
}

// Bit 5
func (d *Display) WindowEnabled() bool {
	return bits.Test(d.LCDC(), 5)
}

// Test Bit 1
func (d *Display) ObjEnabled() bool {
	return bits.Test(d.LCDC(), 1)
}

// Bit 6
func (d *Display) WindowTileMapArea() []uint8 {
	if bits.Test(d.LCDC(), 6) {
		// 1 = $9C00-$9FFF
		return d.ram.Slice(0x9C00, 0x9FFF)
	} else {
		// 0 = $9800-$9BFF
		return d.ram.Slice(0x9800, 0x9BFF)
	}
}

// Bit 4
func (d *Display) UnsignedAddressMode() bool {
	return bits.Test(d.LCDC(), 4)
}

// Bit 2
func (d *Display) ObjHeight() uint8 {
	if bits.Test(d.LCDC(), 2) {
		return 16
	} else {
		return 8
	}
}

// Bit 3
func (d *Display) BGTileMapArea() []uint8 {
	if bits.Test(d.LCDC(), 3) {
		// 1 = $9C00-$9FFF
		return d.ram.Slice(0x9C00, 0x9FFF)
	} else {
		// 0 = $9800-$9BFF
		return d.ram.Slice(0x9800, 0x9BFF)
	}
}

func (d *Display) SetLCDStatus() bool {
	status := d.ram.Ptr(ADDR_STAT)
	ly := d.ram.Ptr(ADDR_LY)

	// TODO(006): check this routine
	if !d.LCDEnabled() {
		// d.ClearScreen()
		d.Dots = 456
		*ly = 0
		*status &= 0b1111_1100
		return false
	}

	currentLine := *ly
	currentMode := *status & 0b11

	var mode uint8
	requestInterrupt := false
	switch {
	case currentLine >= 144:
		mode = 1
		*status = bits.Set(*status, 0)
		*status = bits.Reset(*status, 1)
		requestInterrupt = bits.Test(*status, 4)
	case d.Dots >= LCDMode2Bound:
		mode = 2
		*status = bits.Set(*status, 1)
		*status = bits.Reset(*status, 0)
		requestInterrupt = bits.Test(*status, 5)
	case d.Dots >= LCDMode3Bound:
		mode = 3
		*status = bits.Set(*status, 1)
		*status = bits.Set(*status, 0)
		// This seems to be the wrong place to draw the scanline,
	default:
		mode = 0
		*status = bits.Reset(*status, 1)
		*status = bits.Reset(*status, 0)
		requestInterrupt = bits.Test(*status, 3)
	}

	if currentLine == d.ram.ReadU8(ADDR_LYC) {
		*status = bits.Set(*status, 2)
		requestInterrupt = bits.Test(*status, 6)
	} else {
		*status = bits.Reset(*status, 2)
	}
	return requestInterrupt && mode != currentMode
}

func (d *Display) DrawScanline(scanline uint8) {
	if d.BGWindowEnabled() {
		d.RenderTiles(scanline)
	}
	if d.ObjEnabled() {
		d.RenderSprites(int32(scanline))
	}
}

const TileSize = 16

// 00: White
// 01: Light Grey
// 10: Dark Grey
// 11: Black

func (d *Display) RenderTiles(scanline uint8) {
	scrollY := d.ram.ReadU8(ADDR_SCROLLY)
	scrollX := d.ram.ReadU8(ADDR_SCROLLX)
	windowY := d.ram.ReadU8(ADDR_WINDOWY)
	windowX := d.ram.ReadU8(ADDR_WINDOWX) - 7
	ly := d.ram.ReadU8(ADDR_LY)

	unsignedAddrMode := d.UnsignedAddressMode()
	usingWindow := d.WindowEnabled() && windowY <= ly

	// TODO: verify
	var bgMemory []uint8
	if usingWindow {
		bgMemory = d.WindowTileMapArea()
	} else {
		bgMemory = d.BGTileMapArea()
	}

	var baseTileAddr uint16
	if unsignedAddrMode {
		baseTileAddr = 0x8000
	} else {
		baseTileAddr = 0x8800
	}

	yPos := uint8(0)

	if usingWindow {
		yPos = ly - windowY
	} else {
		yPos = scrollY + ly
	}

	tileRow := uint16(yPos/8) * 32

	for p := uint8(0); p < 160; p++ {
		xPos := p + scrollX
		if usingWindow && p >= windowX {
			xPos = p - windowX
		}
		tileCol := uint16(xPos / 8)
		tileAddr := tileRow + tileCol
		tileNum := bgMemory[tileAddr]
		tileLocation := baseTileAddr
		if unsignedAddrMode {
			tileLocation += uint16(tileNum) * 16
		} else {
			tileLocation += (uint16(int8(tileNum)) + 128) * 16
		}
		line := uint16((yPos % 8) * 2)
		d1 := d.ram.ReadU8(tileLocation + line)
		d2 := d.ram.ReadU8(tileLocation + line + 1)

		// TODO: check this
		colorBit := xPos % 8
		colorBit -= 8
		colorBit = ^colorBit

		colorIndex := (bits.B(bits.Test(d2, colorBit))<<1 | bits.B(bits.Test(d1, colorBit)))
		// TODO: palette
		color := RGB_COLORS[d.GetColor(ColorIndex(colorIndex), 0xFF47)]

		// Safety check.
		if ly > 143 || p > 159 {
			slog.Error("OOB write attempted", "ly", ly, "p", p)
			continue
		}
		d.BlitPixel(color, p, ly)
	}
}

const MAX_SPRITES = 40

func (d *Display) RenderSprites(scanline int32) {
	for i := 0; i < MAX_SPRITES; i++ {
		idx := uint16(i * 4)
		// TODO: refactor constants
		yPos := d.ram.ReadU8(0xFE00+idx) - 16
		xPos := d.ram.ReadU8(0xFE00+idx+1) - 8
		tileLocation := d.ram.ReadU8(0xFE00 + idx + 2)
		attributes := d.ram.ReadU8(0xFE00 + idx + 3)

		yFlip := bits.Test(attributes, 6)
		xFlip := bits.Test(attributes, 5)

		scanline := d.ram.ReadU8(ADDR_LY)
		spriteHeight := d.ObjHeight()
		if scanline >= yPos && scanline < (yPos+spriteHeight) {
			// TODO: do these ops in u8
			line := int(scanline) - int(yPos)
			if yFlip {
				line -= int(spriteHeight)
				line *= -1
			}
			line *= 2
			dataAddr := (0x8000 + uint16(tileLocation)*16) + uint16(line)
			d1 := d.ram.ReadU8(dataAddr)
			d2 := d.ram.ReadU8(dataAddr + 1)
			for p := 7; p >= 0; p-- {
				// TODO: do these ops in u8
				colorBit := p
				if xFlip {
					colorBit -= 7
					colorBit *= -1
				}
				colorIndex := bits.B(bits.Test(d2, uint8(colorBit))) << 1
				colorIndex |= bits.B(bits.Test(d1, uint8(colorBit)))

				// TODO: const
				colorAddr := uint16(0xFF48)
				if bits.Test(attributes, 4) {
					colorAddr = 0xFF49
				}

				color := RGB_COLORS[d.GetColor(ColorIndex(colorIndex), colorAddr)]

				// Transparent for sprites
				if colorIndex == 0 {
					continue
				}
				xPix := 0 - p
				xPix += 7
				pixel := xPos + uint8(xPix)
				if scanline > 143 || p > 159 {
					slog.Error("OOB write attempted", "ly", scanline, "p", p)
					continue
				}
				d.BlitPixel(color, pixel, scanline)
			}
		}
	}
}

type ColorIndex uint8

const (
	WHITE      = ColorIndex(0)
	LIGHT_GRAY = ColorIndex(1)
	DARK_GRAY  = ColorIndex(2)
	BLACK      = ColorIndex(3)
)

var RGB_COLORS = [4]color.RGBA{
	{0xE0, 0xF8, 0xD0, 0xFF}, // WHITE
	{0x88, 0xC0, 0x70, 0xFF}, // LIGHT_GRAY
	{0x34, 0x68, 0x56, 0xFF}, // DARK_GRAY
	{0x08, 0x18, 0x20, 0xFF}, //BLACK
}

func (d *Display) BlitPixel(color color.RGBA, x uint8, y uint8) {
	d.screenData.Set(int(x), int(y), color)
}

func (d *Display) GetColor(index ColorIndex, addr uint16) ColorIndex {
	palette := d.ram.ReadU8(addr)
	hi, lo := uint8(0), uint8(0)
	switch index {
	case WHITE:
		hi, lo = 1, 0
	case LIGHT_GRAY:
		hi, lo = 3, 2
	case DARK_GRAY:
		hi, lo = 5, 4
	case BLACK:
		hi, lo = 7, 6
	}
	colorIndex := (bits.B(bits.Test(palette, hi))<<1 | bits.B(bits.Test(palette, lo)))
	return ColorIndex(colorIndex)
}

// TODO
// Swap and Clear the framebuffer
func (d *Display) Swap() {
	d.Frame = d.screenData
	// d.screenData = Surface{}
}

func (d *Display) DumpPNG(path string) {
	fp, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	img := d.Frame
	err = png.Encode(fp, img)
	if err != nil {
		log.Fatal(err)
	}
	fp.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
