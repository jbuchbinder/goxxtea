// GOXXTEA
// https://github.com/jbuchbinder/goxxtea
//
// vim: tabstop=4:softtabstop=4:shiftwidth=4:noexpandtab

package xxtea

// XXTEA algo cribbed from http://code.google.com/p/xxtea-algorithm/

const (
	delta = 0x9E3779B9
)

// XXTEA algorithm functions

type XXTEA struct{}

func (x XXTEA) toBytes(v []uint32, includeLength bool) []byte {
	length := uint32(len(v))
	n := length << 2
	if includeLength {
		m := v[length-1]
		n -= 4
		if (m < n-3) || (m > n) {
			return nil
		}
		n = m
	}
	bytes := make([]byte, n)
	for i := uint32(0); i < n; i++ {
		bytes[i] = byte(v[i>>2] >> ((i & 3) << 3))
	}
	return bytes
}

func (x XXTEA) toUint32s(bytes []byte, includeLength bool) (v []uint32) {
	length := uint32(len(bytes))
	n := length >> 2
	if length&3 != 0 {
		n++
	}
	if includeLength {
		v = make([]uint32, n+1)
		v[n] = length
	} else {
		v = make([]uint32, n)
	}
	for i := range length {
		v[i>>2] |= uint32(bytes[i]) << ((i & 3) << 3)
	}
	return v
}

func (x XXTEA) mx(sum uint32, y uint32, z uint32, p uint32, e uint32, k []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (k[p&3^e] ^ z))
}

func (x XXTEA) fixk(k []uint32) []uint32 {
	if len(k) < 4 {
		key := make([]uint32, 4)
		copy(key, k)
		return key
	}
	return k
}

func (x XXTEA) encrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = x.fixk(k)
	var y, z, sum, e, p, q uint32
	z = v[n]
	sum = 0
	for q = 6 + 52/length; q > 0; q-- {
		sum += delta
		e = sum >> 2 & 3
		for p = 0; p < n; p++ {
			y = v[p+1]
			v[p] += x.mx(sum, y, z, p, e, k)
			z = v[p]
		}
		y = v[0]
		v[n] += x.mx(sum, y, z, p, e, k)
		z = v[n]
	}
	return v
}

func (x XXTEA) decrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = x.fixk(k)
	var y, z, sum, e, p, q uint32
	y = v[0]
	q = 6 + 52/length
	for sum = q * delta; sum != 0; sum -= delta {
		e = sum >> 2 & 3
		for p = n; p > 0; p-- {
			z = v[p-1]
			v[p] -= x.mx(sum, y, z, p, e, k)
			y = v[p]
		}
		z = v[n]
		v[0] -= x.mx(sum, y, z, p, e, k)
		y = v[0]
	}
	return v
}

// Encrypt the data with key.
// data is the bytes to be encrypted.
// key is the encrypt key. It is the same as the decrypt key.
func (x XXTEA) Encrypt(data []byte, key []byte) []byte {
	if len(data) == 0 {
		return data
	}
	return x.toBytes(x.encrypt(x.toUint32s(data, true), x.toUint32s(key, false)), false)
}

// Decrypt the data with key.
// data is the bytes to be decrypted.
// key is the decrypted key. It is the same as the encrypt key.
func (x XXTEA) Decrypt(data []byte, key []byte) []byte {
	if len(data) == 0 {
		return data
	}
	return x.toBytes(x.decrypt(x.toUint32s(data, false), x.toUint32s(key, false)), true)
}

// Convenience functions

func Encrypt(data []byte, key []byte) []byte {
	return XXTEA{}.Encrypt(data, key)
}

func Decrypt(data []byte, key []byte) []byte {
	return XXTEA{}.Decrypt(data, key)
}
