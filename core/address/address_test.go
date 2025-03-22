package address_test

import (
	"crypto/ed25519"
	"go-lucid/core/address"
	"strings"
	"testing"
)

var (
	testIdentifier     = byte(0x42)
	testSeed           = strings.Repeat("5", 32)
	testAddress        = "3E3CFmPFKUsif14ANVob1vZ5ZFVTWWtoNSxXQkgjumKpkGvtWF1"
	testInvalidAddress = "1e3cfMP1KUsiN14ANVob1Vz5ZABCDwTonsXxqKGJUMkpkGvtWF2"
)

func TestRandomAddress(t *testing.T) {
	address := address.RandomAddress(testIdentifier)

	t.Log("generated address: ", address)
	t.Log("length: ", len(address))

	if !address.Validate() {
		t.Errorf("validation should pass for random address")
	}
}

func TestDeterministicInvalid(t *testing.T) {
	address := address.Address(testInvalidAddress)

	if address.Validate() {
		t.Errorf("validation should fail for invalid address")
	}
}

func TestDeterministicValid(t *testing.T) {
	key := ed25519.NewKeyFromSeed([]byte(testSeed))
	pubKey := address.PrivateKey{key}
	address := pubKey.Address(testIdentifier)

	t.Log("generated address: ", address)
	t.Log("length: ", len(address))

	if !address.Validate() {
		t.Errorf("validation should pass for valid address")
	}

	if testAddress != address.String() {
		t.Errorf("generated address does not match expected address")
	}
}
