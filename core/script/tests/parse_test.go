package script_test

import (
	"go-lucid/core/script"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	scriptText := "OP_DUP OP_EQUALVERIFY"
	s := script.Script{Data: scriptText}

	items, err := s.Parse()
	if err != nil {
		t.Fatalf("failed to parse script: %v", err)
	}

	t.Logf("parsed script: %v", items)
}

func TestInterpreterFail(t *testing.T) {
	scriptText := "OP_DUP OP_EQUALVERIFY"
	s := script.Script{Data: scriptText}

	items, err := s.Parse()
	if err != nil {
		t.Fatalf("failed to parse script: %v", err)
	}

	stack := script.Stack{Items: make([]script.StackItem, 0)}
	interpreter := script.Interpreter{
		LockScript:   script.Script{Data: scriptText},
		UnlockScript: script.Script{Data: ""},
		Stack:        &stack,
		Msg:          []byte("test"),
	}

	interpreter.Run()
	if interpreter.Err == nil {
		t.Fatalf("expected interpreter to fail")
	}

	t.Logf("parsed script: %v", items)
	t.Logf("interpreter failed as expected: %v", interpreter.Err)
}

func TestParseMaxScriptOps(t *testing.T) {
	scriptText := strings.Repeat("OP_DUP ", script.MAX_SCRIPT_OPS+1)
	s := script.Script{Data: scriptText}

	_, err := s.Parse()
	if err == nil {
		t.Fatalf("expected script to be too long")
	}

	t.Logf("failed as expected: %v", err)
}

func TestInterpreterMaxScriptOps(t *testing.T) {
	scriptText := strings.Repeat("OP_DUP ", script.MAX_SCRIPT_OPS+1)

	interpreter := script.Interpreter{
		LockScript:   script.Script{Data: scriptText},
		UnlockScript: script.Script{Data: ""},
		Stack:        &script.Stack{},
		Msg:          []byte("test"),
	}
	interpreter.Run()
	if interpreter.Err == nil {
		t.Fatalf("expected interpreter to fail")
	}

	t.Logf("failed as expected: %v", interpreter.Err)
}
