# Audio support for Go language

[![LICENSE](https://img.shields.io/github/license/qiniu/audio.svg)](https://github.com/qiniu/audio/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/qiniu/audio.png?branch=master)](https://travis-ci.org/qiniu/audio)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/audio)](https://goreportcard.com/report/github.com/qiniu/audio)
[![GitHub release](https://img.shields.io/github/v/tag/qiniu/audio.svg?label=release)](https://github.com/qiniu/audio/releases)
[![Coverage Status](https://codecov.io/gh/qiniu/audio/branch/master/graph/badge.svg)](https://codecov.io/gh/qiniu/audio)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/qiniu/audio)

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

	"github.com/hajimehoshi/oto/v2"

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
