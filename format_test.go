package audio_test

import (
	"bytes"
	"testing"

	"github.com/qiniu/audio"
	_ "github.com/qiniu/audio/mp3"
)

func Test(t *testing.T) {
	b := bytes.NewReader(nil)
	audio.DecodeConfig(b)
}
