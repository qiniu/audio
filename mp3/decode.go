package mp3

import (
	"io"

	"github.com/qiniu/audio"
	"github.com/qiniu/x/bufiox"

	mp3 "github.com/hajimehoshi/go-mp3"
)

// -------------------------------------------------------------------------------------

// Decode decodes a mp3 audio.
func Decode(r io.ReadSeeker) (audio.Decoded, error) {
	b := bufiox.NewReader(r)
	dec, err := mp3.NewDecoder(b)
	return dec, err
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
