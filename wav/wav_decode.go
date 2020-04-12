package wav

import (
	"bytes"
	"fmt"
	"io"

	"github.com/qiniu/audio"
	"github.com/qiniu/x/bufiox"
)

type stream struct {
	src            io.ReadSeeker
	headerSize     int64
	dataSize       int64
	remaining      int64
	channelNum     int
	bytesPerSample int
	sampleRate     int
}

// Read is implementation of io.Reader's Read.
func (s *stream) Read(p []byte) (int, error) {
	if s.remaining <= 0 {
		return 0, io.EOF
	}
	if s.remaining < int64(len(p)) {
		p = p[0:s.remaining]
	}
	n, err := s.src.Read(p)
	s.remaining -= int64(n)
	return n, err
}

// Seek is implementation of io.Seeker's Seek.
func (s *stream) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		offset = offset + s.headerSize
	case io.SeekCurrent:
	case io.SeekEnd:
		offset = s.headerSize + s.dataSize + offset
		whence = io.SeekStart
	}
	n, err := s.src.Seek(offset, whence)
	if err != nil {
		return 0, err
	}
	if n-s.headerSize < 0 {
		return 0, fmt.Errorf("wav: invalid offset")
	}
	s.remaining = s.dataSize - (n - s.headerSize)
	// There could be a tail in wav file.
	if s.remaining < 0 {
		s.remaining = 0
		return s.dataSize, nil
	}
	return n - s.headerSize, nil
}

// Length returns the size of decoded stream in bytes.
func (s *stream) Length() int64 {
	return s.dataSize
}

// SampleRate returns the sample rate like 44100.
func (s *stream) SampleRate() int {
	return s.sampleRate
}

// Channels returns the number of channels. One channel is mono playback.
// Two channels are stereo playback. No other values are supported.
func (s *stream) Channels() int {
	return s.channelNum
}

// BytesPerSample returns the number of bytes per sample per channel.
// The usual value is 2. Only values 1 and 2 are supported.
func (s *stream) BytesPerSample() int {
	return s.bytesPerSample
}

func decode(src io.ReadSeeker) (*stream, error) {
	buf := make([]byte, 12)
	n, err := io.ReadFull(src, buf)
	if n != len(buf) {
		return nil, fmt.Errorf("wav: invalid header")
	}
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(buf[0:4], []byte("RIFF")) {
		return nil, fmt.Errorf("wav: invalid header: 'RIFF' not found")
	}
	if !bytes.Equal(buf[8:12], []byte("WAVE")) {
		return nil, fmt.Errorf("wav: invalid header: 'WAVE' not found")
	}

	// Read chunks
	dataSize := int64(0)
	headerSize := int64(len(buf))
	sampleRate := int64(0)
	channelNum := 0
	bitsPerSample := 0
chunks:
	for {
		buf := make([]byte, 8)
		n, err := io.ReadFull(src, buf)
		if n != len(buf) {
			return nil, fmt.Errorf("wav: invalid header")
		}
		if err != nil {
			return nil, err
		}
		headerSize += 8
		size := int64(buf[4]) | int64(buf[5])<<8 | int64(buf[6])<<16 | int64(buf[7])<<24
		switch {
		case bytes.Equal(buf[0:4], []byte("fmt ")):
			// Size of 'fmt' header is usually 16, but can be more than 16.
			if size < 16 {
				return nil, fmt.Errorf("wav: invalid header: maybe non-PCM file?")
			}
			buf := make([]byte, size)
			n, err := io.ReadFull(src, buf)
			if n != len(buf) {
				return nil, fmt.Errorf("wav: invalid header")
			}
			if err != nil {
				return nil, err
			}
			format := int(buf[0]) | int(buf[1])<<8
			if format != 1 {
				return nil, fmt.Errorf("wav: format must be linear PCM")
			}
			channelNum = int(buf[2]) | int(buf[3])<<8
			switch channelNum {
			case 1, 2:
			default:
				return nil, fmt.Errorf("wav: channel num must be 1 or 2 but was %d", channelNum)
			}
			bitsPerSample = int(buf[14]) | int(buf[15])<<8
			if bitsPerSample != 8 && bitsPerSample != 16 {
				return nil, fmt.Errorf("wav: bits per sample must be 8 or 16 but was %d", bitsPerSample)
			}
			sampleRate = int64(buf[4]) | int64(buf[5])<<8 | int64(buf[6])<<16 | int64(buf[7])<<24
			headerSize += size
		case bytes.Equal(buf[0:4], []byte("data")):
			dataSize = size
			break chunks
		default:
			buf := make([]byte, size)
			n, err := io.ReadFull(src, buf)
			if n != len(buf) {
				return nil, fmt.Errorf("wav: invalid header")
			}
			if err != nil {
				return nil, err
			}
			headerSize += size
		}
	}
	s := &stream{
		src:            src,
		headerSize:     headerSize,
		dataSize:       dataSize,
		remaining:      dataSize,
		sampleRate:     int(sampleRate),
		channelNum:     channelNum,
		bytesPerSample: bitsPerSample >> 3,
	}
	return s, nil
}

// -------------------------------------------------------------------------------------

// Decode decodes a wav audio.
func Decode(r io.ReadSeeker) (audio.Decoded, error) {
	b := bufiox.NewReader(r)
	dec, err := decode(b)
	return dec, err
}

// DecodeConfig is not implemented.
func DecodeConfig(r io.ReadSeeker) (cfg audio.Config, err error) {
	err = audio.ErrFormat
	return
}

func init() {
	audio.RegisterFormat("wav", "RIFF????WAVE", Decode, DecodeConfig)
}

// -------------------------------------------------------------------------------------
