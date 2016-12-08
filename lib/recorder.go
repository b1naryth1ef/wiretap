package wiretap

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"time"
)

const (
	MODE_RECORD = iota
	MODE_REPLAY
)

type FrameHeader struct {
	Offset int64
	Size   uint32
}

type Recording struct {
	Mode int
	Path string

	start time.Time
	file  *os.File
}

func NewRecording(path string, mode int) *Recording {
	return &Recording{
		Mode: mode,
		Path: path,
	}
}

func (r *Recording) Open() (err error) {
	r.start = time.Now()

	if r.Mode == MODE_RECORD {
		file, err := os.Create(r.Path)
		if err == nil {
			r.file = file
		}
	} else {
		file, err := os.Open(r.Path)
		if err == nil {
			r.file = file
		}
	}

	return
}

// TODO: this entire function is a performance nightmare
func (r *Recording) Play(out io.Writer) {
	var header FrameHeader

	headerBuffer := make([]byte, binary.Size(header))

	for {
		// Read a single frame header at a time
		_, err := r.file.Read(headerBuffer)
		if err != nil {
			return
		}

		reader := bytes.NewReader(headerBuffer)
		err = binary.Read(reader, binary.BigEndian, &header)
		if err != nil {
			return
		}

		buffer := make([]byte, header.Size)
		r.file.Read(buffer)

		when := r.start.Add(time.Duration(header.Offset) * time.Nanosecond)
		time.Sleep(when.Sub(time.Now()))
		out.Write(buffer)
	}
}

func (r *Recording) Close() {
	r.file.Close()
}

func (r *Recording) Write(data []byte) {
	r.writeFrame(data)
}

func (r *Recording) writeFrame(data []byte) {
	if len(data) == 0 {
		return
	}

	frame := FrameHeader{
		Offset: (time.Now().Sub(r.start)).Nanoseconds(),
		Size:   uint32(len(data)),
	}

	binary.Write(r.file, binary.BigEndian, &frame)
	r.file.Write(data)
}
