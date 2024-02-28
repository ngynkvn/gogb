package cpu

import "gogb/pkg/bits"

const ADDR_JOYPAD = 0xFF00

const (
	BIT_BUTTON = 5
	BIT_DPAD   = 4
	// The meaning of bits down here are determined by above.
	BIT_START_DOWN = 3
	BIT_SELECT_UP  = 2
	BIT_B_LEFT     = 1
	BIT_A_RIGHT    = 0
)

// TODO: Refactor this, copied from codeslinger
func (c *CPU) GetJoypadState() uint8 {
	res := ^(*c.ram.Ptr(ADDR_JOYPAD))

	// Check for buttons
	if !bits.Test(res, BIT_DPAD) {
		topJoypad := (c.JoypadState >> 4) | 0xF0
		res &= topJoypad
	} else if !bits.Test(res, BIT_BUTTON) {
		bottomJoypad := (c.JoypadState & 0x0F) | 0xF0
		res &= bottomJoypad
	}
	return res
}

func (c *CPU) KeyPressed(key uint8) {
	prevUnset := !bits.Test(c.JoypadState, key)
	c.JoypadState = bits.Reset(c.JoypadState, key)

	buttonPress := key > 3

	keyReq := *c.ram.Ptr(0xFF00)

	requestInterrupt :=
		(buttonPress && !bits.Test(keyReq, BIT_BUTTON)) ||
			(!buttonPress && !bits.Test(keyReq, BIT_DPAD))

	if requestInterrupt && prevUnset {
		c.RequestInterrupt(0b0100)
	}
}

func (c *CPU) KeyReleased(key uint8) {
	c.JoypadState = bits.Set(c.JoypadState, key)
}
