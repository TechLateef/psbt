package types

type KV struct {
	Key   []byte
	Value []byte
}

type PSBTMap struct {
	Pairs []KV
}

type PSBT struct {
	Version uint32

	Global  PSBTMap
	Inputs  []PSBTMap
	Outputs []PSBTMap
}

const (
	PSBT_IN_NON_WITNESS_UTXO      = 0x00
	PSBT_IN_WITNESS_UTXO          = 0x01
	PSBT_IN_SIGHASH_TYPE          = 0x03
	PSBT_GLOBAL_UNSIGNED_TX       = 0x00
	PSBT_GLOBAL_INPUT_COUNT  byte = 0x02
	PSBT_GLOBAL_OUTPUT_COUNT byte = 0x03
	PSBT_GLOBAL_VERSION      byte = 0xFB
)

type UnsignedTxMeta struct {
	Version     uint32
	InputCount  uint64
	OutputCount uint64
	Locktime    uint32
}
