package consensus

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"go-lucid/core/block"
)

const difficulty = 4 // Number of leading zeros required in the hash

func calculateHash(block block.Block) ([]byte, error) {
	record, err := block.Serialize()
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hashed, nil
}

func MineBlock(block *block.Block) {
	start := time.Now()
	for {
		block.BlockHeader.Nonce++
		blockHash, err := calculateHash(*block)
		if err != nil {
			log.Println(err)
			continue
		}
		if strings.HasPrefix(hex.EncodeToString(blockHash), strings.Repeat("0", difficulty)) {
			block.Hash = blockHash
			log.Printf("Block mined: %s\n", blockHash)
			log.Printf("Time taken: %s\n", time.Since(start))
			break
		}
	}
}

func IsValidBlock(newBlock, oldBlock block.Block) bool {
	if bytes.Compare(oldBlock.Hash, newBlock.PrevBlock) != 0 {
		return false
	}

	blockHash, err := calculateHash(newBlock)
	if err != nil {
		log.Println(err)
		return false
	}
	if bytes.Compare(blockHash, newBlock.Hash) != 0 {
		return false
	}
	if !strings.HasPrefix(hex.EncodeToString(blockHash), strings.Repeat("0", difficulty)) {
		return false
	}

	return true
}
