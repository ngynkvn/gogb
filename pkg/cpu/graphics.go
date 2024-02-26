package cpu

func (c *CPU) Graphics(mCycles uint) {
	c.display.SetLCDStatus()
	if !c.display.LCDEnabled() {
		return
	}
	c.display.ScanlineCounter -= int(mCycles)
	if c.display.ScanlineCounter <= 0 {

	}
}
