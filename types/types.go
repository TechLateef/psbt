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
	PSBT_IN_NON_WITNESS_UTXO           = 0x00
	PSBT_IN_WITNESS_UTXO               = 0x01
	PSBT_IN_SIGHASH_TYPE               = 0x03
	PSBT_GLOBAL_UNSIGNED_TX            = 0x00
	PSBT_GLOBAL_INPUT_COUNT       byte = 0x04
	PSBT_GLOBAL_OUTPUT_COUNT      byte = 0x05
	PSBT_GLOBAL_VERSION           byte = 0xFB
	PSBT_GLOBAL_FALLBACK_LOCKTIME      = 0x03
	PSBT_GLOBAL_TX_VERSION             = 0x02
	PSBT_IN_OUTPUT_INDEX               = 0x0F
	PSBT_IN_PREVIOUS_TXID              = 0x0E
	PSBT_IN_SEQUENCE                   = 0x10
	PSBT_OUT_AMOUNT                    = 0x00
	PSBT_OUT_SCRIPT                    = 0x01
)

type UnsignedTxMeta struct {
	Version     uint32
	InputCount  uint64
	OutputCount uint64
	Locktime    uint32
}
