// GOXXTEA
// https://github.com/jbuchbinder/goxxtea

package xxtea

import (
	"testing"
)

func Test_XXTEA_Roundtrip(t *testing.T) {
	key := "xxteaTEST"
	orig := "01234567890123456789"

	intermediate, _ := Encrypt(orig, key)

	decoded, _ := Decrypt(intermediate, key)

	if decoded != orig {
		t.Error("Round trip encode failed, result = " + decoded)
	} else {
		t.Log("Round trip xxtea encode passed")
	}
}

func Test_XXTEA_Native_Roundtrip(t *testing.T) {
	key := ([]byte)("xxteaTEST")
	orig := ([]byte)("01234567890123456789")

	intermediate, _ := XxteaEncrypt(orig, key)

	decoded, _ := XxteaDecrypt(intermediate, key)

	if string(decoded) != string(orig) {
		t.Error("Round trip encode failed, result = " + string(decoded))
	} else {
		t.Log("Round trip xxtea encode passed")
	}
}
