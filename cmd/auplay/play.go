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
	fmt.Printf("Format: %s\nSampleRate: %d\n", format, d.SampleRate())

	c, err := oto.NewContext(d.SampleRate(), d.Channels(), d.BytesPerSample(), 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	fmt.Printf("Length: %d[bytes]\n", d.Length())
	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: auplay <AudioFilePath>\n")
		return
	}
	if err := play(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

// -------------------------------------------------------------------------------------
