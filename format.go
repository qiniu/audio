package audio

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"

	"github.com/qiniu/x/bufiox"
)

// ErrFormat indicates that decoding encountered an unknown format.
var ErrFormat = errors.New("audio: unknown format")

// -------------------------------------------------------------------------------------

// Decoded represents a decoded audio.
type Decoded interface {
	io.ReadSeeker

	// SampleRate returns the sample rate like 44100.
	SampleRate() int

	// Length returns the total size in bytes. It returns -1 when the total size is not
	// available. e.g. when the given source is not io.Seeker.
	Length() int64
}

// Config holds an audio's configurations.
type Config struct {
	// TODO:
}

// DecodeFunc prototype.
type DecodeFunc = func(io.ReadSeeker) (Decoded, error)

// DecodeConfigFunc prototype.
type DecodeConfigFunc = func(io.ReadSeeker) (Config, error)

// -------------------------------------------------------------------------------------

// A format holds an audio format's name, magic header and how to decode it.
type format struct {
	name, magic  string
	decode       DecodeFunc
	decodeConfig DecodeConfigFunc
}

// Formats is the list of registered formats.
var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value
)

// RegisterFormat registers an audio format for use by Decode.
// Name is the name of the format, like "mp3" or "riff".
// Magic is the magic prefix that identifies the format's encoding. The magic
// string can contain "?" wildcards that each match any one byte.
// Decode is the function that decodes the encoded audio.
// DecodeConfig is the function that decodes just its configuration.
func RegisterFormat(name, magic string, decode DecodeFunc, decodeConfig DecodeConfigFunc) {
	formatsMu.Lock()
	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{name, magic, decode, decodeConfig}))
	formatsMu.Unlock()
}

// -------------------------------------------------------------------------------------

// A reader is an io.Reader that can also peek ahead.
type reader interface {
	io.ReadSeeker
	Peek(int) ([]byte, error)
}

// asReader converts an io.ReadSeeker to a reader.
func asReader(r io.ReadSeeker) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufiox.NewReader(r)
}

// Match reports whether magic matches b. Magic may contain "?" wildcards.
func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

// Sniff determines the format of r's data.
func sniff(r reader) format {
	formats, _ := atomicFormats.Load().([]format)
	for _, f := range formats {
		b, err := r.Peek(len(f.magic))
		if err == nil && match(f.magic, b) {
			return f
		}
	}
	return format{}
}

// Decode decodes an audio that has been encoded in a registered format.
// The string returned is the format name used during format registration.
// Format registration is typically done by an init function in the codec-
// specific package.
func Decode(r io.ReadSeeker) (Decoded, string, error) {
	rr := asReader(r)
	f := sniff(rr)
	if f.decode == nil {
		return nil, "", ErrFormat
	}
	m, err := f.decode(rr)
	return m, f.name, err
}

// DecodeConfig decodes the basic configurations of an audio that has
// been encoded in a registered format. The string returned is the format name
// used during format registration. Format registration is typically done by
// an init function in the codec-specific package.
func DecodeConfig(r io.ReadSeeker) (Config, string, error) {
	rr := asReader(r)
	f := sniff(rr)
	if f.decodeConfig == nil {
		return Config{}, "", ErrFormat
	}
	c, err := f.decodeConfig(rr)
	return c, f.name, err
}

// -------------------------------------------------------------------------------------
