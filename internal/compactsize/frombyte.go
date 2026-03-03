package compactsize

import "fmt"

func FromBytes(b []byte) (uint64, error) {
	if len(b) != 1 {
		return 0, fmt.Errorf("invalid byte length")
	}
	return uint64(b[0]), nil
}
