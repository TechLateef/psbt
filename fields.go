package psbt

import (
	"encoding/binary"
	"fmt"

	"github.com/Techlateef/psbt/internal/compactsize"
	"github.com/Techlateef/psbt/types"
)

func getField(m types.PSBTMap, keyType byte) ([]byte, bool) {
	for _, kv := range m.Pairs {
		if len(kv.Key) == 1 && kv.Key[0] == keyType {
			return kv.Value, true
		}
	}
	return nil, false
}

func GetTxVersion(p types.PSBT) (uint32, error) {
	val, ok := getField(p.Global, types.PSBT_GLOBAL_TX_VERSION)
	if !ok {
		return 0, fmt.Errorf("missing tx version")
	}
	v, err := compactsize.FromBytes(val)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func GetFallbackLocktime(p types.PSBT) (uint32, error) {
	val, ok := getField(p.Global, types.PSBT_GLOBAL_FALLBACK_LOCKTIME)
	if !ok {
		return 0, fmt.Errorf("missing fallback locktime")
	}
	v, err := compactsize.FromBytes(val)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func GetPrevTxID(p types.PSBTMap) ([]byte, error) {
	val, ok := getField(p, types.PSBT_IN_PREVIOUS_TXID)
	if !ok {
		return nil, fmt.Errorf("missing previous txid")
	}
	if len(val) != 32 {
		return nil, fmt.Errorf("invalid previous txid length")
	}
	return val, nil
}

func GetOutputIndex(p types.PSBTMap) (uint32, error) {
	val, ok := getField(p, types.PSBT_IN_OUTPUT_INDEX)
	if !ok {
		return 0, fmt.Errorf("missing output index")
	}
	if len(val) != 4 {
		return 0, fmt.Errorf("invalid output index length")
	}
	return binary.LittleEndian.Uint32(val), nil
}

func GetSequence(m types.PSBTMap) (uint32, error) {
	val, ok := getField(m, types.PSBT_IN_SEQUENCE)
	if !ok {
		return 0, fmt.Errorf("missing sequence")
	}

	if len(val) != 4 {
		return 0, fmt.Errorf("invalid sequence length")
	}

	return binary.LittleEndian.Uint32(val), nil
}

func GetOutputAmount(m types.PSBTMap) (uint64, error) {
	val, ok := getField(m, types.PSBT_OUT_AMOUNT)
	if !ok {
		return 0, fmt.Errorf("missing output amount")
	}

	if len(val) != 8 {
		return 0, fmt.Errorf("invalid output amount length")
	}

	return binary.LittleEndian.Uint64(val), nil
}

func GetOutputScript(output types.PSBTMap) ([]byte, error) {
	val, ok := getField(output, types.PSBT_OUT_SCRIPT)
	if !ok {
		return nil, fmt.Errorf("missing output script")
	}
	return val, nil
}
