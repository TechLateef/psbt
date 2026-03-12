# PSBT (Go Implementation)

A minimal, spec-focused, and educational implementation of the Bitcoin Partially Signed Bitcoin Transaction (PSBT) format written in Go.

This library aims to provide a clean and readable reference implementation of PSBT while staying close to the official Bitcoin Improvement Proposals.

## Specification Compliance

This project implements the PSBT specifications defined in:

**[BIP 174](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)** — PSBT v0 (Partially Signed Bitcoin Transactions)

**[BIP 370](https://github.com/bitcoin/bips/blob/master/bip-0370.mediawiki)** — PSBT v2 (PSBT Version 2)

These specifications define a standardized format for constructing, signing, and finalizing Bitcoin transactions across multiple participants.

## Goals

🎓 **Educational** — clear and readable implementation for learning PSBT internals

📜 **Spec-focused** — follows the BIP specifications closely

🧱 **Simple architecture** — minimal abstractions and easy to understand

🔄 **Multi-version support** — supports PSBT v0 and PSBT v2

🛡 **Safe defaults** — built-in validation and structural checks

## Current Features

### Core PSBT Support

✅ PSBT decoding from binary format

✅ Global/Input/Output map parsing

✅ Version detection (v0 / v2)

✅ Duplicate key detection

✅ CompactSize integer parsing

✅ Structural validation

### Transaction Utilities

✅ Unsigned transaction parsing (v0)

✅ Transaction reconstruction (v2)

✅ Field extraction helpers

### Developer-Friendly API

✅ `Decode()` — parse PSBT

✅ `DecodeAndValidate()` — safe decoding

✅ `GetVersion()` — detect PSBT version

✅ `IsV0()` / `IsV2()` helpers

✅ `Reconstruct()` — build transaction from PSBT

## Project Structure

```
psbt/
├── psbt.go               # PSBT decoding logic
├── reconstruct.go        # Transaction reconstruction
├── fields.go             # PSBT field extraction helpers
├── version.go            # Version helpers
├── api.go                # Public convenience APIs
├── errors.go             # Custom error types
│
├── types/
│   └── types.go          # Core PSBT + transaction structures
│
├── internal/
│   ├── parser/
│   │   └── map.go        # PSBT map parsing
│   │
│   ├── compactsize/      # Bitcoin CompactSize integers
│   │   ├── read.go
│   │   ├── write.go
│   │   └── frombytes.go
│   │
│   └── validate/
│       └── validate.go   # Structural PSBT validation
│
├── go.mod
└── README.md
```

## Installation

```bash
go get github.com/Techlateef/psbt
```

## Example Usage

### Decode and Validate a PSBT

```go
package main

import (
	"fmt"
	"os"

	"github.com/Techlateef/psbt"
)

func main() {

	file, err := os.Open("example.psbt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	psbtData, err := psbt.DecodeAndValidate(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"PSBT contains %d inputs and %d outputs\n",
		len(psbtData.Inputs),
		len(psbtData.Outputs),
	)
}
```

### Detect PSBT Version

```go
version, err := psbt.GetVersion(p)
if err != nil {
	panic(err)
}

fmt.Println("PSBT Version:", version)

isV2, _ := psbt.IsV2(p)
if isV2 {
	fmt.Println("PSBT is version 2")
}
```

### Reconstruct a Transaction (PSBT v2)

```go
tx, err := psbt.Reconstruct(p)
if err != nil {
	panic(err)
}

fmt.Println("Transaction Version:", tx.Version)
fmt.Println("Inputs:", len(tx.Inputs))
fmt.Println("Outputs:", len(tx.Outputs))
```

### Working with PSBT Data

```go
package main

import (
	"fmt"
	"os"
	
	"github.com/Techlateef/psbt"
)

func main() {
	file, err := os.Open("example.psbt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p, err := psbt.Decode(file)
	if err != nil {
		panic(err)
	}

	// Access global map data
	for _, kv := range p.Global.Pairs {
		fmt.Printf("Global key: %x, value: %x\n", kv.Key, kv.Value)
	}

	// Access input maps
	for i, input := range p.Inputs {
		fmt.Printf("Input %d has %d key-value pairs\n", i, len(input.Pairs))
	}

	// Access output maps
	for i, output := range p.Outputs {
		fmt.Printf("Output %d has %d key-value pairs\n", i, len(output.Pairs))
	}
}
```

## Core Types

```go
type PSBT struct {
	Global  PSBTMap
	Inputs  []PSBTMap
	Outputs []PSBTMap
}

type PSBTMap struct {
	Pairs []KV
}

type KV struct {
	Key   []byte
	Value []byte
}
```

Transaction types:

```go
type Transaction struct {
	Version  uint32
	Inputs   []TxInput
	Outputs  []TxOutput
	LockTime uint32
}

type TxInput struct {
	PreviousOutput OutPoint
	ScriptSig      []byte
	Sequence       uint32
}

type TxOutput struct {
	Value        uint64
	ScriptPubKey []byte
}

type OutPoint struct {
	Hash  [32]byte
	Index uint32
}
```

## Validation

The validation system checks:

- Duplicate keys in PSBT maps
- Proper map structure
- PSBT magic bytes (`psbt\xff`)
- Correct CompactSize encodings
- Required fields for PSBT v0 and v2

## Status

🚧 **Active Development**

Current focus areas:
- improving validation
- extending PSBT utilities
- expanding test coverage

## Roadmap

### Near Term

- [ ] PSBT serialization (encoding)
- [ ] PSBT combining
- [ ] Improved validation rules
- [ ] Comprehensive test suite

### Future Features

- [ ] Transaction signing
- [ ] PSBT finalization
- [ ] Hardware wallet compatibility
- [ ] CLI utilities

## Contributing

Contributions are welcome. This project values:

- spec-accurate implementations
- readable and well-structured code
- meaningful validation checks
- thorough testing

## License

MIT License

## References

- [BIP 174](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki) — Partially Signed Bitcoin Transactions
- [BIP 370](https://github.com/bitcoin/bips/blob/master/bip-0370.mediawiki) — PSBT Version 2