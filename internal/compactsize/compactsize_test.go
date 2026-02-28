package compactsize

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRead_SingleByte(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint64
		wantErr  bool
	}{
		{"zero", []byte{0x00}, 0, false},
		{"one", []byte{0x01}, 1, false},
		{"max single byte", []byte{0xFC}, 252, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			result, err := Read(r)

			if tc.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestRead_TwoBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint64
		wantErr  bool
	}{
		{"min 2-byte", []byte{0xFD, 0xFD, 0x00}, 253, false},
		{"max 2-byte", []byte{0xFD, 0xFF, 0xFF}, 65535, false},
		{"mid 2-byte", []byte{0xFD, 0x00, 0x01}, 256, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			result, err := Read(r)

			if tc.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestRead_FourBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint64
		wantErr  bool
	}{
		{"min 4-byte", []byte{0xFE, 0x00, 0x00, 0x01, 0x00}, 65536, false},
		{"max 4-byte", []byte{0xFE, 0xFF, 0xFF, 0xFF, 0xFF}, 4294967295, false},
		{"mid 4-byte", []byte{0xFE, 0x00, 0x00, 0x00, 0x01}, 16777216, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			result, err := Read(r)

			if tc.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestRead_EightBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected uint64
		wantErr  bool
	}{
		{"min 8-byte", []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, 4294967296, false},
		{"large 8-byte", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, 18446744073709551615, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			result, err := Read(r)

			if tc.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestRead_NonMinimalEncoding(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{"2-byte non-minimal", []byte{0xFD, 0xFC, 0x00}},                                     // 252 encoded as 2-byte
		{"4-byte non-minimal", []byte{0xFE, 0xFF, 0xFF, 0x00, 0x00}},                         // 65535 encoded as 4-byte
		{"8-byte non-minimal", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00}}, // 4294967295 encoded as 8-byte
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			_, err := Read(r)

			if err == nil {
				t.Errorf("expected non-minimal encoding error but got none")
			}
			if !strings.Contains(err.Error(), "non-minimal") {
				t.Errorf("expected non-minimal error, got: %v", err)
			}
		})
	}
}

func TestRead_EOFErrors(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{"empty reader", []byte{}},
		{"incomplete 2-byte", []byte{0xFD, 0x00}},
		{"incomplete 4-byte", []byte{0xFE, 0x00, 0x00}},
		{"incomplete 8-byte", []byte{0xFF, 0x00, 0x00, 0x00}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.input)
			_, err := Read(r)

			if err == nil {
				t.Errorf("expected EOF error but got none")
			}
		})
	}
}

func TestWrite_SingleByte(t *testing.T) {
	testCases := []uint64{0, 1, 127, 252}

	for _, val := range testCases {
		t.Run("value_"+string(rune(val)), func(t *testing.T) {
			var buf bytes.Buffer
			err := Write(&buf, val)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			result := buf.Bytes()
			if len(result) != 1 {
				t.Errorf("expected 1 byte, got %d", len(result))
			}
			if uint64(result[0]) != val {
				t.Errorf("expected %d, got %d", val, result[0])
			}
		})
	}
}

func TestWrite_TwoBytes(t *testing.T) {
	testCases := []uint64{253, 256, 1000, 65535}

	for _, val := range testCases {
		t.Run("value_"+string(rune(val)), func(t *testing.T) {
			var buf bytes.Buffer
			err := Write(&buf, val)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			result := buf.Bytes()
			if len(result) != 3 {
				t.Errorf("expected 3 bytes, got %d", len(result))
			}
			if result[0] != 0xFD {
				t.Errorf("expected prefix 0xFD, got 0x%02x", result[0])
			}
		})
	}
}

func TestWrite_FourBytes(t *testing.T) {
	testCases := []uint64{65536, 100000, 1000000, 4294967295}

	for _, val := range testCases {
		t.Run("value_"+string(rune(val)), func(t *testing.T) {
			var buf bytes.Buffer
			err := Write(&buf, val)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			result := buf.Bytes()
			if len(result) != 5 {
				t.Errorf("expected 5 bytes, got %d", len(result))
			}
			if result[0] != 0xFE {
				t.Errorf("expected prefix 0xFE, got 0x%02x", result[0])
			}
		})
	}
}

func TestWrite_EightBytes(t *testing.T) {
	testCases := []uint64{4294967296, 1000000000000, 18446744073709551615}

	for _, val := range testCases {
		t.Run("large_value", func(t *testing.T) {
			var buf bytes.Buffer
			err := Write(&buf, val)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			result := buf.Bytes()
			if len(result) != 9 {
				t.Errorf("expected 9 bytes, got %d", len(result))
			}
			if result[0] != 0xFF {
				t.Errorf("expected prefix 0xFF, got 0x%02x", result[0])
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	testValues := []uint64{
		0, 1, 127, 252, // single byte
		253, 256, 1000, 65535, // two bytes
		65536, 100000, 4294967295, // four bytes
		4294967296, 1000000000000, // eight bytes
	}

	for _, val := range testValues {
		t.Run("roundtrip", func(t *testing.T) {
			var buf bytes.Buffer

			// Write
			err := Write(&buf, val)
			if err != nil {
				t.Errorf("write error: %v", err)
				return
			}

			// Read back
			result, err := Read(&buf)
			if err != nil {
				t.Errorf("read error: %v", err)
				return
			}

			if result != val {
				t.Errorf("roundtrip failed: wrote %d, read %d", val, result)
			}
		})
	}
}

func TestWrite_WriterError(t *testing.T) {
	// Create a writer that always returns an error
	errorWriter := &failingWriter{}

	err := Write(errorWriter, 100)
	if err == nil {
		t.Errorf("expected writer error but got none")
	}
}

// Helper type for testing writer errors
type failingWriter struct{}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrShortWrite
}
