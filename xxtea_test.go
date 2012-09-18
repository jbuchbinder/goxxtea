// GOXXTEA
// https://github.com/jbuchbinder/goxxtea

package xxtea

import (
	"testing"
)

func Test_XXTEA_Roundtrip(t *testing.T) {
	key := "xxteaTEST"
	orig := "The quick brown fox jumped over the lazy dog."

	intermediate, _ := Encrypt(orig, key)

	decoded, _ := Decrypt(intermediate, key)

	if decoded != orig {
		t.Error("Round trip encode failed, intermediate = " + intermediate)
	} else {
		t.Log("Round trip xxtea encode passed")
	}
}
