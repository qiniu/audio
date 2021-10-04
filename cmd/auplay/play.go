package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/oto"

	"github.com/qiniu/audio"
	_ "github.com/qiniu/audio/mp3"
	_ "github.com/qiniu/audio/wav"
	_ "github.com/qiniu/audio/wav/adpcm"
)

// -------------------------------------------------------------------------------------

func play(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	d, format, err := audio.Decode(f)
	if err != nil {
		return err
	}
	fmt.Printf(
		"Format: %s\nSampleRate: %d\nChannels: %d\nBytesPerSample: %d\n",
		format, d.SampleRate(), d.Channels(), d.BytesPerSample())

	c, err := oto.NewContext(d.SampleRate(), d.Channels(), d.BytesPerSample(), 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	fmt.Printf("Length: %d[bytes]\n", d.Length())
	p := c.NewPlayer()
	defer p.Close()

	_, err = io.Copy(p, d)
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: auplay <AudioFilePath>\n\n")
		return
	}
	if err := play(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

// -------------------------------------------------------------------------------------
