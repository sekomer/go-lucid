package transaction

import "crypto/sha256"

// * RawTransaction

func (tx *RawTransaction) AddTxIn(in TxIn) {
	tx.TxIns = append(tx.TxIns, in)
}

func (tx *RawTransaction) AddTxOut(out TxOut) {
	tx.TxOuts = append(tx.TxOuts, out)
}

func (tx *TxInModel) ToTxIn() TxIn {
	return TxIn{
		TransactionID:  tx.TransactionID,
		Coinbase:       tx.Coinbase,
		PreviousOutput: tx.PreviousOutput,
		ScriptSig:      tx.ScriptSig,
		Sequence:       tx.Sequence,
	}
}

func (tx *TxOutModel) ToTxOut() TxOut {
	return TxOut{
		TransactionID: tx.TransactionID,
		Value:         tx.Value,
		PkScript:      tx.PkScript,
	}
}

func (tx *RawTransaction) GetHash() ([]byte, error) {
	ser, err := tx.SerializeWithoutHash()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(ser)
	return hash[:], nil
}

// * RawTransactionModel

func (tx *RawTransactionModel) AddTxIn(in TxInModel) {
	tx.TxIns = append(tx.TxIns, in)
}

func (tx *RawTransactionModel) AddTxOut(out TxOutModel) {
	tx.TxOuts = append(tx.TxOuts, out)
}

func (tx *RawTransactionModel) ToRawTransaction() RawTransaction {
	txIns := make([]TxIn, len(tx.TxIns))
	for i, txIn := range tx.TxIns {
		txIns[i] = txIn.ToTxIn()
	}
	txOuts := make([]TxOut, len(tx.TxOuts))
	for i, txOut := range tx.TxOuts {
		txOuts[i] = txOut.ToTxOut()
	}

	return RawTransaction{
		Hash:       tx.Hash,
		Version:    tx.Version,
		TxInCount:  tx.TxInCount,
		TxOutCount: tx.TxOutCount,
		TxIns:      txIns,
		TxOuts:     txOuts,
		BlockID:    tx.BlockID,
		LockTime:   tx.LockTime,
	}
}
