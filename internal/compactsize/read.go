package compactsize

import (
	"encoding/binary"
	"fmt"
	"io"
)

func Read(r io.Reader) (uint64, error) {
	var prefix [1]byte
	if _, err := io.ReadFull(r, prefix[:]); err != nil {
		return 0, err
	}

	switch prefix[0] {

	case 0xFF:
		var buf [8]byte
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return 0, err
		}
		val := binary.LittleEndian.Uint64(buf[:])
		if val <= 0xFFFFFFFF {
			return 0, fmt.Errorf("non-minimal compact size")
		}
		return val, nil

	case 0xFE:
		var buf [4]byte
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return 0, err
		}
		val := uint64(binary.LittleEndian.Uint32(buf[:]))
		if val <= 0xFFFF {
			return 0, fmt.Errorf("non-minimal compact size")
		}
		return val, nil

	case 0xFD:
		var buf [2]byte
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return 0, err
		}
		val := uint64(binary.LittleEndian.Uint16(buf[:]))
		if val < 0xFD {
			return 0, fmt.Errorf("non-minimal compact size")
		}
		return val, nil

	default:
		return uint64(prefix[0]), nil
	}
}
