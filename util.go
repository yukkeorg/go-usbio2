package usbio2


func boolarray_to_byte(ba []bool) byte {
	var b byte = 0
	for i, p := range ba {
		if p {
			b |= (1 << uint(i))
		}
	}
	return b
}

func byte_to_boolarray(b byte) []bool {
	ba := make([]bool, 8, 8)
	for i := 0; b != 0; i++ {
		if (b & 0x01) == 0x01 {
			ba[i] = true
		} else {
			ba[i] = false
		}
		b = b >> 1
	}
	return ba
}
