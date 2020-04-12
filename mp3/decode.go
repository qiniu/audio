package mp3

import (
	"io"

	"github.com/qiniu/audio"
	"github.com/qiniu/x/bufiox"

	mp3 "github.com/hajimehoshi/go-mp3"
)

// -------------------------------------------------------------------------------------

type decoded struct {
	mp3.Decoder
}

// Channels returns the number of channels. One channel is mono playback.
// Two channels are stereo playback. No other values are supported.
func (p *decoded) Channels() int {
	return 2
}

// BytesPerSample returns the number of bytes per sample per channel.
// The usual value is 2. Only values 1 and 2 are supported.
func (p *decoded) BytesPerSample() int {
	return 2
}

// Decode decodes a mp3 audio.
func Decode(r io.ReadSeeker) (audio.Decoded, error) {
	b := bufiox.NewReader(r)
	dec, err := mp3.NewDecoder(b)
	return &decoded{Decoder: *dec}, err
}

// DecodeConfig is not implemented.
func DecodeConfig(r io.ReadSeeker) (cfg audio.Config, err error) {
	err = audio.ErrFormat
	return
}

func init() {
	audio.RegisterFormat("mp3", "ID3", Decode, DecodeConfig)
}

// -------------------------------------------------------------------------------------
