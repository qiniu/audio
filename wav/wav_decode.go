package wav

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/qiniu/audio"
	"github.com/qiniu/x/bufiox"
)

// -------------------------------------------------------------------------------------

// Config type.
type Config struct {
	Format        int
	SampleRate    int
	DataSize      int64
	HeaderSize    int64
	Channels      int
	BitsPerSample int
	BlockAlign    int
}

// SamplesPerBlock returns samples per block.
func (p *Config) SamplesPerBlock() int {
	return (p.BlockAlign << 1) - p.Channels*7
}

// DecodeFunc prototype.
type DecodeFunc = func(r *bufiox.Reader, cfg *Config) (audio.Decoded, error)

// A format holds an audio format's name, magic header and how to decode it.
type format struct {
	tag    int
	decode DecodeFunc
}

// Formats is the list of registered formats.
var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value
)

// RegisterFormat registers a wav decoder extension.
func RegisterFormat(tag int, decode DecodeFunc) {
	formatsMu.Lock()
	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{tag, decode}))
	formatsMu.Unlock()
}

func decodeEx(r *bufiox.Reader, cfg *Config) (audio.Decoded, error) {
	formats, _ := atomicFormats.Load().([]format)
	for _, f := range formats {
		if f.tag == cfg.Format {
			return f.decode(r, cfg)
		}
	}
	return nil, audio.ErrFormat
}

// -------------------------------------------------------------------------------------

type stream struct {
	src            *bufiox.Reader
	headerSize     int64
	dataSize       int64
	remaining      int64
	channelNum     int
	bytesPerSample int
	sampleRate     int
}

func (s *stream) Read(p []byte) (int, error) {
	if s.remaining < int64(len(p)) {
		if s.remaining <= 0 {
			return 0, io.EOF
		}
		p = p[0:s.remaining]
	}
	n, err := s.src.Read(p)
	s.remaining -= int64(n)
	return n, err
}

func (s *stream) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		offset = offset + s.headerSize
	case io.SeekCurrent:
	case io.SeekEnd:
		offset = s.headerSize + s.dataSize + offset
		whence = io.SeekStart
	}
	n, err := s.src.Seek(offset, whence)
	if err != nil {
		return 0, err
	}
	if n-s.headerSize < 0 {
		return 0, fmt.Errorf("wav: invalid offset")
	}
	s.remaining = s.dataSize - (n - s.headerSize)
	// There could be a tail in wav file.
	if s.remaining < 0 {
		s.remaining = 0
		return s.dataSize, nil
	}
	return n - s.headerSize, nil
}

// Length returns the size of decoded stream in bytes.
func (s *stream) Length() int64 {
	return s.dataSize
}

// SampleRate returns the sample rate like 44100.
func (s *stream) SampleRate() int {
	return s.sampleRate
}

// Channels returns the number of channels. One channel is mono playback.
// Two channels are stereo playback. No other values are supported.
func (s *stream) Channels() int {
	return s.channelNum
}

// BytesPerSample returns the number of bytes per sample per channel.
// The usual value is 2. Only values 1 and 2 are supported.
func (s *stream) BytesPerSample() int {
	return s.bytesPerSample
}

var (
	errInvalidFormat = errors.New("wav: invalid header")
	errNotRIFF       = errors.New("wav: invalid header: 'RIFF' not found")
	errNotWAVE       = errors.New("wav: invalid header: 'WAVE' not found")
)

func decode(src *bufiox.Reader) (audio.Decoded, error) {
	b := make([]byte, 16)
	if _, err := src.ReadFull(b[:12]); err != nil {
		return nil, errInvalidFormat
	}
	if !bytes.Equal(b[0:4], []byte("RIFF")) {
		return nil, errNotRIFF
	}
	if !bytes.Equal(b[8:12], []byte("WAVE")) {
		return nil, errNotWAVE
	}

	// Read chunks
	cfg := Config{HeaderSize: int64(len(b))}
chunks:
	for {
		if _, err := src.ReadFull(b[:8]); err != nil {
			return nil, errInvalidFormat
		}
		cfg.HeaderSize += 8
		size := int64(b[4]) | int64(b[5])<<8 | int64(b[6])<<16 | int64(b[7])<<24
		switch {
		case bytes.Equal(b[0:4], []byte("fmt ")):
			if size < 16 { // Size of 'fmt' header is usually 16, but can be more than 16.
				return nil, errInvalidFormat
			}
			buf, err := src.Peek(int(size))
			if err != nil {
				return nil, errInvalidFormat
			}
			cfg.Format = int(buf[0]) | int(buf[1])<<8
			cfg.Channels = int(buf[2]) | int(buf[3])<<8
			switch cfg.Channels {
			case 1, 2:
			default:
				return nil, fmt.Errorf("wav: channel num must be 1 or 2 but was %d", cfg.Channels)
			}
			cfg.BlockAlign = int(buf[12]) | int(buf[13])<<8
			cfg.BitsPerSample = int(buf[14]) | int(buf[15])<<8
			cfg.SampleRate = int(buf[4]) | int(buf[5])<<8 | int(buf[6])<<16 | int(buf[7])<<24
			src.Discard(int(size))
			cfg.HeaderSize += size
		case bytes.Equal(b[0:4], []byte("data")):
			cfg.DataSize = size
			break chunks
		default:
			if _, err := src.Discard(int(size)); err != nil {
				return nil, err
			}
			cfg.HeaderSize += size
		}
	}
	if cfg.Format != 1 {
		return decodeEx(src, &cfg)
	}
	if cfg.BitsPerSample != 8 && cfg.BitsPerSample != 16 {
		return nil, fmt.Errorf("wav: bits per sample must be 8 or 16 but was %d", cfg.BitsPerSample)
	}
	s := &stream{
		src:            src,
		headerSize:     cfg.HeaderSize,
		dataSize:       cfg.DataSize,
		remaining:      cfg.DataSize,
		sampleRate:     int(cfg.SampleRate),
		channelNum:     cfg.Channels,
		bytesPerSample: cfg.BitsPerSample >> 3,
	}
	return s, nil
}

// -------------------------------------------------------------------------------------

// Decode decodes a wav audio.
func Decode(r io.ReadSeeker) (audio.Decoded, error) {
	b := bufiox.NewReader(r)
	return decode(b)
}

// DecodeConfig is not implemented.
func DecodeConfig(r io.ReadSeeker) (cfg audio.Config, err error) {
	err = audio.ErrFormat
	return
}

func init() {
	audio.RegisterFormat("wav", "RIFF????WAVE", Decode, DecodeConfig)
}

// -------------------------------------------------------------------------------------
