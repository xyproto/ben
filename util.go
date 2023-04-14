package main

func uint16ToBytes(value uint16) []byte {
	return []byte{byte(value >> 8), byte(value & 0xFF)}
}

func uint32ToBytes(value uint32) []byte {
	return []byte{byte(value >> 24), byte((value >> 16) & 0xFF), byte((value >> 8) & 0xFF), byte(value & 0xFF)}
}
