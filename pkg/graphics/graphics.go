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
	ram             *mem.RAM
	Frame           image.Image
	screenData      *image.RGBA
	ScanlineCounter int
}

func NewDisplay(ram *mem.RAM) *Display {
	screenData := image.NewRGBA(image.Rect(0, 0, SCREEN_W, SCREEN_H))
	return &Display{
		ram: ram,
		// TODO: proper interface
		Frame:           screenData,
		screenData:      screenData,
		ScanlineCounter: 456,
	}
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
		d.ScanlineCounter = 456
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
	case d.ScanlineCounter >= LCDMode2Bound:
		mode = 2
		*status = bits.Set(*status, 1)
		*status = bits.Reset(*status, 0)
		requestInterrupt = bits.Test(*status, 5)
	case d.ScanlineCounter >= LCDMode3Bound:
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

	// windowArea := d.WindowTileMapArea()
	bgArea := d.BGTileMapArea()
	unsignedAddrMode := d.UnsignedAddressMode()
	usingWindow := d.WindowEnabled() && windowY <= ly

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
		tileNum := bgArea[tileAddr]
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

		colorNum := (bits.B(bits.Test(d2, colorBit))<<1 | bits.B(bits.Test(d1, colorBit)))
		// TODO: palette
		color := d.GetColor(uint8(colorNum), 0xFF47)

		// Safety check.
		if ly > 143 || p > 159 {
			slog.Error("OOB write attempted", "ly", ly, "p", p)
			continue
		}
		d.BlitPixel(RGB_COLORS[color], p, ly)
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
				colorNum := bits.B(bits.Test(d2, uint8(colorBit))) << 1
				colorNum |= bits.B(bits.Test(d1, uint8(colorBit)))

				// TODO: const
				colorAddr := uint16(0xFF48)
				if bits.Test(attributes, 4) {
					colorAddr = 0xFF49
				}

				colorValue := RGB_COLORS[d.GetColor(uint8(colorNum), colorAddr)]

				// Transparent for sprites
				if colorValue == RGB_COLORS[WHITE] {
					continue
				}
				xPix := 0 - p
				xPix += 7
				pixel := xPos + uint8(xPix)
				if scanline > 143 || p > 159 {
					slog.Error("OOB write attempted", "ly", scanline, "p", p)
					continue
				}
				d.BlitPixel(colorValue, pixel, scanline)
			}
		}
	}
}

type Color uint8

const (
	WHITE      = Color(0)
	LIGHT_GRAY = Color(1)
	DARK_GRAY  = Color(2)
	BLACK      = Color(3)
)

var RGB_COLORS = [4]color.RGBA{
	{0xFF, 0xFF, 0xFF, 0xFF}, //WHITE
	{0xCC, 0xCC, 0xCC, 0xFF}, // LIGHT_GRAY
	{0x77, 0x77, 0x77, 0xFF}, // DARK_GRAY
	{0x00, 0x00, 0x00, 0xFF}, // BLACK
}

func (d *Display) BlitPixel(color color.RGBA, x uint8, y uint8) {
	d.screenData.Set(int(x), int(y), color)
}

func (d *Display) GetColor(colorNum uint8, addr uint16) Color {
	palette := d.ram.ReadU8(addr)
	hi, lo := uint8(0), uint8(0)
	switch colorNum {
	case 0:
		hi, lo = 1, 0
	case 1:
		hi, lo = 3, 2
	case 2:
		hi, lo = 5, 4
	case 3:
		hi, lo = 7, 6
	}
	colorVal := (bits.B(bits.Test(palette, hi))<<1 | bits.B(bits.Test(palette, lo)))
	return Color(colorVal)
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
