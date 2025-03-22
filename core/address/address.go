package address

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"go-lucid/util"

	"github.com/mr-tron/base58"
)

type Address string

type PrivateKey struct {
	Key ed25519.PrivateKey
}

type PublicKey struct {
	Key ed25519.PublicKey
}

func (a Address) Validate() bool {
	decoded, err := base58.Decode(a.String())
	if err != nil {
		return false
	}

	checksum := decoded[len(decoded)-4:]
	versionedHash := decoded[:len(decoded)-4]

	hash1 := sha256.Sum256(versionedHash)
	hash2 := sha256.Sum256(hash1[:])

	return bytes.Equal(checksum, hash2[:4])
}

func (a Address) String() string {
	return string(a)
}

func RandomAddress(identifier byte) Address {
	key := ed25519.NewKeyFromSeed(util.GenerateRandomBytes(32))
	pubKey := PublicKey{Key: key.Public().(ed25519.PublicKey)}
	return pubKey.Address(identifier)
}

func (p PrivateKey) Sign(msg []byte) []byte {
	return ed25519.Sign(p.Key, msg)
}

func (p PublicKey) Verify(msg []byte, sig []byte) bool {
	return ed25519.Verify(p.Key, msg, sig)
}

func (p PublicKey) Address(identifier byte) Address {
	hash1 := sha256.Sum256([]byte(p.Key))
	hash2 := sha256.Sum256(hash1[:])

	// Add chain identifier byte
	// (0x42 for mainnet) (0x45 for testnet)
	versionedHash := append([]byte{identifier}, hash2[:]...)

	checksum1 := sha256.Sum256(versionedHash)
	checksum2 := sha256.Sum256(checksum1[:])
	checksumBytes := checksum2[:4]

	finalBytes := append(versionedHash, checksumBytes...)

	return Address(base58.Encode(finalBytes))
}

func (p PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: p.Key.Public().(ed25519.PublicKey)}
}

func (p PrivateKey) Address(identifier byte) Address {
	return p.PublicKey().Address(identifier)
}
