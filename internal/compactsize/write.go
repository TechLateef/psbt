package compactsize

import (
	"encoding/binary"
	"io"
)

func Write(w io.Writer, n uint64) error {
	if n < 0xFD {

		_, err := w.Write([]byte{byte(n)})
		return err
	} else if n <= 0xFFFF {
		if _, err := w.Write([]byte{0xFD}); err != nil {
			return err
		}
		var buf [2]byte
		binary.LittleEndian.PutUint16(buf[:], uint16(n))
		_, err := w.Write(buf[:])
		return err
	} else if n <= 0xFFFFFFFF {
		if _, err := w.Write([]byte{0xFE}); err != nil {
			return err
		}
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], uint32(n))
		_, err := w.Write(buf[:])
		return err
	} else {
		if _, err := w.Write([]byte{0xFF}); err != nil {
			return err
		}
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], n)
		_, err := w.Write(buf[:])
		return err
	}
}
