package psbt

import (
	"fmt"

	"github.com/Techlateef/psbt/types"
)

func Reconstruct(p *types.PSBT) (*types.Transaction, error) {

	version, err := GetTxVersion(*p)
	if err != nil {
		return nil, err
	}

	fallbackLocktime, err := GetFallbackLocktime(*p)
	if err != nil {
		return nil, err
	}

	tx := &types.Transaction{
		Version:  version,
		Inputs:   make([]types.TxInput, 0, len(p.Inputs)),
		Outputs:  make([]types.TxOutput, 0, len(p.Outputs)),
		LockTime: fallbackLocktime,
	}
	// ----- Inputs ----- //
	for _, in := range p.Inputs {
		txidBytes, err := GetPrevTxID(in)
		if err != nil {
			return nil, err
		}
		if len(txidBytes) != 32 {
			return nil, fmt.Errorf("invalid txid length")
		}
		var hash [32]byte
		copy(hash[:], txidBytes)
		index, err := GetOutputIndex(in)
		if err != nil {
			return nil, err
		}
		sequence, err := GetSequence(in)
		if err != nil {
			return nil, err
		}

		tx.Inputs = append(tx.Inputs, types.TxInput{
			PreviousOutput: types.OutPoint{
				Hash:  hash,
				Index: index,
			},
			ScriptSig: nil, // unsigned PSBT doesn't have scriptSig
			Sequence:  sequence,
		})
	}

	// ----- Outputs ----- //
	for _, out := range p.Outputs {
		amount, err := GetOutputAmount(out)
		if err != nil {
			return nil, err
		}
		script, err := GetOutputScript(out)
		if err != nil {
			return nil, err
		}
		tx.Outputs = append(tx.Outputs, types.TxOutput{
			Value:        amount,
			ScriptPubKey: script,
		})
	}
	return tx, nil
}
