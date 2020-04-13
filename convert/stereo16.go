package convert

import (
	"io"

	"github.com/qiniu/x/bufiox"
)

// Stereo16 class.
type Stereo16 struct {
	source *bufiox.Reader
	mono   bool
	eight  bool
}

// NewStereo16 func.
func NewStereo16(source io.ReadSeeker, isMono, eightBitsPerSample bool) *Stereo16 {
	return &Stereo16{
		source: bufiox.NewReader(source),
		mono:   isMono,
		eight:  eightBitsPerSample,
	}
}

func (s *Stereo16) Read(b []uint8) (int, error) {
	l := len(b)
	if s.mono {
		l >>= 1
	}
	if s.eight {
		l >>= 1
	}
	buf := b[len(b)-l:]
	n, err := s.source.ReadFull(buf) // ReadFull: forbidden to read odd bytes (when !mono or !eight).
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		} else if err != io.EOF {
			return 0, err
		}
	}
	switch {
	case s.mono && s.eight:
		for i := 0; i < n; i++ {
			v := int16(int(buf[i])*0x101 - (1 << 15))
			b[4*i] = uint8(v)
			b[4*i+1] = uint8(v >> 8)
			b[4*i+2] = uint8(v)
			b[4*i+3] = uint8(v >> 8)
		}
	case s.mono && !s.eight:
		for i := 0; i < (n >> 1); i++ {
			b[4*i] = buf[2*i]
			b[4*i+1] = buf[2*i+1]
			b[4*i+2] = buf[2*i]
			b[4*i+3] = buf[2*i+1]
		}
	case !s.mono && s.eight:
		for i := 0; i < (n >> 1); i++ {
			v0 := int16(int(buf[2*i])*0x101 - (1 << 15))
			v1 := int16(int(buf[2*i+1])*0x101 - (1 << 15))
			b[4*i] = uint8(v0)
			b[4*i+1] = uint8(v0 >> 8)
			b[4*i+2] = uint8(v1)
			b[4*i+3] = uint8(v1 >> 8)
		}
	}
	if s.mono {
		n <<= 1
	}
	if s.eight {
		n <<= 1
	}
	return n, err
}

// Seek func.
func (s *Stereo16) Seek(offset int64, whence int) (int64, error) {
	if s.mono {
		offset >>= 1
	}
	if s.eight {
		offset >>= 1
	}
	return s.source.Seek(offset, whence)
}
