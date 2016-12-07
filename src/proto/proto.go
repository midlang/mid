package proto

type Message interface {
	MessageName() string
	Encode(Writer) error
	Decode([]byte) (int, error)
}
