package psbt

import (
	"fmt"

	"github.com/Techlateef/psbt/internal/compactsize"
	"github.com/Techlateef/psbt/types"
)

func GetVersion(p *types.PSBT) (uint64, error) {
	val, ok := getField(p.Global, types.PSBT_GLOBAL_VERSION)
	if !ok {
		return 0, nil // default version is 0 if not specified
	}
	v, err := compactsize.FromBytes(val)
	if err != nil {
		return 0, fmt.Errorf("invalid PSBT version encoding")
	}
	return v, nil
}

func IsV0(p *types.PSBT) (bool, error) {
	version, err := GetVersion(p)
	if err != nil {
		return false, err
	}
	return version == 0, nil
}

func IsV2(p *types.PSBT) (bool, error) {
	version, err := GetVersion(p)
	if err != nil {
		return false, err
	}
	return version == 2, nil
}
