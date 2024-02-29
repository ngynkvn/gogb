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

func (c *CPU) Interrupts() {
	if !c.IME && !c.Halt {
		return
	}

	flags := c.ram.ReadU8(ADDR_IF) | 0xE0
	enabled := c.ram.ReadU8(ADDR_IE)

	for i := uint8(0); i < 5; i++ {
		if bits.Test(flags, i) && bits.Test(enabled, i) {
			c.SpinCycle(2)
			c.ServiceInterrupt(i)
		}
	}

}

const (
	BIT_VBLANK = 0
	BIT_LCD    = 1
	BIT_TIMER  = 2
	BIT_SERIAL = 3
	BIT_JOYPAD = 4
)

func (c *CPU) RequestInterrupt(intAddr uint8) {
	flags := c.ram.ReadU8(ADDR_IF) | 0xE0
	flags = bits.Set(flags, intAddr)
	c.ram.WriteU8(ADDR_IF, flags)
}

func (c *CPU) ServiceInterrupt(intAddr uint8) {
	if !c.IME && c.Halt {
		c.Halt = false
		return
	}
	c.IME = false
	c.Halt = false
	flags := c.ram.ReadU8(ADDR_IF)
	flags = bits.Reset(flags, intAddr)
	c.ram.WriteU8(ADDR_IF, flags)

	c.PushStack(c.PC)
	c.PC = INT_ADDR_MAP[intAddr]
	c.SpinCycle(1)
}
