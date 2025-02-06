package block

import (
	"bytes"
	"encoding/gob"
)

func (b *Block) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *Block) Deserialize(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	return dec.Decode(b)
}
