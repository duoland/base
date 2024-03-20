package io

import "bufio"

type BufferWriter struct {
	bw  *bufio.Writer
	err error
}

func NewBufferWriter(bw *bufio.Writer) *BufferWriter {
	return &BufferWriter{
		bw:  bw,
		err: nil,
	}
}

func (w *BufferWriter) WriteString(str string) {
	if w.err != nil {
		return
	}
	_, err := w.bw.WriteString(str)
	w.err = err
}

func (w *BufferWriter) Error() error {
	return w.err
}

func (w *BufferWriter) Flush() (err error) {
	return w.bw.Flush()
}
