package ble

func ReadUint16LE(bytes []byte, offset int) uint16 {
	return uint16(bytes[offset]) | uint16(bytes[offset+1])<<8
}

func ReadUint8(bytes []byte, offset int) uint8 {
	return uint8(bytes[offset])
}

func ReadInt16LE(bytes []byte, offset int) int16 {
	return int16(bytes[offset]) | int16(bytes[offset+1])<<8
}

func ReadUInt32BE(bytes []byte, offset int) uint32 {
	return uint32(bytes[offset])<<24 | uint32(bytes[offset+1])<<16 | uint32(bytes[offset+2])<<8 | uint32(bytes[offset+3])
}

func ReadInt8(bytes []byte, offset int) int8 {
	return int8(bytes[offset])
}

func ReadUint16BE(bytes []byte, offset int) uint16 {
	return uint16(bytes[offset])<<8 | uint16(bytes[offset+1])
}

func Reverse(input []byte) []byte {
	inputLength := len(input)
	output := make([]byte, inputLength)

	for i := 0; i < inputLength; i++ {
		output[inputLength-i-1] = input[i]
	}
	return output
}
