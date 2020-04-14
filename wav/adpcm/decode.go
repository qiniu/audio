package adpcm

import (
	"errors"
	"fmt"

	"github.com/qiniu/audio"
	"github.com/qiniu/audio/wav"
	"github.com/qiniu/x/bufiox"
)

var (
	errNotImpl = errors.New("not impl")
)

// -------------------------------------------------------------------------------------

type decoded struct {
	src *bufiox.Reader

	samples   []byte
	remaining int

	length     int64
	sampleRate int
	channelNum int
	blockAlign int
}

func newDecoded(src *bufiox.Reader, cfg *wav.Config) *decoded {
	samplesPerBlock := cfg.SamplesPerBlock()
	return &decoded{
		src:        bufiox.NewReaderSize(src, cfg.BlockAlign),
		sampleRate: cfg.SampleRate,
		channelNum: cfg.Channels,
		blockAlign: cfg.BlockAlign,
		samples:    make([]byte, samplesPerBlock<<1),
		length:     cfg.DataSize / int64(cfg.BlockAlign) * int64(samplesPerBlock<<1),
	}
}

func (p *decoded) nextBlock() error {
	block, err := p.src.Peek(p.blockAlign)
	if err != nil {
		return err
	}
	loadBlock(p.channelNum, block, p.samples)
	p.src.Discard(p.blockAlign)
	p.remaining = p.blockAlign
	return nil
}

func (p *decoded) Read(b []byte) (n int, err error) {
	if p.remaining <= 0 {
		if err = p.nextBlock(); err != nil {
			return
		}
	}
	n = copy(b, p.samples[p.blockAlign-p.remaining:])
	p.remaining -= n
	return
}

func (p *decoded) Seek(offset int64, whence int) (newoff int64, err error) {
	return 0, errNotImpl
}

// Length returns the size of decoded stream in bytes.
func (p *decoded) Length() int64 {
	return p.length
}

// SampleRate returns the sample rate like 44100.
func (p *decoded) SampleRate() int {
	return p.sampleRate
}

// Channels returns the number of channels. One channel is mono playback.
// Two channels are stereo playback. No other values are supported.
func (p *decoded) Channels() int {
	return p.channelNum
}

// BytesPerSample returns the number of bytes per sample per channel.
// The usual value is 2. Only values 1 and 2 are supported.
func (p *decoded) BytesPerSample() int {
	return 2
}

// -------------------------------------------------------------------------------------

func decode(src *bufiox.Reader, cfg *wav.Config) (dec audio.Decoded, err error) {
	if cfg.BitsPerSample != 4 {
		return nil, fmt.Errorf("adpcm wav: bits per sample must be 4 but was %d", cfg.BitsPerSample)
	}
	d := newDecoded(src, cfg)
	return d, nil
}

const (
	adpcmFormat = 0x11
)

func init() {
	wav.RegisterFormat(adpcmFormat, decode)
}

// -------------------------------------------------------------------------------------
