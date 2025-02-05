package gomcu

import (
	"strings"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var (
	header    = []byte{0x00, 0x00, 0x66, 0x14}
	header_xt = []byte{0x00, 0x00, 0x66, 0x15}
	header_c4 = []byte{0x00, 0x00, 0x66, 0x17}

	SysExMessages = map[string][]byte{
		"Query":       {0x00},
		"GoOffline":   {0x0F, 0x7F},
		"Version":     {0x13, 0x00},
		"ResetFaders": {0x61},
		"ResetLEDs":   {0x62},
		"Reset":       {0x63},
	}
)

// Reset resets and quickly triggers all the available features on the control surface.
// I recommend running this both to avoid some errors and as a sanity check to make sure that your entire control surface is working.
func Reset(output drivers.Out) error {

	send, err := midi.SendTo(output)
	if err != nil {
		return err
	}
	var m []midi.Message
	for i := 0; i < LenIDs; i++ {
		m = append(m, SetLED(Switch(i), StateOn))
	}
	for i := 0; i < LenChannels; i++ {
		m = append(m, SetFaderPos(Channel(i), 0x1FFF))
	}
	for i := 0; i < LenChannels-1; i++ {
		m = append(m, SetVPot(Channel(i), VPotMode3, VPot6+VPotDot))
	}
	for i := 0; i < LenChannels-1; i++ {
		m = append(m, SetMeter(Channel(i), Clipping))
	}
	for i := 0; i < LenDigits; i++ {
		m = append(m, SetDigit(Digit(uint8(i)+0x40), Char0+DigitDot))
	}
	for i := 0; i < LenLines; i++ {
		m = append(m, SetLCD(i, " "))
	}

	for _, msg := range m {
		err := send(msg)
		if err != nil {
			return err
		}
	}

	time.Sleep(100 * time.Millisecond)

	m = []midi.Message{}

	for i := 0; i < LenIDs; i++ {
		m = append(m, SetLED(Switch(i), StateOff))
	}
	for i := 0; i < LenChannels; i++ {
		m = append(m, SetFaderPos(Channel(i), 0x0))
	}
	for i := 0; i < LenChannels-1; i++ {
		m = append(m, SetVPot(Channel(i), VPotMode0, VPot0))
	}
	for i := 0; i < LenChannels-1; i++ {
		m = append(m, SetMeter(Channel(i), ClipOff))
	}
	for i := 0; i < LenChannels-1; i++ {
		m = append(m, SetMeter(Channel(i), LessThan60))
	}
	for i := 0; i < LenDigits; i++ {
		m = append(m, SetDigit(Digit(uint8(i)+0x40), SymbolSpace))
	}
	for i := 0; i < LenLines; i++ {
		m = append(m, SetLCD(i, " "))
	}
	for _, msg := range m {
		err := send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetLED sets the chosen button's LED to the chosen State.
func SetLED(led Switch, state State) midi.Message {
	return midi.NoteOn(0, uint8(led), byte(state))
}

func SendOff(btn Switch) midi.Message {
	return midi.NoteOff(0, uint8(btn))
}

func SendOffVelocity(btn Switch, val byte) midi.Message {
	return midi.NoteOffVelocity(0, uint8(btn), val)
}

// SetFaderPos sets the position of the chosen fader to a number between 0 (bottom) and 16382 (top).
func SetFaderPos(fader Channel, pos uint16) midi.Message {
	p := int16(pos) - 8191
	return midi.Pitchbend(uint8(fader), p)
}

// SetTimeDisplay sets multiple characters on the timecode display.
// Note: letters is limited to ten characters and is right aligned.
// Refer to timecode.Digit for valid characters.
func SetTimeDisplay(letters string) (m []midi.Message) {
	bytes := []byte(strings.ToUpper(letters))
	if len(bytes) > 10 {
		bytes = bytes[:10]
	}

	for i, char := range bytes {
		if char >= 0x40 && char <= 0x60 {
			bytes[i] = char - 0x40
		}
	}

	for i := len(bytes)/2 - 1; i >= 0; i-- {
		opp := len(bytes) - 1 - i
		bytes[i], bytes[opp] = bytes[opp], bytes[i]
	}

	for i := uint8(0); int(i) < len(bytes); i++ {
		m = append(m, midi.ControlChange(15, i+0x40, bytes[i]))
	}
	return
}

// SetDigit sets an individual digit on the timecode or Assignment section.
// Refer to Char for more information on valid characters.
func SetDigit(digit Digit, char Char) midi.Message {
	if (char >= 0x40 && char <= 0x60) || (char >= 0x80 && char <= 0xA0) {
		char = char - 0x40
	}
	return midi.ControlChange(15, byte(digit), byte(char))
}

// SetLCD sets the text (an ASCII string) found on the LCD starting from the specified offset.
func SetLCDC4(offset int, row int, text string) midi.Message {
	rowu := 0x30 + uint8(row)
	return midi.SysEx(append(append(header_c4, rowu, uint8(offset)), []byte(text)...))
}

// SetLCD sets the text (an ASCII string) found on the LCD starting from the specified offset.
func SetLCDXT(offset int, text string) midi.Message {
	return midi.SysEx(append(append(header_xt, 0x12, uint8(offset)), []byte(text)...))
}

// SetLCD sets the text (an ASCII string) found on the LCD starting from the specified offset.
func SetLCD(offset int, text string) midi.Message {
	return midi.SysEx(append(append(header, 0x12, uint8(offset)), []byte(text)...))
}

// SetVPot sets the LEDs around the knobs (VPots).
// Refer to VPotMode for an explanation of the various VPot modes.
func SetVPot(ch Channel, mode VPotMode, led VPotLED) midi.Message {
	return midi.ControlChange(0, byte(ch+0x30), byte(mode)+byte(led))
}

// SetMeter sets the level meter for the selected Channel to the desired value.
func SetMeter(ch Channel, value MeterLevel) midi.Message {
	return midi.AfterTouch(0, byte(ch<<4)+byte(value))
}
