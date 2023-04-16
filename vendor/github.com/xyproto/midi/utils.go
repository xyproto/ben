package midi

import (
	"io"
)

func uint16ToBytes(value uint16) []byte {
	return []byte{byte(value >> 8), byte(value & 0xFF)}
}

func uint32ToBytes(value uint32) []byte {
	return []byte{byte(value >> 24), byte((value >> 16) & 0xFF), byte((value >> 8) & 0xFF), byte(value & 0xFF)}
}

func writeBytes(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}

func writeUint32(w io.Writer, value uint32) error {
	return writeBytes(w, uint32ToBytes(value))
}
