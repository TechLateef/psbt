package psbt

import (
	"fmt"
	"io"

	"github.com/Techlateef/psbt/internal/validate"
	"github.com/Techlateef/psbt/types"
)

func DecodeAndValidate(r io.Reader) (*types.PSBT, error) {
	p, err := Decode(r)
	if err != nil {
		return nil, err
	}
	if err := validate.Validate(p); err != nil {
		return nil, fmt.Errorf("PSBT validation failed: %w", err)
	}
	return p, nil
}
