// GOXXTEA
// https://github.com/jbuchbinder/goxxtea
//
// vim: tabstop=4:softtabstop=4:shiftwidth=4:noexpandtab

package xxtea

// XXTEA algo cribbed from http://code.google.com/p/xxtea-algorithm/

import (
	"fmt"
)

const (
	XXTEA_DELTA = 0x9e3779b9
)

// XXTEA algorithm functions

func xxteaLongEncrypt(o []uint32, length uint32, k []uint32) []uint32 {
	v := o[:]

	p := uint32(0)
	n := length - 1
	z := v[n]
	y := v[0]
	q := uint32(6 + 52/(n+1))
	sum := uint32(0)
	if n < 1 {
		return v
	}
	for 0 < q {
		q = q - 1
		sum += XXTEA_DELTA
		e := (sum >> 2) & 3
		for p = 0; p < n; p++ {
			y = v[p+1]
			v[p] += (((z >> 5) ^ (y << 2)) + ((y >> 3) ^ (z << 4))) ^ ((sum ^ y) + (k[(p&3)^e] ^ z))
			z = v[p]
		}
		y = v[0]
		v[n] += (((z >> 5) ^ (y << 2)) + ((y >> 3) ^ (z << 4))) ^ ((sum ^ y) + (k[(p&3)^e] ^ z))
		z = v[n]
	}
	return v
}

func xxteaLongDecrypt(o []uint32, length uint32, k []uint32) []uint32 {
	v := o[:]

	p := uint32(0)
	n := length - 1
	z := v[n]
	y := v[0]
	q := uint32(6 + 52/(n+1))
	sum := uint32(q * XXTEA_DELTA)
	if n < 1 {
		return v
	}
	for sum != 0 {
		e := uint32((sum >> 2) & 3)
		for p = n; p > 0; p-- {
			z = v[p-1]
			v[p] -= (((z >> 5) ^ (y << 2)) + ((y >> 3) ^ (z << 4))) ^ ((sum ^ y) + (k[(p&3)^e] ^ z))
			y = v[p]
		}
		z = v[n]
		v[0] -= (((z >> 5) ^ (y << 2)) + ((y >> 3) ^ (z << 4))) ^ ((sum ^ y) + (k[(p&3)^e] ^ z))
		y = v[0]
		sum -= XXTEA_DELTA
	}
	return v
}

func xxteaToLongArray(data []byte, length uint32, include_length bool) (result []uint32, ret_len uint32) {
	fmt.Printf("xxteaToLongArray %s, %d\n", string(data), length)
	n := length >> 2
	if (length & 3) != 0 {
		n = n + 1
	}
	if include_length {
		result = make([]uint32, (n+1)<<2)
		result[n] = length
		ret_len = n + 1
	} else {
		result = make([]uint32, n<<2)
		ret_len = n
	}
	for iter := 0; uint32(iter) < n<<2; iter++ {
		result[iter] = 0
	}
	var i uint32
	for i = 0; i < length; i++ {
		if i < uint32(len(data)) {
			result[i>>2] |= uint32(data[i]) << ((i & 3) << 3)
			fmt.Printf("len(result) = %d, pos = %d, value = %d\n", len(result), i>>2, result[i>>2])
		}
	}
	return
}

func xxteaToByteArray(data []uint32, length uint32, include_length bool) (result []byte, ret_len uint32) {
	m := uint32(0)
	n := uint32(length << 2)
	if include_length {
		m = data[length-1]
		if (m < n-7) || (m > n-4) {
			return
		}
		n = m
	}
	result = make([]byte, n+1)
	var i uint32
	for i = 0; i < n; i++ {
		result[i] = byte((data[i>>2] >> ((i & 3) << 3)) & 0xff)
		fmt.Printf("data[%d] / byte %d = %d (orig = %d)\n", i >> 2, i, uint32(result[i]), uint32(data[i>>2]))
	}
	result[n] = byte(0)
	ret_len = n
	return
}

func XxteaEncrypt(data, key []byte) (result []byte, ret_len uint32) {
	v, v_len := xxteaToLongArray(data, uint32(len(data)), true)
	k, _ := xxteaToLongArray(key, 16, false)
	n := xxteaLongEncrypt(v, v_len, k)
	result, ret_len = xxteaToByteArray(n, v_len, false)
	return
}

func XxteaDecrypt(data, key []byte) (result []byte, ret_len uint32) {
	v, v_len := xxteaToLongArray(data, uint32(len(data)), false)
	k, _ := xxteaToLongArray(key, 16, false)
	n := xxteaLongDecrypt(v, v_len, k)
	result, ret_len = xxteaToByteArray(n, v_len, true)
	return
}
