package adpcm

import (
	"log"

	"github.com/qiniu/audio"
	"github.com/qiniu/audio/wav"
	"github.com/qiniu/x/bufiox"
)

// -------------------------------------------------------------------------------------

func decode(r *bufiox.Reader, cfg *wav.Config) (dec audio.Decoded, err error) {
	log.Println("adpcm.decode: TODO")
	return nil, audio.ErrFormat
}

const (
	adpcmFormat = 0x11
)

func init() {
	wav.RegisterFormat(adpcmFormat, decode)
}

// -------------------------------------------------------------------------------------
