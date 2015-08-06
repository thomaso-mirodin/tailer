package tailer

import (
	"io"
	"time"
)

// Read as much data is available in the file into the ring buffer
func (t *Tailer) fill() error {
	t.fmu.Lock()
	_, err := io.Copy(t.ring, t.file)
	t.fmu.Unlock()
	switch err {
	case nil, io.ErrShortWrite, io.EOF:
		return nil
	default:
		return err
	}
}

func (t *Tailer) pollForChanges() {
	for {
		if t.closed {
			break
		}

		if err := t.fill(); err != nil {
			if err = t.reopenFile(); err != nil {
				t.errc <- err
			}
		}

		time.Sleep(pollIntervalFast)
	}
}

// func (t *Tailer) notifyForChanges() {
// 	// tbd
// }