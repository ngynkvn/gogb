package audio

const (
	// Master audio
	ADDR_NR52 = 0xFF26
	// Sound panning
	ADDR_NR51 = 0xFF25
	// Master volume & VIN panning
	ADDR_NR50 = 0xFF24
)

// ADDR_NR52
const (
	BIT_ENABLED    = 7
	BIT_CH4_ACTIVE = 3
	BIT_CH3_ACTIVE = 2
	BIT_CH2_ACTIVE = 1
	BIT_CH1_ACTIVE = 0
)

// ADDR_NR51
const (
	BIT_CH4_LEFT  = 7
	BIT_CH3_LEFT  = 6
	BIT_CH2_LEFT  = 5
	BIT_CH1_LEFT  = 4
	BIT_CH4_RIGHT = 3
	BIT_CH3_RIGHT = 2
	BIT_CH2_RIGHT = 1
	BIT_CH1_RIGHT = 0
)
