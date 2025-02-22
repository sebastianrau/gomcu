package gomcu

// A Channel is used to define which Fader, Meter or VPot should be modified.
type Channel byte

// TODO Check if Channel ID is used correctly
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

	ChannelIDs = map[string]Channel{
		"Channel1": Channel1,
		"Channel2": Channel2,
		"Channel3": Channel3,
		"Channel4": Channel4,
		"Channel5": Channel5,
		"Channel6": Channel6,
		"Channel7": Channel7,
		"Channel8": Channel8,
		"Master":   Master,
	}
)
