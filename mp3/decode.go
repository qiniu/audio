package mp3

import (
	"io"

	"github.com/qiniu/audio"

	mp3 "github.com/hajimehoshi/go-mp3"
)

// -------------------------------------------------------------------------------------

func decode(r io.ReadSeeker) (audio.Decoder, error) {
	dec, err := mp3.NewDecoder(r)
	return dec, err
}

func decodeConfig(r io.ReadSeeker) (cfg audio.Config, err error) {
	err = audio.ErrFormat
	return
}

func init() {
	audio.RegisterFormat("mp3", "ID3", decode, decodeConfig)
}

// -------------------------------------------------------------------------------------
