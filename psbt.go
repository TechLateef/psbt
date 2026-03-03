package psbt

import (
	"fmt"
	"io"

	"github.com/Techlateef/psbt/internal/compactsize"
	"github.com/Techlateef/psbt/internal/parser"
	"github.com/Techlateef/psbt/types"
)

func Decode(r io.Reader) (*types.PSBT, error) {

	globalMap, err := parser.ReadPSBTHeader(r)
	if err != nil {
		return nil, err
	}
	var txByte []byte
	unsignedTxCount := 0
	var version uint64 = 0
	var versionCount int
	var inputCount, outputCount uint64
	var inputCountFound, outputCountFound bool
	var inputCountFieldCount, outputCountFieldCount int
	// var txVersionFound bool
	// var fallbackLocktimeFound bool
	for _, kv := range globalMap {
		if len(kv.Key) == 1 {
			switch kv.Key[0] {
			case types.PSBT_GLOBAL_UNSIGNED_TX:
				txByte = kv.Value
				unsignedTxCount++
			case types.PSBT_GLOBAL_VERSION:
				v, err := compactsize.FromBytes(kv.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid PSBT version encoding")
				}
				version = v
				versionCount++
			case types.PSBT_GLOBAL_INPUT_COUNT: // for v2
				inputCount, err = compactsize.FromBytes(kv.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid input count encoding")
				}
				inputCountFound = true
				inputCountFieldCount++
			case types.PSBT_GLOBAL_OUTPUT_COUNT: // for v2
				outputCount, err = compactsize.FromBytes(kv.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid output count encoding")
				}
				outputCountFound = true
				outputCountFieldCount++
			}
		} else if kv.Key[0] == types.PSBT_GLOBAL_UNSIGNED_TX {
			return nil, fmt.Errorf("invalid key for unsigned transaction: %x", kv.Key)
		}
	}
	if versionCount > 1 {
		return nil, fmt.Errorf("multiple PSBT_GLOBAL_VERSION keys")
	}
	if inputCountFieldCount > 1 {
		return nil, fmt.Errorf("multiple PSBT_GLOBAL_INPUT_COUNT keys")
	}
	if outputCountFieldCount > 1 {
		return nil, fmt.Errorf("multiple PSBT_GLOBAL_OUTPUT_COUNT keys")
	}
	if version == 0 {
		if unsignedTxCount != 1 {
			return nil, fmt.Errorf("v0 requires exactly one unsigned tx")
		}
	}

	if version == 2 {
		if unsignedTxCount != 0 {
			return nil, fmt.Errorf("v2 forbids PSBT_GLOBAL_UNSIGNED_TX")
		}

		if !inputCountFound || !outputCountFound {
			return nil, fmt.Errorf("v2 requires explicit input/output counts")
		}

	}

	switch version {
	case 0:
		if txByte == nil {
			return nil, fmt.Errorf("missing transaction bytes")
		}
		inputCount, outputCount, err = parser.CountTxInputsOutputs(txByte)
		if err != nil {
			return nil, err
		}
	case 2:
	default:
		return nil, fmt.Errorf("unsupported version")
	}
	psbt := &types.PSBT{
		Global:  types.PSBTMap{Pairs: globalMap},
		Inputs:  make([]types.PSBTMap, 0, inputCount),
		Outputs: make([]types.PSBTMap, 0, outputCount),
	}
	for i := uint64(0); i < inputCount; i++ {
		m, err := parser.ReadPSBTMap(r)
		if err != nil {
			return nil, err
		}
		psbt.Inputs = append(psbt.Inputs, types.PSBTMap{Pairs: m})
	}
	for i := uint64(0); i < outputCount; i++ {
		m, err := parser.ReadPSBTMap(r)
		if err != nil {
			return nil, err
		}
		psbt.Outputs = append(psbt.Outputs, types.PSBTMap{Pairs: m})
	}

	// Verify we read exactly the expected number of input/output maps
	if uint64(len(psbt.Inputs)) != inputCount {
		return nil, fmt.Errorf("expected %d inputs, got %d", inputCount, len(psbt.Inputs))
	}
	if uint64(len(psbt.Outputs)) != outputCount {
		return nil, fmt.Errorf("expected %d outputs, got %d", outputCount, len(psbt.Outputs))
	}

	var extra [1]byte
	if _, err := r.Read(extra[:]); err != io.EOF {
		return nil, fmt.Errorf("trailing data")
	}
	return psbt, nil
}
