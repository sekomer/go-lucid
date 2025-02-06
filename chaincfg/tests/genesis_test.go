package chaincfg_test

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"reflect"
	"testing"

	"go-lucid/chaincfg"
	"go-lucid/core/block"
)

func TestGenesisBlockHash(t *testing.T) {
	t.Parallel()

	genesis := chaincfg.GenesisBlock

	blockHashBytes, err := genesis.GetHash()
	blockHash := hex.EncodeToString(blockHashBytes)
	if err != nil {
		t.Fatalf("Failed to get hash of genesis block: %v", err)
	}
	if blockHash != hex.EncodeToString(genesis.Hash) {
		t.Fatalf("Hash of genesis block does not match")
	}

}

func TestGenesisBlockSerialization(t *testing.T) {
	t.Parallel()

	genesis := chaincfg.GenesisBlock
	serialized, err := genesis.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize genesis block: %v", err)
	}

	deserialized := block.Block{}
	err = deserialized.Deserialize(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize genesis block: %v", err)
	}

	if !reflect.DeepEqual(genesis, deserialized) {
		t.Fatalf("Genesis block does not match deserialized block")
	}
}

func TestGenesisBlockSize(t *testing.T) {
	t.Parallel()

	genesis := chaincfg.GenesisBlock
	serialized, err := genesis.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize genesis block: %v", err)
	}

	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	gz.Write(serialized)
	gz.Close()

	t.Log("compressed genesis block size:", compressed.Len())
	t.Log("uncompressed genesis block size:", len(serialized))
}
