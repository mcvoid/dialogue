package program

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/mcvoid/dialogue/internal/types/asm"
)

func compareProgram(a, b Program) bool {
	if a.Start != b.Start {
		return false
	}
	if len(a.Code) != len(b.Code) {
		return false
	}
	for i := range a.Code {
		if a.Code[i] != b.Code[i] {
			return false
		}
	}
	if len(a.Code) != len(b.Code) {
		return false
	}
	for i := range a.Code {
		if a.Code[i] != b.Code[i] {
			return false
		}
	}
	if len(a.Funcs) != len(b.Funcs) {
		return false
	}
	for i := range a.Funcs {
		if len(a.Funcs[i]) != len(b.Funcs[i]) {
			return false
		}
		for j := range a.Funcs[i] {
			if a.Funcs[i][j] != b.Funcs[i][j] {
				return false
			}
		}
	}
	for i := range b.Funcs {
		if len(a.Funcs[i]) != len(b.Funcs[i]) {
			return false
		}
		for j := range a.Funcs[i] {
			if a.Funcs[i][j] != b.Funcs[i][j] {
				return false
			}
		}
	}
	return true
}

type mockErrReader struct{}

func (m mockErrReader) Read(p []byte) (n int, err error) { return 0, fmt.Errorf("error") }

func TestReadProgram(t *testing.T) {
	tests := map[string]struct {
		str           string
		expected      Program
		errorExpected bool
	}{
		"valid": {
			`{
				"start": 0,
				"code": [
					["PushNumber", ["number", 5]],
					["PushBool", ["boolean", true]],
					["PushString", ["string", "abc"]],
					["PushNull"],
					["EnterNode", ["symbol", "abc"]],
					["EndDialogue"]
				],
				"funcs": {
					"func1": ["string", "number", "boolean", "null"]
				}
			}`,
			Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 5}},
					{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
					{Opcode: asm.PushNull, Arg: asm.Value{}},
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "abc"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
				Funcs: map[string][]asm.Type{
					"func1": {asm.StringType, asm.NumberType, asm.BooleanType, asm.NullType},
				},
			},
			false,
		},
		"invalid JSON - missing comma": {
			`
		{
			"start": 1
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"invalid JSON - instructions are string": {
			`
		{
			"start": 2,
			"code": "[]"
		}
	`,
			Program{},
			true,
		},
		"one instruction is a string": {
			`
		{
			"start": 3,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				"EndDialogue"
			]
		}
	`,
			Program{},
			true,
		},
		"Empty instruction": {
			`
		{
			"start": 5,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				[]
			]
		}
	`,
			Program{},
			true,
		},
		"number for nullary opcode": {
			`
		{
			"start": 6,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				[5]
			]
		}
	`,
			Program{},
			true,
		},
		"number for unary opcode": {
			`
		{
			"start": 7,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				[5, ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"number for value literal": {
			`
		{
			"start": 8,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", 5],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"number for type tag": {
			`
		{
			"start": 10,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", [5, "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"invalid type tag": {
			`
		{
			"start": 11,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", ["float", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"number for string": {
			`
		{
			"start": 12,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", 5]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"string for number": {
			`
		{
			"start": 13,
			"code": [
				["PushNumber", ["number", "5"]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"string for boolean": {
			`
		{
			"start": 14,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", "true"]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"number for symbol": {
			`
		{
			"start": 15,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", 5]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"invalid opcode": {
			`
		{
			"start": 17,
			"code": [
				["PushNumber", ["number", 5]],
				["PushBool", ["boolean", true]],
				["Thing", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"bool in number value": {
			`
		{
			"start": 19,
			"code": [
				["PushNumber", ["number", true]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
		"null in number value": {
			`
		{
			"start": 20,
			"code": [
				["PushNumber", ["number", null]],
				["PushBool", ["boolean", true]],
				["PushString", ["string", "abc"]],
				["PushNull"],
				["EnterNode", ["symbol", "abc"]],
				["EndDialogue"]
			]
		}
	`,
			Program{},
			true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := Program{}
			_, err := (&actual).ReadFrom(strings.NewReader(test.str))
			errorActual := err != nil
			if errorActual != test.errorExpected {
				t.Fatalf("Expected %v got %v", test.errorExpected, errorActual)
			}
			if !errorActual && !compareProgram(actual, test.expected) {
				t.Errorf("Expected\n%v\ngot\n%v\n", test.expected, actual)
			}
		})
	}
}

func TestBadReader(t *testing.T) {
	actual := Program{}
	_, err := (&actual).ReadFrom(mockErrReader{})
	if err == nil {
		t.Errorf("expected error on bad reader")
	}
}

func TestWriteProgram(t *testing.T) {
	for name, test := range map[string]struct {
		input         Program
		expected      string
		errorExpected bool
	}{
		"trivial": {
			Program{},
			`{"start": 0, "code": null }`,
			false,
		},
		"a little less trivial": {
			Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			`{"start": 0, "code": [
				["PushBool", ["boolean", true]],
				["EndDialogue"]
			] }`,
			false,
		},
		"bad opcode": {
			Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
					{Opcode: "DoSomething", Arg: asm.Value{}},
				},
			},
			``,
			true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			var actualUncompacted, actual, expected bytes.Buffer
			_, err := test.input.WriteTo(&actualUncompacted)

			if (err != nil) != test.errorExpected {
				t.Fatalf("expected error %v got %v", test.errorExpected, err != nil)
			}

			json.Compact(&actual, actualUncompacted.Bytes())
			json.Compact(&expected, []byte(test.expected))

			a, e := actual.String(), expected.String()
			if a != e {
				t.Errorf("expected\n%v\ngot\n%v\n", e, a)
			}
		})

	}
}
