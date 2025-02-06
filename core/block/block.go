package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// * BlockHeader

func (h *BlockHeader) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(h)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *BlockHeader) Deserialize(b []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(b))
	return dec.Decode(h)
}

// * Block

func (b *Block) GetHash() ([]byte, error) {
	headerBytes, err := b.BlockHeader.Serialize()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(headerBytes)
	return hash[:], nil
}

func (b *Block) GetMerkleRoot() ([]byte, error) {
	if len(b.Txs) == 0 {
		return nil, nil
	}

	hashes := make([][]byte, len(b.Txs))
	for i, tx := range b.Txs {
		serialized, err := tx.Serialize()
		if err != nil {
			return nil, err
		}
		hash := sha256.Sum256(serialized)
		hashes[i] = hash[:]
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}

		newLevel := make([][]byte, 0, len(hashes)/2)
		for i := 0; i < len(hashes); i += 2 {
			hash := sha256.Sum256(append(hashes[i], hashes[i+1]...))
			newLevel = append(newLevel, hash[:])
		}
		hashes = newLevel
	}

	var root []byte
	copy(root, hashes[0])
	return root, nil
}

// * BlockModel
