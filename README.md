# PSBT (Go Implementation)

A minimal, spec-focused, and educational implementation of the Bitcoin Partially Signed Bitcoin Transaction (PSBT) format in Go.

## Specification Compliance

This project implements the PSBT specifications defined in:
- **[BIP 174](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)** - PSBT v0 (Partially Signed Bitcoin Transactions)
- **[BIP 370](https://github.com/bitcoin/bips/blob/master/bip-0370.mediawiki)** - PSBT v2 (Version 2 Partially Signed Bitcoin Transactions)

## Goals

- 🎓 **Educational**: Provide a clean and readable PSBT implementation for learning
- 📋 **Spec-focused**: Follow the BIP specifications closely without unnecessary abstractions  
- 🏗️ **Simple architecture**: Keep the codebase maintainable and understandable
- 🔄 **Multi-version support**: Handle both PSBT v0 and PSBT v2 formats
- ✅ **Robust validation**: Comprehensive error checking and input validation

## Current Features

- ✅ **PSBT decoding** from binary format
- ✅ **Global/Input/Output map parsing** with proper key-value handling
- ✅ **Version detection** (automatically detects v0/v2 format)
- ✅ **Structural validation** (magic bytes, format compliance)
- ✅ **Duplicate key detection** across all maps
- ✅ **CompactSize integer** encoding/decoding
- ✅ **Unsigned transaction parsing** and metadata extraction

## Project Structure

```
psbt/
├── psbt.go                    # Main PSBT decoding functionality
├── errors.go                 # Custom error types
├── types/
│   └── types.go              # Core PSBT data structures
├── internal/
│   ├── parser/
│   │   └── map.go           # PSBT map and header parsing
│   ├── compactsize/         # Bitcoin CompactSize integer handling
│   │   ├── read.go
│   │   ├── write.go
│   │   └── frombyte.go
│   ├── validate/
│   │   └── validate.go      # PSBT validation logic
│   ├── keys/                # Key handling utilities
│   └── serialize/           # Serialization utilities
├── go.mod
└── README.md
```

## Installation

```bash
go get github.com/Techlateef/psbt
```

## Example Usage

### Basic PSBT Decoding and Validation

```go
package main

import (
    "os"
    "fmt"
    
    "github.com/Techlateef/psbt"
    "github.com/Techlateef/psbt/internal/validate" 
)

func main() {
    // Open PSBT file
    file, err := os.Open("example.psbt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Decode PSBT
    p, err := psbt.Decode(file)
    if err != nil {
        panic(err)
    }

    // Validate PSBT structure
    err = validate.Validate(p)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Successfully decoded PSBT with %d inputs and %d outputs\n", 
        len(p.Inputs), len(p.Outputs))
}
```

### Working with PSBT Data

```go
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
```

## API Reference

### Main Functions

- `psbt.Decode(r io.Reader) (*types.PSBT, error)` - Decode PSBT from binary format
- `validate.Validate(p *types.PSBT) error` - Validate PSBT structure and detect duplicates

### Core Types

```go
type PSBT struct {
    Version uint32
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

## Validation Features

The validation system checks for:
- **Duplicate keys** within the same map (forbidden by spec)
- **Proper key ordering** (keys must be in strictly increasing order)
- **Valid PSBT magic bytes** (`psbt` + separator)
- **Correct CompactSize encoding** for integers
- **Segwit transaction restrictions** in unsigned transactions

## Status

🚧 **Work in progress** - This implementation is actively being developed.

## Development Roadmap

### 🚧 Work in Progress
- [ ] **Transaction reconstruction** from PSBT data
- [ ] **Signing support** with private keys
- [ ] **PSBT finalization** (combining signatures)
- [ ] **PSBT serialization** (encoding back to binary)

### 🔮 Future Enhancements
- [ ] **Hardware wallet integration**
- [ ] **Advanced validation rules**
- [ ] **PSBT combining/merging**
- [ ] **Comprehensive test suite**
- [ ] **CLI tools**

## Contributing

This project is designed to be educational and welcomes contributions that:
- Follow the PSBT specifications closely
- Include comprehensive tests
- Maintain code readability and documentation
- Add meaningful validation or functionality

## License

MIT License - see [LICENSE](LICENSE) file for details.

## References

- [BIP 174 - PSBT v0](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [BIP 370 - PSBT v2](https://github.com/bitcoin/bips/blob/master/bip-0370.mediawiki)
- [Bitcoin Core implementation](https://github.com/bitcoin/bitcoin/blob/master/src/psbt.h)