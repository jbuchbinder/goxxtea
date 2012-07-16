package xxtea

func Encrypt(data []byte, key []byte) []byte {
	if len(data) == 0 {
		return data
	}
	return toByteArray(EncryptInt(toIntArray(data, true), toIntArray(key, false)), false)
}

func Decrypt(data []byte, key []byte) []byte {
	if len(data) == 0 {
		return data
	}
	return toByteArray(DecryptInt(toIntArray(data, false), toIntArray(key, false)), true)
}

func EncryptInt(v []int, k []int) []int {
	n := len(v) - 1

	if n < 1 {
		return v
	}
	if len(k) < 4 {
		key := make([]int, 4)
		key = k[0:len(k)]
		k = key
	}
	z := v[n]
	y := v[0]
	delta := 0x9E3779B9
	sum := 0
	e := 0
	p := 0
	q := 6 + 52/(n+1)

	for q-1 > 0 {
		q = q - 1
		sum = sum + delta
		e = sum >> 2 & 3
		for p = 0; p < n; p++ {
			y = v[p+1]
			v[p] += (z>>5 ^ y<<2) + (y>>3 ^ z<<4) ^ (sum ^ y) + (k[p&3^e] ^ z)
			z = v[p]
		}
		y = v[0]
		v[n] += (z>>5 ^ y<<2) + (y>>3 ^ z<<4) ^ (sum ^ y) + (k[p&3^e] ^ z)
		z = v[n]
	}
	return v
}

func DecryptInt(v []int, k []int) []int {
	n := len(v) - 1

	if n < 1 {
		return v
	}
	if len(k) < 4 {
		key := make([]int, 4)
		key = k[0:len(k)]
		k = key
	}
	z := v[n]
	y := v[0]
	delta := 0x9E3779B9
	sum := 0
	e := 0
	p := 0
	q := 6 + 52/(n+1)

	sum = q * delta
	for sum != 0 {
		e = sum >> 2 & 3
		for p = n; p > 0; p-- {
			z = v[p-1]
			v[p] -= (z>>5 ^ y<<2) + (y>>3 ^ z<<4) ^ (sum ^ y) + (k[p&3^e] ^ z)
			y = v[p]
		}
		z = v[n]
		v[0] -= (z>>5 ^ y<<2) + (y>>3 ^ z<<4) ^ (sum ^ y) + (k[p&3^e] ^ z)
		y = v[0]
		sum = sum - delta
	}
	return v
}

func toIntArray(data []byte, includeLength bool) []int {
	n := 0
	if (len(data) & 3) == 0 {
		n = (len(data) >> 2)
	} else {
		n = ((len(data) >> 2) + 1)
	}
	l := n
	if includeLength {
		l = n + 1
	}
	result := make([]int, l)
	if includeLength {
		result[n] = len(data)
	}
	n = len(data)
	for i := 0; i < n; i++ {
		result[i>>2] |= (0x000000ff & data[i]) << ((i & 3) << 3)
	}
	return result
}

func toByteArray(data []int, includeLength bool) []byte {
	n := len(data) << 2

	if includeLength {
		m := data[len(data)-1]

		if m > n {
			return nil
		} else {
			n = m
		}
	}
	result := make([]byte, n)

	for i := 0; i < n; i++ {
		result[i] = byte((data[i>>2] >> uint((i & 3) << 3)) & 0xff)
	}
	return result
}
