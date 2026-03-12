package types

// OutPoint represents a reference to a previous transaction output
type OutPoint struct {
	Hash  [32]byte // Transaction hash
	Index uint32   // Output index
}

// TxInput represents a transaction input
type TxInput struct {
	PreviousOutput OutPoint // Reference to the previous transaction output
	ScriptSig      []byte   // Script signature
	Sequence       uint32   // Sequence number
}

// TxOutput represents a transaction output
type TxOutput struct {
	Value        uint64 // Amount in satoshis
	ScriptPubKey []byte // Script public key
}

type Transaction struct {
	Version  uint32
	Inputs   []TxInput
	Outputs  []TxOutput
	LockTime uint32
}
