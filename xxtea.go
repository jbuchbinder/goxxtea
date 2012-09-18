// GOXXTEA
// https://github.com/jbuchbinder/goxxtea

package xxtea

// XXTEA algo cribbed from http://code.google.com/p/xxtea-algorithm/

// #cgo LDFLAGS: -lc
/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <sys/types.h>
#include <sys/stat.h>

typedef unsigned int xxtea_long;

#define XXTEA_MX (( (z >> 5) ^ (y << 2)) + ((y >> 3) ^ (z << 4))) ^ ((sum ^ y) + (k[(p & 3) ^ e] ^ z))
#define XXTEA_DELTA 0x9e3779b9

unsigned char *char_to_unsigned_array(char *orig) {
  char *p = orig;
  while (p != '\0') {
    *p = (unsigned char) *p;
    p++;
  }
  return (unsigned char *) orig;
}

char *unsigned_to_char_array(unsigned char *orig) {
  unsigned char *p = orig;
  while (p != '\0') {
    *p = (char) *p;
    p++;
  }
  return (char *) orig;
}

void xxtea_long_encrypt(xxtea_long *v, xxtea_long len, xxtea_long *k) {
    xxtea_long n = len - 1;
    xxtea_long z = v[n], y = v[0], p, q = 6 + 52 / (n + 1), sum = 0, e;
    if (n < 1) {
        return;
    }
    while (0 < q--) {
        sum += XXTEA_DELTA;
        e = sum >> 2 & 3;
        for (p = 0; p < n; p++) {
            y = v[p + 1];
            z = v[p] += XXTEA_MX;
        }
        y = v[0];
        z = v[n] += XXTEA_MX;
    }
}

void xxtea_long_decrypt(xxtea_long *v, xxtea_long len, xxtea_long *k) {
    xxtea_long n = len - 1;
    xxtea_long z = v[n], y = v[0], p, q = 6 + 52 / (n + 1), sum = q * XXTEA_DELTA, e;
    if (n < 1) {
        return;
    }
    while (sum != 0) {
        e = sum >> 2 & 3;
        for (p = n; p > 0; p--) {
            z = v[p - 1];
            y = v[p] -= XXTEA_MX;
        }
        z = v[n];
        y = v[0] -= XXTEA_MX;
        sum -= XXTEA_DELTA;
    }
}

static xxtea_long *xxtea_to_long_array(const unsigned char *data, xxtea_long len, int include_length, xxtea_long *ret_len) {
    xxtea_long i, n, *result;
        n = len >> 2;
    n = (((len & 3) == 0) ? n : n + 1);
    if (include_length) {
        result = (xxtea_long *)malloc((n + 1) << 2);
        result[n] = len;
                *ret_len = n + 1;
        } else {
        result = (xxtea_long *)malloc(n << 2);
                *ret_len = n;
    }
        memset(result, 0, n << 2);
        for (i = 0; i < len; i++) {
        result[i >> 2] |= (xxtea_long)data[i] << ((i & 3) << 3);
    }
    return result;
}

static unsigned char *xxtea_to_byte_array(xxtea_long *data, xxtea_long len, int include_length, xxtea_long *ret_len) {
    xxtea_long i, n, m;
    unsigned char *result;
    n = len << 2;
    if (include_length) {
        m = data[len - 1];
        if ((m < n - 7) || (m > n - 4)) return NULL;
        n = m;
    }
    result = (unsigned char *)malloc(n + 1);
        for (i = 0; i < n; i++) {
        result[i] = (unsigned char)((data[i >> 2] >> ((i & 3) << 3)) & 0xff);
    }
        result[n] = '\0';
        *ret_len = n;
        return result;
}

unsigned char *xxtea_encrypt(const unsigned char *data, xxtea_long len, unsigned char *key, xxtea_long *ret_len)
{
    unsigned char *result;
    xxtea_long *v, *k, v_len, k_len;
    v = xxtea_to_long_array(data, len, 1, &v_len);
    k = xxtea_to_long_array(key, 16, 0, &k_len);
    xxtea_long_encrypt(v, v_len, k);
    result = xxtea_to_byte_array(v, v_len, 0, ret_len);
    free(v);
    free(k);
    return result;
}

unsigned char *xxtea_decrypt(const unsigned char *data, xxtea_long len, unsigned char *key, xxtea_long *ret_len)
{
    unsigned char *result;
    xxtea_long *v, *k, v_len, k_len;
    v = xxtea_to_long_array(data, len, 0, &v_len);
    k = xxtea_to_long_array(key, 16, 0, &k_len);
    xxtea_long_decrypt(v, v_len, k);
    result = xxtea_to_byte_array(v, v_len, 1, ret_len);
    free(v);
    free(k);
    return result;
}
*/
import "C"

import (
	"unsafe"
)

// Encrypts a string using XXTEA with a specified key.
func Encrypt(text, key string) (out string, err error) {
	ctext := C.char_to_unsigned_array(C.CString(text))
	defer C.free(unsafe.Pointer(ctext))
	ckey := C.char_to_unsigned_array(C.CString(key))
	defer C.free(unsafe.Pointer(ckey))
	var cretlen C.xxtea_long

	out = C.GoString(C.unsigned_to_char_array(C.xxtea_encrypt(ctext, C.xxtea_long(len(text)), ckey, &cretlen)))
	return
}

// Decrypts a string using XXTEA with a specified key.
func Decrypt(text, key string) (out string, err error) {
	ctext := C.char_to_unsigned_array(C.CString(text))
	defer C.free(unsafe.Pointer(ctext))
	ckey := C.char_to_unsigned_array(C.CString(key))
	defer C.free(unsafe.Pointer(ckey))
	var cretlen C.xxtea_long

	out = C.GoString(C.unsigned_to_char_array(C.xxtea_decrypt(ctext, C.xxtea_long(len(text)), ckey, &cretlen)))
	return
}
