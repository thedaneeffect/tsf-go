package tsf

//#define TML_STATIC 1
//#define TML_IMPLEMENTATION
//#include "tml.h"
import "C"

//https://www.midi.org/specifications

const (
	NoteOff         = 0x80
	NoteOn          = 0x90
	KeyPressure     = 0xA0
	ControlChange   = 0xB0
	ProgramChange   = 0xC0
	ChannelPressure = 0xD0
	PitchBend       = 0xE0
	SetTempo        = 0x51
)

var MessageTypeToString = map[int]string{
	NoteOff:         "NoteOff",
	NoteOn:          "NoteOn",
	KeyPressure:     "KeyPressure",
	ControlChange:   "ControlChange",
	ProgramChange:   "ProgramChange",
	ChannelPressure: "ChannelPressure",
	PitchBend:       "PitchBend",
	SetTempo:        "SetTempo",
}

var ControlChangeToString = map[int]string{
	BankSelectMSB:      "BankSelectMSB",
	ModulationMSB:      "ModulationMSB",
	BreathMSB:          "BreathMSB",
	FootMSB:            "FootMSB",
	PortamentoTimeMSB:  "PortamentoTimeMSB",
	DataEntryMSB:       "DataEntryMSB",
	VolumeMSB:          "VolumeMSB",
	BalanceMSB:         "BalanceMSB",
	PanMSB:             "PanMSB",
	ExpressionMSB:      "ExpressionMSB",
	Effects1MSB:        "Effects1MSB",
	Effects2MSB:        "Effects2MSB",
	GPC1MSB:            "GPC1MSB",
	GPC2MSB:            "GPC2MSB",
	GPC3MSB:            "GPC3MSB",
	GPC4MSB:            "GPC4MSB",
	BankSelectLSB:      "BankSelectLSB",
	ModulationWheelLSB: "ModulationWheelLSB",
	BreathLSB:          "BreathLSB",
	FootLSB:            "FootLSB",
	PortamentoTimeLSB:  "PortamentoTimeLSB",
	DataEntryLSB:       "DataEntryLSB",
	VolumeLSB:          "VolumeLSB",
	BalanceLSB:         "BalanceLSB",
	PanLSB:             "PanLSB",
	ExpressionLSB:      "ExpressionLSB",
	Effects1LSB:        "Effects1LSB",
	Effects2LSB:        "Effects2LSB",
	GPC1LSB:            "GPC1LSB",
	GPC2LSB:            "GPC2LSB",
	GPC3LSB:            "GPC3LSB",
	GPC4LSB:            "GPC4LSB",
	SustainSwitch:      "SustainSwitch",
	PortamentoSwitch:   "PortamentoSwitch",
	SostenutoSwitch:    "SostenutoSwitch",
	SoftPedalSwitch:    "SoftPedalSwitch",
	LegatoSwitch:       "LegatoSwitch",
	Hold2Switch:        "Hold2Switch",
	SoundCtrl1:         "SoundCtrl1",
	SoundCtrl2:         "SoundCtrl2",
	SoundCtrl3:         "SoundCtrl3",
	SoundCtrl4:         "SoundCtrl4",
	SoundCtrl5:         "SoundCtrl5",
	SoundCtrl6:         "SoundCtrl6",
	SoundCtrl7:         "SoundCtrl7",
	SoundCtrl8:         "SoundCtrl8",
	SoundCtrl9:         "SoundCtrl9",
	SoundCtrl10:        "SoundCtrl10",
	GPC5:               "GPC5",
	GPC6:               "GPC6",
	GPC7:               "GPC7",
	GPC8:               "GPC8",
	PortamentoCtrl:     "PortamentoCtrl",
	EffectsDepth1:      "EffectsDepth1",
	EffectsDepth2:      "EffectsDepth2",
	EffectsDepth3:      "EffectsDepth3",
	EffectsDepth4:      "EffectsDepth4",
	EffectsDepth5:      "EffectsDepth5",
	DataEntryIncr:      "DataEntryIncr",
	DataEntryDecr:      "DataEntryDecr",
	NRPNLSB:            "NRPNLSB",
	NRPNMSB:            "NRPNMSB",
	RPNLSB:             "RPNLSB",
	RPNMSB:             "RPNMSB",
	AllSoundOff:        "AllSoundOff",
	AllCtrlOff:         "AllCtrlOff",
	LocalControl:       "LocalControl",
	AllNotesOff:        "AllNotesOff",
	OmniOff:            "OmniOff",
	OmniOn:             "OmniOn",
	PolyOff:            "PolyOff",
	PolyOn:             "PolyOn",
}

const (
	BankSelectMSB      = 0x00
	ModulationMSB      = 0x01
	BreathMSB          = 0x02
	FootMSB            = 0x04
	PortamentoTimeMSB  = 0x05
	DataEntryMSB       = 0x06
	VolumeMSB          = 0x07
	BalanceMSB         = 0x08
	PanMSB             = 0x0A
	ExpressionMSB      = 0x0B
	Effects1MSB        = 0x0C
	Effects2MSB        = 0x0D
	GPC1MSB            = 0x10 /* general purpose controller */
	GPC2MSB            = 0x11
	GPC3MSB            = 0x12
	GPC4MSB            = 0x13
	BankSelectLSB      = 0x20
	ModulationWheelLSB = 0x21
	BreathLSB          = 0x22
	FootLSB            = 0x24
	PortamentoTimeLSB  = 0x25
	DataEntryLSB       = 0x26
	VolumeLSB          = 0x27
	BalanceLSB         = 0x28
	PanLSB             = 0x2A
	ExpressionLSB      = 0x2B
	Effects1LSB        = 0x2C
	Effects2LSB        = 0x2D
	GPC1LSB            = 0x30
	GPC2LSB            = 0x31
	GPC3LSB            = 0x32
	GPC4LSB            = 0x33
	SustainSwitch      = 0x40
	PortamentoSwitch   = 0x41
	SostenutoSwitch    = 0x42
	SoftPedalSwitch    = 0x43
	LegatoSwitch       = 0x44
	Hold2Switch        = 0x45
	SoundCtrl1         = 0x46 // default: sound variation
	SoundCtrl2         = 0x47 // default: timbre/harmonic intensity
	SoundCtrl3         = 0x48 // default: release time
	SoundCtrl4         = 0x49 // default: attack time
	SoundCtrl5         = 0x4A // default: brightness
	SoundCtrl6         = 0x4B // default: decay time
	SoundCtrl7         = 0x4C // default: vibrato rate
	SoundCtrl8         = 0x4D // default: vibrato depth
	SoundCtrl9         = 0x4E // default: vibrato delay
	SoundCtrl10        = 0x4F // default: undefined
	GPC5               = 0x50 // general purpose controller
	GPC6               = 0x51
	GPC7               = 0x52
	GPC8               = 0x53
	PortamentoCtrl     = 0x54
	EffectsDepth1      = 0x5B
	EffectsDepth2      = 0x5C
	EffectsDepth3      = 0x5D
	EffectsDepth4      = 0x5E
	EffectsDepth5      = 0x5F
	DataEntryIncr      = 0x60
	DataEntryDecr      = 0x61
	NRPNLSB            = 0x62 // non-registered parameter number
	NRPNMSB            = 0x63
	RPNLSB             = 0x64 // registered parameter number
	RPNMSB             = 0x65
	AllSoundOff        = 0x78
	AllCtrlOff         = 0x79 // reset all controlers
	LocalControl       = 0x7A
	AllNotesOff        = 0x7B
	OmniOff            = 0x7C
	OmniOn             = 0x7D
	PolyOff            = 0x7E
	PolyOn             = 0x7F
)

type Message struct {
	message *C.tml_message
}

func (m Message) Time() int            { return int(m.message.time) }
func (m Message) Type() int            { return int(m.message._type) }
func (m Message) Channel() int         { return int(m.message.channel) }
func (m Message) Key() int             { return int(m.message.anon0[0]) }
func (m Message) Control() int         { return int(m.message.anon0[0]) }
func (m Message) Program() int         { return int(m.message.anon0[0]) }
func (m Message) ChannelPressure() int { return int(m.message.anon0[0]) }
func (m Message) Velocity() int        { return int(m.message.anon0[1]) }
func (m Message) KeyPressure() int     { return int(m.message.anon0[1]) }
func (m Message) ControlValue() int    { return int(m.message.anon0[1]) }
func (m Message) PitchBend() int {
	return (int(m.message.anon0[0]) << 8) | int(m.message.anon0[1])
}
func (m Message) HasNext() bool { return m.message.next != nil }
func (m Message) Next() Message { return Message{m.message.next} }
func (m Message) IsNil() bool   { return m.message == nil }

func LoadMidiFile(filename string) Message {
	return Message{C.tml_load_filename(C.CString(filename))}
}

func LoadMidiMemory(mem []byte) Message {
	return Message{C.tml_load_memory(C.CBytes(mem), C.int(len(mem)))}
}
