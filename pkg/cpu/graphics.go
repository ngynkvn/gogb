package cpu

import (
	"gogb/pkg/graphics"
)

func (c *CPU) Graphics(mCycles uint) {
	// TODO: is this right?
	interruptRequest := c.display.SetLCDStatus()
	if interruptRequest {
		c.RequestInterrupt(0b00)
	}
	if !c.display.LCDEnabled() {
		return
	}
	ly := c.ram.Ptr(graphics.ADDR_LY)
	c.display.ScanlineCounter -= int(mCycles)
	if c.display.ScanlineCounter <= 0 {
		*ly++
		currentLine := *ly
		c.display.ScanlineCounter = 456
		// VBlank
		switch {
		case currentLine == 144 && !interruptRequest:
			c.RequestInterrupt(0b00)
		case currentLine > 153:
			c.display.Swap()
		case currentLine < 144:
			c.display.DrawScanline(currentLine)
		}
	}
}
