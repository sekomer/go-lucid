package script_test

import (
	"crypto/sha256"
	"encoding/hex"
	"go-lucid/core/script"
	"testing"
)

func TestTrue(t *testing.T) {
	unlockScript := script.Script{Data: ""}
	lockScript := script.Script{Data: "OP_TRUE"}

	interpreter := script.Interpreter{
		UnlockScript: unlockScript,
		LockScript:   lockScript,
		Stack:        &script.Stack{},
	}
	interpreter.Run()

	if interpreter.Err != nil {
		t.Fatalf("error: %v", interpreter.Err)
	}
}

func TestHash256(t *testing.T) {
	/*
		check if given data given in unlock script satisfies the condition given in lock script:
			UnlockScript: <data_to_hash>
			LockScript:   OP_HASH256 <hash_of_data> OP_EQUALVERIFY
	*/

	dataToHash := "lUcIdDrEaM"
	hashOfData := sha256.Sum256([]byte(dataToHash))
	hashOfHash := sha256.Sum256(hashOfData[:])
	hashOfHashString := hex.EncodeToString(hashOfHash[:])

	unlockScript := script.Script{Data: dataToHash}
	lockScript := script.Script{Data: "OP_HASH256 " + hashOfHashString + " OP_EQUALVERIFY"}

	interpreter := script.Interpreter{
		UnlockScript: unlockScript,
		LockScript:   lockScript,
		Stack:        &script.Stack{},
	}
	interpreter.Run()

	if interpreter.Err != nil {
		t.Fatalf("error: %v", interpreter.Err)
	}
}

func TestEqualVerify(t *testing.T) {
	unlockScript := script.Script{Data: "neo"}
	lockScript := script.Script{Data: "neo OP_EQUALVERIFY"}

	interpreter := script.Interpreter{
		UnlockScript: unlockScript,
		LockScript:   lockScript,
		Stack:        &script.Stack{},
	}
	interpreter.Run()

	if interpreter.Err != nil {
		t.Fatalf("error: %v", interpreter.Err)
	}
}
