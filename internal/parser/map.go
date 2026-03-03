package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Techlateef/psbt/internal/compactsize"
	"github.com/Techlateef/psbt/types"
)

func ReadPSBTMap(r io.Reader) ([]types.KV, error) {

	result := make([]types.KV, 0)
	var prevKey []byte
	for {

		keyLen, err := compactsize.Read(r)
		if err != nil {
			return nil, err
		}
		if keyLen == 0 {
			break
		}

		key := make([]byte, keyLen)
		if _, err := io.ReadFull(r, key); err != nil {
			return nil, err
		}
		if prevKey != nil {
			cmp := bytes.Compare(prevKey, key)
			if cmp >= 0 {
				return nil, fmt.Errorf("psbt keys not strictly increasing")
			}
		}
		prevKey = append([]byte(nil), key...)

		valueLen, err := compactsize.Read(r)
		if err != nil {
			return nil, err
		}

		value := make([]byte, valueLen)
		if _, err := io.ReadFull(r, value); err != nil {
			return nil, err
		}

		result = append(result, types.KV{Key: key, Value: value})
	}
	return result, nil
}

func ReadPSBTHeader(r io.Reader) ([]types.KV, error) {
	var magic [5]byte
	_, err := io.ReadFull(r, magic[:])
	if err != nil {
		return nil, err
	}
	if magic != [5]byte{0x70, 0x73, 0x62, 0x74, 0xff} {
		return nil, fmt.Errorf("invalid psbt magic")
	}
	globalMap, err := ReadPSBTMap(r)
	if err != nil {
		return nil, err
	}
	return globalMap, nil
}

func CountTxInputsOutputs(txBytes []byte) (uint64, uint64, error) {
	// Create a bytes.Reader from txBytes
	reader := bytes.NewReader(txBytes)
	// skip the first 4 byte which is the version
	_, err := reader.Seek(4, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}
	// Read Compactsize the input countTxInputsOutput
	inputCount, err := compactsize.Read(reader)
	if err != nil {
		return 0, 0, err
	}

	if inputCount == 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		if b == 0x01 {
			return 0, 0, fmt.Errorf("segwit serializationnot allowed in PSBT unsigned tx")
		}
		// Not segwith-> unread the byte
		if err := reader.UnreadByte(); err != nil {
			return 0, 0, err
		}
	}
	// Skip over each Input to reach the Output countTxInputsOutputs
	for range inputCount {
		// Skip Output: 32 bytes (TXID) + 4 bytes (VOUT) = 36 bytes
		_, err := reader.Seek(36, io.SeekCurrent)
		if err != nil {
			return 0, 0, err
		}
		// Read the script length to know how much to skip
		scriptLen, err := compactsize.Read(reader)
		if err != nil {
			return 0, 0, err
		}
		if scriptLen != 0 {
			return 0, 0, fmt.Errorf("PSBT with non-empty scriptSig is not supported")
		}
		// Skip the Scrip + 4 byte for sequece
		_, err = reader.Seek(int64(scriptLen)+4, io.SeekCurrent)
		if err != nil {
			return 0, 0, err
		}
	}
	outputCount, err := compactsize.Read(reader)
	if err != nil {
		return 0, 0, err
	}
	for range outputCount {
		// Skip Output: 8 bytes (value)
		_, err := reader.Seek(8, io.SeekCurrent)
		if err != nil {
			return 0, 0, err
		}
		// Read the script length to know how much to skip
		scriptLen, err := compactsize.Read(reader)
		if err != nil {
			return 0, 0, err
		}
		// Skip the Scrip
		_, err = reader.Seek(int64(scriptLen), io.SeekCurrent)
		if err != nil {
			return 0, 0, err
		}

	}
	// Skip the locktime
	_, err = reader.Seek(4, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}
	if reader.Len() != 0 {
		return 0, 0, fmt.Errorf("unexpected trailing data in unsigned transaction")
	}

	return inputCount, outputCount, nil

}

func ParseUnsignedTx(txBytes []byte) (*types.UnsignedTxMeta, error) {
	reader := bytes.NewReader(txBytes)

	// 1. Version
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil {
		return nil, err
	}

	// 2. Input count
	inputCount, err := compactsize.Read(reader)
	if err != nil {
		return nil, err
	}

	// 3. Segwit check (if inputCount == 0)
	if inputCount == 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 0x01 {
			return nil, fmt.Errorf("segwit serialization not allowed in PSBT unsigned tx")
		}
		if err := reader.UnreadByte(); err != nil {
			return nil, err
		}
	}

	// 4. Skip inputs
	for range inputCount {
		if _, err := reader.Seek(36, io.SeekCurrent); err != nil { // TXID + VOUT
			return nil, err
		}
		scriptLen, err := compactsize.Read(reader)
		if err != nil {
			return nil, err
		}
		if scriptLen != 0 {
			return nil, fmt.Errorf("PSBT with non-empty scriptSig is not supported")
		}
		if _, err := reader.Seek(int64(scriptLen)+4, io.SeekCurrent); err != nil { // script + sequence
			return nil, err
		}
	}

	// 5. Output count
	outputCount, err := compactsize.Read(reader)
	if err != nil {
		return nil, err
	}

	// 6. Skip outputs
	for range outputCount {
		if _, err := reader.Seek(8, io.SeekCurrent); err != nil { // value
			return nil, err
		}
		scriptLen, err := compactsize.Read(reader)
		if err != nil {
			return nil, err
		}
		if _, err := reader.Seek(int64(scriptLen), io.SeekCurrent); err != nil { // script
			return nil, err
		}
	}

	// 7. Locktime
	var locktime uint32
	if err := binary.Read(reader, binary.LittleEndian, &locktime); err != nil {
		return nil, err
	}

	// 8. Trailing data check
	if reader.Len() != 0 {
		return nil, fmt.Errorf("unexpected trailing data in unsigned transaction")
	}

	// 9. Return meta
	return &types.UnsignedTxMeta{
		Version:     version,
		InputCount:  inputCount,
		OutputCount: outputCount,
		Locktime:    locktime,
	}, nil
}
