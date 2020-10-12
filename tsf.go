package tsf

//#cgo LDFLAGS: -lm
//#define TSF_STATIC 1
//#define TSF_IMPLEMENTATION 1
//#include "tsf.h"
import "C"

// Supported output modes by the render methods
type OutputMode int

const (
	// Two channels with single left/right samples one after another
	StereoInterleaved OutputMode = C.TSF_STEREO_INTERLEAVED
	// Two channels with all samples for the left channel first then right
	StereoUnweaved OutputMode = C.TSF_STEREO_UNWEAVED
	// A single channel (stereo instruments are mixed into center)
	Mono OutputMode = C.TSF_MONO
)

type SoundFont struct {
	tsf *C.tsf
}

// Directly load a SoundFont from a .sf2 file path
func LoadFilename(filename string) *SoundFont {
	return &SoundFont{C.tsf_load_filename(C.CString(filename))}
}

// Load a SoundFont from a block of memory
func LoadMemory(data []byte) *SoundFont {
	return &SoundFont{C.tsf_load_memory(C.CBytes(data), C.int(len(data)))}
}

// Free the memory related to this tsf instance
func (f *SoundFont) Close() {
	C.tsf_close(f.tsf)
}

// Stop all playing notes immediatly and reset all channel parameters
func (f *SoundFont) Reset() {
	C.tsf_reset(f.tsf)
}

// Returns the preset index from a bank and preset number, or -1 if it does not exist in the loaded SoundFont
func (f *SoundFont) GetPresetIndex(bank, preset int) int {
	return int(C.tsf_get_presetindex(f.tsf, C.int(bank), C.int(preset)))
}

// Returns the number of presets in the loaded SoundFont
func (f *SoundFont) GetPresetCount() int {
	return int(C.tsf_get_presetcount(f.tsf))
}

// Returns the name of a preset index >= 0 and < f.GetPresetCount()
func (f *SoundFont) GetPresetName(preset int) string {
	return C.GoString(C.tsf_get_presetname(f.tsf, C.int(preset)))
}

// Returns the name of a preset by bank and preset number
func (f *SoundFont) GetBankPresetName(bank, preset int) string {
	return C.GoString(C.tsf_bank_get_presetname(f.tsf, C.int(bank), C.int(preset)))
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
func (f *SoundFont) SetOutput(mode OutputMode, sampleRate int, gain float32) {
	C.tsf_set_output(f.tsf, uint32(mode), C.int(sampleRate), C.float(gain))
}

// Set the global gain as a volume factor
// volume: the desired volume where 1.0 is 100%
func (f *SoundFont) SetVolume(volume float32) {
	C.tsf_set_volume(f.tsf, C.float(volume))
}

// Set the maximum number of voices to play simultaneously
// Depending on the soundfond, one note can cause many new voices to be started,
// so don't keep this number too low or otherwise sounds may not play.
// max_voices: maximum number to pre-allocate and set the limit to
func (f *SoundFont) SetMaxVoices(max int) {
	C.tsf_set_max_voices(f.tsf, C.int(max))
}

// Start playing a note
// preset: preset index >= 0 and < f.GetPresetCount()
// key: note value between 0 and 127 (60 being middle C)
// vel: velocity as a float between 0.0 (equal to note off) and 1.0 (full)
func (f *SoundFont) NoteOn(preset, key int, velocity float32) {
	C.tsf_note_on(f.tsf, C.int(preset), C.int(key), C.float(velocity))
}

// Start playing a note
// bank: instrument bank number (alternative to preset_index)
// preset: preset index >= 0 and < f.GetPresetCount()
// key: note value between 0 and 127 (60 being middle C)
// vel: velocity as a float between 0.0 (equal to note off) and 1.0 (full)
// returns 0 if preset does not exist, otherwise 1
func (f *SoundFont) BankNoteOn(bank, preset, key int, velocity float32) int {
	return int(C.tsf_bank_note_on(f.tsf, C.int(bank), C.int(preset), C.int(key), C.float(velocity)))
}

// Stop playing a note
func (f *SoundFont) NoteOff(preset, key int) {
	C.tsf_note_off(f.tsf, C.int(preset), C.int(key))
}

// Stop playing a note
// returns 0 if preset does not exist, otherwise 1
func (f *SoundFont) BankNoteOff(bank, preset, key int) int {
	return int(C.tsf_bank_note_off(f.tsf, C.int(bank), C.int(preset), C.int(key)))
}

// Stop playing all notes (end with sustain and release)
func (f *SoundFont) AllNotesOff() {
	C.tsf_note_off_all(f.tsf)
}

// Returns the number of active voices
func (f *SoundFont) ActiveVoiceCount() int {
	return int(C.tsf_active_voice_count(f.tsf))
}

// Render output samples into a buffer
// dst: target buffer of size samples * output_channels * sizeof(type)
// samples: number of samples to render
// mix: if 0 clear the buffer first, otherwise mix into existing data
func (f *SoundFont) RenderShort(dst []int16, samples, mix int) {
	C.tsf_render_short(f.tsf, (*C.short)(&dst[0]), C.int(samples), C.int(mix))
}

// Render output samples into a buffer
// dst: target buffer of size samples * output_channels * sizeof(type)
// samples: number of samples to render
// mix: if 0 clear the buffer first, otherwise mix into existing data
func (f *SoundFont) RenderFloat(dst []float32, samples, mix int) {
	C.tsf_render_float(f.tsf, (*C.float)(&dst[0]), C.int(samples), C.int(mix))
}
