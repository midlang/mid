package proto

type Message interface {
	MessageName() string
	Encode() ([]byte, error)
	Decode([]byte) (int, error)
}
