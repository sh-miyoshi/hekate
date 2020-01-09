package util

// Int2bytes convert uint64 to an array of bytes
func Int2bytes(v uint64) []byte {
	if v == 0 {
		return []byte{0}
	}

	bytes := make([]byte, 8)
	i := 0
	for ; i < 8; i++ {
		if v <= 0 {
			break
		}
		bytes[7-i] = byte(v & 255)
		v >>= 8
	}
	return bytes[8-i : 8]
}
