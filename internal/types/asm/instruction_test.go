package asm

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	for name, test := range map[string]struct {
		input       string
		expected    Instruction
		errExpected bool
	}{
		"nullary": {
			input:       `["EndDialogue"]`,
			expected:    Instruction{Opcode: EndDialogue, Arg: Value{}},
			errExpected: false,
		},
		"unary (number val)": {
			input:       `["PushNumber", ["number", 5]]`,
			expected:    Instruction{Opcode: PushNumber, Arg: Value{Type: NumberType, Val: 5}},
			errExpected: false,
		},
		"unary (bool val)": {
			input:       `["PushBool", ["boolean", true]]`,
			expected:    Instruction{Opcode: PushBool, Arg: Value{Type: BooleanType, Val: true}},
			errExpected: false,
		},
		"unary (string val)": {
			input:       `["PushString", ["string", "5"]]`,
			expected:    Instruction{Opcode: PushString, Arg: Value{Type: StringType, Val: "5"}},
			errExpected: false,
		},
		"unary (symbol val)": {
			input:       `["EnterNode", ["symbol", "abc123"]]`,
			expected:    Instruction{Opcode: EnterNode, Arg: Value{Type: SymbolType, Val: "abc123"}},
			errExpected: false,
		},
		"unary wrong arity": {
			input:       `["PushNumber"]`,
			expected:    Instruction{},
			errExpected: true,
		},
		"invalid opcode": {
			input:       `["Bogus"]`,
			expected:    Instruction{},
			errExpected: true,
		},
		"invalid arg (number val)": {
			input:       `["PushNumber", ["number", "5"]]`,
			expected:    Instruction{},
			errExpected: true,
		},
		"invalid arg (bool val)": {
			input:       `["PushBool", ["boolean", "true"]]`,
			expected:    Instruction{},
			errExpected: true,
		},
		"invalid arg (string val)": {
			input:       `["PushString", ["string", 5]]`,
			expected:    Instruction{},
			errExpected: true,
		},
		"invalid arg (symbol val)": {
			input:       `["EnterNode", ["symbol", 5]]`,
			expected:    Instruction{},
			errExpected: true,
		},
		"invalid type tag": {
			input:       `["EnterNode", ["float", 5]]`,
			expected:    Instruction{},
			errExpected: true,
		},
		"bad json": {
			input:       `["EnterNode", ["float", 5]`,
			expected:    Instruction{},
			errExpected: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := Instruction{}
			err := actual.UnmarshalJSON([]byte(test.input))
			errActual := err != nil
			if test.errExpected != errActual {
				t.Fatalf("unexpected err result: %v", err)
			}
			if err != nil {
				return
			}
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}

	val := Value{}
	err := val.UnmarshalJSON([]byte(`["string" "abc"]`))
	if err == nil {
		t.Fatalf("unexpected err result: %v", err)
	}
}

func TestMarshal(t *testing.T) {
	for name, test := range map[string]struct {
		input       Instruction
		expected    string
		errExpected bool
	}{
		"nullary op": {
			input:       Instruction{EndDialogue, Value{}},
			expected:    `["EndDialogue"]`,
			errExpected: false,
		},
		"unary op": {
			input:       Instruction{PushNumber, Value{NumberType, 5}},
			expected:    `["PushNumber", ["number", 5]]`,
			errExpected: false,
		},
		"bogus opcode": {
			input:       Instruction{Opcode: Opcode("blah")},
			expected:    "",
			errExpected: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual, err := test.input.MarshalJSON()
			errActual := err != nil
			if test.errExpected != errActual {
				t.Fatalf("unexpected err result: %v", err)
			}

			var compactExpectedBuffer, compactActualBuffer bytes.Buffer
			json.Compact(&compactExpectedBuffer, []byte(test.expected))
			json.Compact(&compactActualBuffer, actual)

			e, a := compactExpectedBuffer.String(), compactActualBuffer.String()

			if e != a {
				t.Errorf("expected %v got %v", e, a)
			}
		})
	}
}
