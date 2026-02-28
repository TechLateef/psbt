type Key []byte
type Value []byte

type Map map[string]Value

type PSBT struct {
	Version uint32

	Global  Map
	Inputs  []Map
	Outputs []Map
}
