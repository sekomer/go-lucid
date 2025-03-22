package script

import (
	"errors"
	"strings"
)

type Script struct {
	Data string
}

func (s *Script) Parse() ([]StackItem, error) {
	if len(s.Data) == 0 {
		return []StackItem{}, nil
	}

	parts := strings.Split(s.Data, " ")
	items := make([]StackItem, 0)

	if len(parts) > MAX_SCRIPT_OPS {
		return nil, errors.New("script is too long")
	}

	for _, part := range parts {
		switch part {
		case OP_DUP:
			items = append(items, OP_DUP)
		case OP_EQUALVERIFY:
			items = append(items, OP_EQUALVERIFY)
		case OP_HASH256:
			items = append(items, OP_HASH256)
		case OP_CHECKSIG:
			items = append(items, OP_CHECKSIG)
		case OP_RETURN:
			items = append(items, OP_RETURN)
		default:
			items = append(items, StackItem(part))
		}
	}
	return items, nil
}
