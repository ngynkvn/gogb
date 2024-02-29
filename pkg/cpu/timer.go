package cpu

import "gogb/pkg/bits"

const (
	ADDR_DIV  = 0xFF04
	ADDR_TIMA = 0xFF05
	ADDR_TMA  = 0xFF06
	ADDR_TAC  = 0xFF07
)

func (c *CPU) DivRegister(mCycles uint) {
	c.DIV += uint(mCycles)
	if c.DIV >= 0xFF {
		c.DIV -= 0xFF
		div := c.ram.Ptr(ADDR_DIV)
		*div++
	}
}

// TODO(005): update with obscure behaviors https://gbdev.io/pandocs/Timer_Obscure_Behaviour.html
func (c *CPU) Timer(mCycles uint) {
	c.DivRegister(mCycles)
	if c.ClockEnabled() {
		c.TimerCount += mCycles
		freq := c.ClockFrequency()
		tima := c.ram.Ptr(ADDR_TIMA)
		for c.TimerCount >= freq {
			c.TimerCount -= freq
			if *tima == 0xFF {
				*tima = c.ram.ReadU8(ADDR_TMA)
				c.RequestInterrupt(BIT_TIMER)
			} else {
				*tima++
			}
		}
	}
}

func (c *CPU) ClockEnabled() bool {
	return bits.Test(c.ram.ReadU8(ADDR_TAC), 2)
}

func (c *CPU) ClockFrequency() uint {
	switch c.ram.ReadU8(ADDR_TAC) & 0b11 {
	case 0b00:
		return 256
	case 0b01:
		return 4
	case 0b10:
		return 16
	case 0b11:
		return 64
	}
	panic("unknown")
}
