package net

type HardwareAddress interface {
	Bytes() []byte
	Len() uint8
	String() string
}

// e.g. IP address
type ProtocolAddress interface {
	Bytes() []byte
	Len() uint8
	String() string
	IsEmpty() bool
}
