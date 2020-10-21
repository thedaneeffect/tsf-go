package tsf

//#cgo LDFLAGS: -lm
//#define TSF_STATIC 1
//#define TSF_IMPLEMENTATION 1
//#include "tsf.h"
import "C"

// The lower this block size is the more accurate the effects are.
// Increasing the value significantly lowers the CPU usage of the voice rendering.
// If LFO affects the low-pass filter it can be hearable even as low as 8.
// See tsf.h:240
const RenderBlockSize = C.TSF_RENDER_EFFECTSAMPLEBLOCK

// Supported output modes by the render methods
type OutputMode int

const (
	// Two channels with single left/right samples one after another
	OutputModeStereoInterleaved OutputMode = C.TSF_STEREO_INTERLEAVED
	// Two channels with all samples for the left channel first then right
	OutputModeStereoUnweaved OutputMode = C.TSF_STEREO_UNWEAVED
	// A single channel (stereo instruments are mixed into center)
	OutputModeMono OutputMode = C.TSF_MONO
)

type SoundFont struct {
	font *C.tsf
}

// Directly load a SoundFont from a .sf2 file path
func LoadSoundFontFile(filename string) SoundFont {
	return SoundFont{C.tsf_load_filename(C.CString(filename))}
}

// Load a SoundFont from a block of memory
func LoadSoundFontMemory(mem []byte) SoundFont {
	return SoundFont{C.tsf_load_memory(C.CBytes(mem), C.int(len(mem)))}
}

func (f SoundFont) IsNil() bool {
	return f.font == nil
}

// Free the memory related to this tsf instance
func (f SoundFont) Close() {
	C.tsf_close(f.font)
}

// Stop all playing notes immediatly and reset all channel parameters
func (f SoundFont) Reset() {
	C.tsf_reset(f.font)
}

// Returns the preset index from a bank and preset number, or -1 if it does not exist in the loaded SoundFont
func (f SoundFont) GetPresetIndex(bank, preset int) int {
	return int(C.tsf_get_presetindex(f.font, C.int(bank), C.int(preset)))
}

// Returns the number of presets in the loaded SoundFont
func (f SoundFont) GetPresetCount() int {
	return int(C.tsf_get_presetcount(f.font))
}

// Returns the name of a preset index >= 0 and < f.GetPresetCount()
func (f SoundFont) GetPresetName(preset int) string {
	return C.GoString(C.tsf_get_presetname(f.font, C.int(preset)))
}

// Returns the name of a preset by bank and preset number
func (f SoundFont) BankGetPresetName(bank, preset int) string {
	return C.GoString(C.tsf_bank_get_presetname(f.font, C.int(bank), C.int(preset)))
}

// Thread safety:
// Your audio output which calls the tsf_render* functions will most likely
// run on a different thread than where the playback tsf_note* functions
// are called. In which case some sort of concurrency control like a
// mutex needs to be used so they are not called at the same time.
// Alternatively, you can pre-allocate a maximum number of voices that can
// play simultaneously by calling tsf_set_max_voices after loading.
// That way memory re-allocation will not happen during tsf_note_on and
// TSF should become mostly thread safe.
// There is a theoretical chance that ending notes would negatively influence
// a voice that is rendering at the time but it is hard to say.
// Also be aware, this has not been tested much.

// Setup the parameters for the voice render methods
// mode: if mono or stereo and how stereo channel data is ordered
// sampleRate: the number of samples per second (output frequency)
// gain: volume gain in decibels (>0 means higher, <0 means lower)
func (f SoundFont) SetOutput(mode OutputMode, sampleRate int, gain float32) {
	C.tsf_set_output(f.font, uint32(mode), C.int(sampleRate), C.float(gain))
}

// Set the global gain as a volume factor
// volume: the desired volume where 1.0 is 100%
func (f SoundFont) SetVolume(volume float32) {
	C.tsf_set_volume(f.font, C.float(volume))
}

// Set the maximum number of voices to play simultaneously
// Depending on the soundfond, one note can cause many new voices to be started,
// so don't keep this number too low or otherwise sounds may not play.
// max_voices: maximum number to pre-allocate and set the limit to
func (f SoundFont) SetMaxVoices(max int) {
	C.tsf_set_max_voices(f.font, C.int(max))
}

// Start playing a note
// preset: preset index >= 0 and < f.GetPresetCount()
// key: note value between 0 and 127 (60 being middle C)
// vel: velocity as a float between 0.0 (equal to note off) and 1.0 (full)
func (f SoundFont) NoteOn(preset, key int, velocity float32) {
	C.tsf_note_on(f.font, C.int(preset), C.int(key), C.float(velocity))
}

// Start playing a note
// bank: instrument bank number (alternative to preset_index)
// preset: preset index >= 0 and < f.GetPresetCount()
// key: note value between 0 and 127 (60 being middle C)
// vel: velocity as a float between 0.0 (equal to note off) and 1.0 (full)
// returns 0 if preset does not exist, otherwise 1
func (f SoundFont) BankNoteOn(bank, preset, key int, velocity float32) int {
	return int(C.tsf_bank_note_on(f.font, C.int(bank), C.int(preset), C.int(key), C.float(velocity)))
}

// Stop playing a note
func (f SoundFont) NoteOff(preset, key int) {
	C.tsf_note_off(f.font, C.int(preset), C.int(key))
}

// Stop playing a note
// returns 0 if preset does not exist, otherwise 1
func (f SoundFont) BankNoteOff(bank, preset, key int) int {
	return int(C.tsf_bank_note_off(f.font, C.int(bank), C.int(preset), C.int(key)))
}

// Stop playing all notes (end with sustain and release)
func (f SoundFont) NoteOffAll() {
	C.tsf_note_off_all(f.font)
}

// Returns the number of active voices
func (f SoundFont) ActiveVoiceCount() int {
	return int(C.tsf_active_voice_count(f.font))
}

// Render output samples into a buffer
// dst: target buffer of size (samples * output_channels * sizeof(type))
// samples: number of samples to render
// mix: if 0 clear the buffer first, otherwise mix into existing data
func (f SoundFont) RenderShort(dst []int16, samples int, mix bool) {
	_mix := 0
	if mix {
		_mix = 1
	}
	C.tsf_render_short(f.font, (*C.short)(&dst[0]), C.int(samples), C.int(_mix))
}

// Render output samples into a buffer
// dst: target buffer of size (samples * output_channels * sizeof(type))
// samples: number of samples to render
// mix: if 0 clear the buffer first, otherwise mix into existing data
func (f SoundFont) RenderFloat(dst []float32, samples int, mix bool) {
	_mix := 0
	if mix {
		_mix = 1
	}
	C.tsf_render_float(f.font, (*C.float)(&dst[0]), C.int(samples), C.int(_mix))
}

//
// Higher level channel based functions, set up channel parameters
//

func (f SoundFont) ChannelSetPresetIndex(channel, preset int) {
	C.tsf_channel_set_presetindex(f.font, C.int(channel), C.int(preset))
}

// midiDrums: 0 for normal channels, otherwise apply MIDI drum channel rules
// returns 0 if preset does not exist, otherwise 1
func (f SoundFont) ChannelSetPresetNumber(channel, preset int, drums bool) int {
	_drums := 0
	if drums {
		_drums = 1
	}
	return int(C.tsf_channel_set_presetnumber(f.font, C.int(channel), C.int(preset), C.int(_drums)))
}

func (f SoundFont) ChannelSetBank(channel, bank int) {
	C.tsf_channel_set_bank(f.font, C.int(channel), C.int(bank))
}

// returns 0 if preset does not exist, otherwise 1
func (f SoundFont) ChannelSetBankPreset(channel, bank, preset int) int {
	return int(C.tsf_channel_set_bank_preset(f.font, C.int(channel), C.int(bank), C.int(preset)))
}

// pan: stereo panning value from 0.0 (left) to 1.0 (right) (default 0.5 center)
func (f SoundFont) ChannelSetPan(channel int, pan float32) {
	C.tsf_channel_set_pan(f.font, C.int(channel), C.float(pan))
}

// volume: linear volume scale factor (default 1.0 full)
func (f SoundFont) ChannelSetVolume(channel int, volume float32) {
	C.tsf_channel_set_volume(f.font, C.int(channel), C.float(volume))
}

// pitch: pitch wheel position 0 to 16383 (default 8192 unpitched)
func (f SoundFont) ChannelSetPitchWheel(channel, pitch int) {
	C.tsf_channel_set_pitchwheel(f.font, C.int(channel), C.int(pitch))
}

// pitch_range: range of the pitch wheel in semitones (default 2.0, total +/- 2 semitones)
func (f SoundFont) ChannelSetPitchRange(channel int, _range float32) {
	C.tsf_channel_set_pitchrange(f.font, C.int(channel), C.float(_range))
}

// tuning: tuning of all playing voices in semitones (default 0.0, standard (A440) tuning)
func (f SoundFont) ChannelSetTuning(channel int, tuning float32) {
	C.tsf_channel_set_tuning(f.font, C.int(channel), C.float(tuning))
}

// starts playing note
// key: note value between 0 and 127 (60 being middle C)
// vel: velocity as a float between 0.0 (equal to note off) and 1.0 (full)
func (f SoundFont) ChannelNoteOn(channel, key int, velocity float32) {
	C.tsf_channel_note_on(f.font, C.int(channel), C.int(key), C.float(velocity))
}

// stops playing note
// key: note value between 0 and 127 (60 being middle C)
func (f SoundFont) ChannelNoteOff(channel, key int) {
	C.tsf_channel_note_off(f.font, C.int(channel), C.int(key))
}

// end with sustain and release
func (f SoundFont) ChannelNoteOffAll(channel int) {
	C.tsf_channel_note_off_all(f.font, C.int(channel))
}

// end immediately
func (f SoundFont) ChannelSoundsOffAll(channel int) {
	C.tsf_channel_sounds_off_all(f.font, C.int(channel))
}

// Apply a MIDI control change to the channel (not all controllers are supported!)
func (f SoundFont) ChannelMidiControl(channel, controller, value int) {
	C.tsf_channel_midi_control(f.font, C.int(channel), C.int(controller), C.int(value))
}

//
// boring getters
//

func (f SoundFont) ChannelGetPresetIndex(channel int) int {
	return int(C.tsf_channel_get_preset_index(f.font, C.int(channel)))
}

func (f SoundFont) ChannelGetPresetBank(channel int) int {
	return int(C.tsf_channel_get_preset_bank(f.font, C.int(channel)))
}

func (f SoundFont) ChannelGetPresetNumber(channel int) int {
	return int(C.tsf_channel_get_preset_number(f.font, C.int(channel)))
}

func (f SoundFont) ChannelGetPan(channel int) float32 {
	return float32(C.tsf_channel_get_pan(f.font, C.int(channel)))
}

func (f SoundFont) ChannelGetVolume(channel int) float32 {
	return float32(C.tsf_channel_get_volume(f.font, C.int(channel)))
}

func (f SoundFont) ChannelGetPitchWheel(channel int) int {
	return int(C.tsf_channel_get_pitchwheel(f.font, C.int(channel)))
}

func (f SoundFont) ChannelGetPitchRange(channel int) float32 {
	return float32(C.tsf_channel_get_pitchrange(f.font, C.int(channel)))
}

func (f SoundFont) ChannelGetTuning(channel int) float32 {
	return float32(C.tsf_channel_get_tuning(f.font, C.int(channel)))
}
