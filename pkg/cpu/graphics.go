package cpu

import (
	"gogb/pkg/graphics"
)

func (c *CPU) Graphics(mCycles uint) {
	// TODO: is this right?
	interruptRequest := c.display.SetLCDStatus()
	if interruptRequest {
		c.RequestInterrupt(BIT_LCD)
	}
	if !c.display.LCDEnabled() {
		return
	}
	ly := c.ram.Ptr(graphics.ADDR_LY)
	c.display.Dots -= int(mCycles)
	if c.display.Dots <= 0 {
		*ly++
		currentLine := *ly
		c.display.Dots = 456
		// VBlank
		switch {
		case currentLine == 144:
			c.RequestInterrupt(BIT_VBLANK)
		case currentLine > 153:
			c.display.Swap()
		case currentLine < 144:
			c.display.DrawScanline(currentLine)
		}
	}
}
