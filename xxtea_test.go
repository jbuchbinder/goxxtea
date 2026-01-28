// GOXXTEA
// https://github.com/jbuchbinder/goxxtea

package xxtea

import (
	"testing"
)

func Test_XXTEA_Roundtrip(t *testing.T) {
	key := "xxteaTEST"
	orig := "01234567890123456789"

	intermediate := Encrypt([]byte(orig), []byte(key))

	decoded := Decrypt(intermediate, []byte(key))

	if string(decoded) != orig {
		t.Error("Round trip encode failed, result = " + string(decoded))
	} else {
		t.Log("Round trip xxtea encode passed")
	}
}
