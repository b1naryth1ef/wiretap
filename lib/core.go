package wiretap

import (
	"io"
)

type Wiretapper struct {
	Complete chan bool

	rec *Recording
	src io.Reader
	dst io.Writer
}

func NewWiretapper(rec *Recording, src io.Reader, dst io.Writer) *Wiretapper {
	return &Wiretapper{
		rec: rec,
		src: src,
		dst: dst,
	}
}

func (w *Wiretapper) Run() {

	go func() {
		w.rec.Open()
		var err error
		var size int
		buffer := make([]byte, 4096)

		for err != io.EOF {
			size, err = w.src.Read(buffer)
			w.dst.Write(buffer[:size])
			w.rec.Write(buffer[:size])
		}

		w.Complete <- true
		w.rec.Close()
	}()
}
