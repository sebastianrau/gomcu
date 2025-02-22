package gomcu

// A Channel is used to define which Fader, Meter or VPot should be modified.
type Channel byte

const (
	Channel1 Channel = iota
	Channel2
	Channel3
	Channel4
	Channel5
	Channel6
	Channel7
	Channel8
	// Master is only a fader and will do nothing if used to set a VPot or a Meter.
	Master

	LenChannels = 9

	FaderMax = 16382
	FaderMin = 0
)

var (
	ChannelNames = []string{
		"Channel1",
		"Channel2",
		"Channel3",
		"Channel4",
		"Channel5",
		"Channel6",
		"Channel7",
		"Channel8",
		"Master",
	}

	ChannelIDtoChannel = map[int]Channel{
		1: 0x10,
		2: 0x11,
		3: 0x12,
		4: 0x13,
		5: 0x14,
		6: 0x15,
		7: 0x16,
		8: 0x17,
		9: 0x18,
	}

	ChannelNamesToChannel = map[string]Channel{
		"Channel1": 0x10,
		"Channel2": 0x11,
		"Channel3": 0x12,
		"Channel4": 0x13,
		"Channel5": 0x14,
		"Channel6": 0x15,
		"Channel7": 0x16,
		"Channel8": 0x17,
		"Master":   0x18,
	}
)
