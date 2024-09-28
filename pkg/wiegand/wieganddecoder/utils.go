package wieganddecoder

func BitsToString(bits []bool) string {
	retval := ""
	for _, b := range bits {
		if b {
			retval += "1"
		} else {
			retval += "0"
		}
	}
	return retval
}

func BitsFromString(s string) []bool {
	bits := make([]bool, len(s))
	for i, c := range s {
		if c == '1' {
			bits[i] = true
		} else {
			bits[i] = false
		}
	}
	return bits
}
