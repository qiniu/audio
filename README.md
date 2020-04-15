# Audio support for Go language

[![Build Status](https://travis-ci.org/qiniu/audio.svg?branch=master)](https://travis-ci.org/qiniu/audio) [![GoDoc](https://godoc.org/github.com/qiniu/audio?status.svg)](https://godoc.org/github.com/qiniu/audio)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://www.qiniu.com/)

The package `github.com/qiniu/audio` is an extensible audio library with simple API for multi platforms in the Go programming language.

## Platforms

* Windows
* macOS
* Linux
* FreeBSD
* Android
* iOS
* Web browsers (Chrome, Firefox, Safari and Edge)
  * GopherJS
  * WebAssembly (Experimental)

## Features

* Pluggable audio decoders. And now it supports the following formats:
  * wav/pcm: `import _ "github.com/qiniu/audio/wav"`
  * wav/adpcm: `import _ "github.com/qiniu/audio/wav/adpcm"`
  * mp3: `import _ "github.com/qiniu/audio/mp3"`
* Audio encoders (TODO).
* Convert decoded audio stream.

## Example

```
import (
	"io"
	"os"

	"github.com/hajimehoshi/oto"

	"github.com/qiniu/audio"
	_ "github.com/qiniu/audio/mp3"
	_ "github.com/qiniu/audio/wav"
	_ "github.com/qiniu/audio/wav/adpcm"
)

func playAudio(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	d, _, err := audio.Decode(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), d.Channels(), d.BytesPerSample(), 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	_, err = io.Copy(p, d)
	return err
}
```

## Document

* See https://godoc.org/github.com/qiniu/audio
