package cpu

import "gogb/pkg/bits"

const (
	INT_VBLANK = 0x0040
	INT_STAT   = 0x0048
	INT_TIMER  = 0x0050
	INT_SERIAL = 0x0058
	INT_JOYPAD = 0x0060
)

var INT_ADDR_MAP = [5]uint16{
	INT_VBLANK,
	INT_STAT,
	INT_TIMER,
	INT_SERIAL,
	INT_JOYPAD,
}

const ADDR_IF = 0xFF0F
const ADDR_IE = 0xFFFF

func (c *CPU) Interrupts() (cycles uint) {
	if !c.IME && !c.halt {
		return 0
	}

	flags := c.ram.ReadU8(ADDR_IF) | 0xE0
	enabled := c.ram.ReadU8(ADDR_IE)

	for i := uint8(0); i < 5; i++ {
		if bits.Test(flags, i) && bits.Test(enabled, i) {
			c.ServiceInterrupt(i)
			return 20
		}
	}
	return 0

}

func (c *CPU) ServiceInterrupt(IntAddr uint8) {
	if !c.IME && c.halt {
		c.halt = false
		return
	}
	c.IME = false
	c.halt = false
	flags := c.ram.ReadU8(ADDR_IF)
	flags = bits.Reset(flags, IntAddr)
	c.ram.WriteU8(ADDR_IF, flags)

	c.PushStack(c.PC)
	c.PC = INT_ADDR_MAP[IntAddr]
}

func (c *CPU) Timer(cycles uint) {

}
