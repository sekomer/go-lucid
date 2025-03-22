package script

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go-lucid/core/address"
)

type Interpreter struct {
	UnlockScript Script // user script
	LockScript   Script // utxo script
	Stack        *Stack
	Err          error
	Msg          []byte // tx being verified
}

func (i *Interpreter) Run() {
	defer func() {
		if r := recover(); r != nil {
			i.Err = errors.New("panic in script interpreter")
		}
	}()

	items, err := i.UnlockScript.Parse()
	if err != nil {
		i.Err = err
		return
	}
	for _, item := range items {
		i.Stack.Push(item)
	}

	if i.LockScript.Data == "" {
		return
	}

	items, err = i.LockScript.Parse()
	if err != nil {
		i.Err = err
		return
	}

	for _, item := range items {
		if i.Err != nil {
			return
		}

		switch item {
		case OP_DUP:
			if i.Stack.Size() < 1 {
				i.Err = errors.New("stack underflow")
				return
			}
			i.Stack.Push(i.Stack.Peek())
		case OP_EQUALVERIFY:
			if i.Stack.Size() < 2 {
				i.Err = errors.New("stack underflow")
				return
			}
			first := i.Stack.Pop()
			second := i.Stack.Pop()
			if first != second {
				i.Err = errors.New("equal verify failed")
				return
			}
			i.Stack.Push(OP_TRUE)
		case OP_HASH256:
			if i.Stack.Size() < 1 {
				i.Err = errors.New("stack underflow")
				return
			}
			top := i.Stack.Pop()
			hash1 := sha256.Sum256([]byte(top))
			hash2 := sha256.Sum256(hash1[:])
			i.Stack.Push(StackItem(hex.EncodeToString(hash2[:])))
		case OP_CHECKSIG:
			if i.Stack.Size() < 2 {
				i.Err = errors.New("stack underflow")
				return
			}
			pubKeyBytes := []byte(i.Stack.Pop())
			sigBytes := []byte(i.Stack.Pop())

			pubKey := address.PublicKey{Key: pubKeyBytes}
			if !pubKey.Verify(i.Msg, sigBytes) {
				i.Stack.Push(OP_FALSE)
			} else {
				i.Stack.Push(OP_TRUE)
			}
		case OP_RETURN:
			i.Err = errors.New("op_return encountered")
			return
		default:
			// Push the value onto the stack
			i.Stack.Push(StackItem(item))
		}
	}

	// * after execution, stack should have at least one item
	if i.Stack.Size() < 1 {
		i.Err = errors.New("empty stack after execution")
		return
	}

	// * check if the top of the stack is truthy
	top := i.Stack.Pop()
	if top != OP_TRUE {
		i.Err = errors.New("script evaluated to false")
		return
	}
}
