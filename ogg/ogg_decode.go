package ogg

/*
import (
	"io"

	"github.com/qiniu/audio"
	"github.com/qiniu/audio/convert"
)

type decoder interface {
	Read([]float32) (int, error)
	SetPosition(int64) error
	Length() int64
	Channels() int
	SampleRate() int
}

type decoded struct {
	totalBytes int
	posInBytes int
	decoder    decoder
	decoderr   io.Reader
}

func (d *decoded) Length() int64 {
	return 0 // TODO
}

func (d *decoded) Read(b []byte) (int, error) {
	if d.decoderr == nil {
		d.decoderr = convert.NewReaderFromFloat32Reader(d.decoder)
	}

	l := d.totalBytes - d.posInBytes
	if l > len(b) {
		l = len(b)
	}
	if l < 0 {
		return 0, io.EOF
	}

retry:
	n, err := d.decoderr.Read(b[:l])
	if err != nil && err != io.EOF {
		return 0, err
	}
	if n == 0 && l > 0 && err != io.EOF {
		// When l is too small, decoder's Read might return 0 for a while. Let's retry.
		goto retry
	}

	d.posInBytes += n
	if d.posInBytes == d.totalBytes || err == io.EOF {
		return n, io.EOF
	}
	return n, nil
}

func (d *decoded) Seek(offset int64, whence int) (int64, error) {
	next := int64(0)
	switch whence {
	case io.SeekStart:
		next = offset
	case io.SeekCurrent:
		next = int64(d.posInBytes) + offset
	case io.SeekEnd:
		next = int64(d.totalBytes) + offset
	}
	// pos should be always even
	next = next / 2 * 2
	d.posInBytes = int(next)
	d.decoder.SetPosition(next / int64(d.decoder.Channels()) / 2)
	d.decoderr = nil
	return next, nil
}

// SampleRate returns the sample rate like 44100.
func (d *decoded) SampleRate() int {
	return d.decoder.SampleRate()
}

// Channels func.
func (d *decoded) Channels() int {
	return d.decoder.Channels()
}

// Decode accepts an ogg stream and returns a decorded audio.
func Decode(in io.ReadSeeker) (audio.Decoded, error) {
	r, err := newDecoder(in)
	if err != nil {
		return nil, err
	}
	d := &decoded{
		// TODO: r.Length() returns 0 when the format is unknown.
		// Should we check that?
		totalBytes: int(r.Length()) * r.Channels() * 2, // 2 means 16bit per sample.
		posInBytes: 0,
		decoder:    r,
	}
	return d, nil
}

// DecodeConfig is not implemented.
func DecodeConfig(r io.ReadSeeker) (cfg audio.Config, err error) {
	err = audio.ErrFormat
	return
}

func init() {
	audio.RegisterFormat("ogg", "OggS", Decode, DecodeConfig)
}
*/
// -------------------------------------------------------------------------------------
