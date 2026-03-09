package validate

import (
	"fmt"

	"github.com/Techlateef/psbt/types"
)

func Validate(p *types.PSBT) error {
	if err := validateMap(p.Global); err != nil {
		return fmt.Errorf("global map validation failed: %w", err)
	}

	for i, in := range p.Inputs {
		if err := validateInput(in); err != nil {
			return fmt.Errorf("input %d validation failed: %w", i, err)
		}
	}

	for i, out := range p.Outputs {
		if err := validateOutput(out); err != nil {
			return fmt.Errorf("output %d validation failed: %w", i, err)
		}
	}

	return nil
}

func validateMap(pmap types.PSBTMap) error {
	seen := make(map[string]bool)

	for _, kv := range pmap.Pairs {
		keyStr := string(kv.Key)
		if seen[keyStr] {
			return fmt.Errorf("duplicate key in PSBT map: %x", kv.Key)
		}
		seen[keyStr] = true
	}
	return nil

}

func validateInput(input types.PSBTMap) error {
	if err := validateMap(input); err != nil {
		return fmt.Errorf("input map validation error: %w", err)
	}
	return nil
}

func validateOutput(output types.PSBTMap) error {
	if err := validateMap(output); err != nil {
		return fmt.Errorf("output map validation error: %w", err)
	}
	return nil
}
