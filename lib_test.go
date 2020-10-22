package tsf

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func TestRawOutput(t *testing.T) {
	msg := LoadMidiFile("expanse.mid")

	if msg.IsNil() {
		t.Error("bad midi")
	}

	font := LoadSoundFontFile("winxp.sf2")

	if font.IsNil() {
		t.Error("bad soundfont")
	}

	sampleRate := 44100
	channels := 2

	font.SetOutput(OutputModeStereoInterleaved, sampleRate, 0)

	// ensure channel 9 is using drum kit 0 (some midis don't set it for some reason)
	font.ChannelSetBank(9, 128)
	font.ChannelSetPresetNumber(9, 0, false)

	buffer := make([]int16, 2048)
	msec := 0.0

	wav, err := os.Create("out.raw")

	if err != nil {
		t.Error(err)
	}

	defer wav.Close()

	out, err := os.Create("out.txt")

	if err != nil {
		t.Error(err)
	}

	defer out.Close()

	step := 1000.0 / float64(sampleRate)

	for !msg.IsNil() {
		sampleBlock := 0
		sampleCount := len(buffer) / channels
		offset := 0

		for {
			sampleBlock = RenderBlockSize

			if sampleCount <= 0 {
				break
			}

			if sampleBlock > sampleCount {
				sampleBlock = sampleCount
			}

			for msec += float64(sampleBlock) * step; !msg.IsNil() && msec >= float64(msg.Time()); msg = msg.Next() {
				time := fmt.Sprintf("%dms", msg.Time())
				channel := fmt.Sprintf("[%d]", msg.Channel())
				_type := "nil"

				var output string

				if t, ok := MessageTypeToString[msg.Type()]; ok {
					_type = t
				} else {
					_type = fmt.Sprintf("unknown type(%d)", msg.Type())
				}

				switch msg.Type() {
				case ProgramChange:
					font.ChannelSetPresetNumber(msg.Channel(), msg.Program(), msg.Channel() == 9)
					output = fmt.Sprintf("Program=%d Drums=%v", msg.Program(), msg.Channel() == 9)
					break
				case NoteOn:
					font.ChannelNoteOn(msg.Channel(), msg.Key(), float32(msg.Velocity())/127.0)
					output = fmt.Sprintf("Key=%d Velocity=%d", msg.Key(), msg.Velocity())
					break
				case NoteOff:
					font.ChannelNoteOff(msg.Channel(), msg.Key())
					output = fmt.Sprintf("Key=%d", msg.Key())
					break
				case PitchBend:
					font.ChannelSetPitchWheel(msg.Channel(), msg.PitchBend())
					output = fmt.Sprintf("Value=%d", msg.PitchBend())
					break
				case ControlChange:
					font.ChannelMidiControl(msg.Channel(), msg.Control(), msg.ControlValue())

					if name, ok := ControlChangeToString[msg.Control()]; ok {
						output = fmt.Sprintf("%s=%d", name, msg.ControlValue())
					} else {
						output = fmt.Sprintf("control %d=%d", msg.Control(), msg.ControlValue())
					}
					break
				default:
					output = fmt.Sprintf("raw: %+v", msg.message.anon0)
					break
				}

				if _, err := out.Write(([]byte)(fmt.Sprintf("%s %s %s %s\n", time, channel, _type, output))); err != nil {
					t.Error(err)
				}
			}

			font.RenderShort(buffer[offset:offset+sampleBlock], sampleBlock, false)
			sampleCount -= sampleBlock
			offset += sampleBlock * channels
		}

		if err := binary.Write(wav, binary.BigEndian, buffer); err != nil {
			t.Error(err)
		}
	}
}
