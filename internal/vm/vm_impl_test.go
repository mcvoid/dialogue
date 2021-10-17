package vm

import (
	"fmt"
	"testing"

	"github.com/mcvoid/dialogue/internal/program"
	"github.com/mcvoid/dialogue/internal/types/asm"
)

var emptyProgram = program.Program{
	Start: 0,
	Code: []asm.Instruction{
		{Opcode: asm.EndDialogue, Arg: asm.Value{}},
	},
	Funcs: map[string][]asm.Type{},
}

func compareValue(a, b []asm.Value) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func compareStrings(a, b []string) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func TestOpcodes(t *testing.T) {
	tests := []struct {
		code     []asm.Instruction
		expected []asm.Value
	}{
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 3}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "3"}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.StringType, Val: "3"}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{asm.Null},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PopValue, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 5}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.JumpIfFalse, Arg: asm.Value{Type: asm.NumberType, Val: 6}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.JumpIfFalse, Arg: asm.Value{Type: asm.NumberType, Val: 6}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{asm.Null, asm.Null, asm.Null, asm.Null},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 30}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
				{Opcode: asm.Add, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 30 + 4}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 30}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
				{Opcode: asm.Subtract, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 30 - 4}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 30}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
				{Opcode: asm.Multiply, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 30 * 4}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 30}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
				{Opcode: asm.Divide, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 30 / 4}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 30}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
				{Opcode: asm.Modulo, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 30 % 4}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.And, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.And, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.And, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.And, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.Or, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.Or, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.Or, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.Or, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.Not, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.Not, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.Increment, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 18}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.Decrement, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: 16}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 34}},
				{Opcode: asm.GreaterThan, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 34}},
				{Opcode: asm.Lessthan, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 34}},
				{Opcode: asm.GreaterThanOrEqual, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 34}},
				{Opcode: asm.LessthanOrEqual, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.GreaterThanOrEqual, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.LessthanOrEqual, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "17"}},
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "34"}},
				{Opcode: asm.Concat, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.StringType, Val: "1734"}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "34"}},
				{Opcode: asm.Concat, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.StringType, Val: "1734"}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.Concat, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.StringType, Val: "17null"}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNull, Arg: asm.Value{}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.Concat, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.StringType, Val: "null17"}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
				{Opcode: asm.Concat, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.StringType, Val: "17true"}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
				{Opcode: asm.Concat, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.StringType, Val: "17false"}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.Equal, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 16}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.Equal, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.NotEqual, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: false}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 16}},
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
				{Opcode: asm.NotEqual, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.BooleanType, Val: true}},
		},
		{
			[]asm.Instruction{
				{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 16}},
				{Opcode: asm.Negative, Arg: asm.Value{}},
				{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			},
			[]asm.Value{{Type: asm.NumberType, Val: -16}},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			vm, _ := New(program.Program{
				Start: 0,
				Code:  test.code,
			})
			err := vm.Run()
			if err != nil {
				t.Error("No error expected on correct program")
			}
			if !compareValue(vm.stack, test.expected) {
				t.Errorf("%d: expected %v got %v", i, test.expected, vm.stack)
			}
		})
	}
}

func TestOpcodeErrors(t *testing.T) {
	tests := [][]asm.Instruction{
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.JumpIfFalse, Arg: asm.Value{Type: asm.NumberType, Val: 2}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.Add, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Add, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.Subtract, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Subtract, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.Multiply, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Multiply, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.Divide, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Divide, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.Modulo, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Modulo, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.GreaterThan, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.GreaterThan, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.Lessthan, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Lessthan, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.GreaterThanOrEqual, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.GreaterThanOrEqual, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.LessthanOrEqual, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.LessthanOrEqual, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
			{Opcode: asm.And, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.And, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
			{Opcode: asm.Or, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Or, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Not, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Increment, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "30"}},
			{Opcode: asm.Decrement, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.Opcode("invalid"), Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
			{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 4}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "3"}},
			{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.StringType, Val: "4"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "3"}},
			{Opcode: asm.Negative, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			vm, _ := New(program.Program{
				Start: 0,
				Code:  test,
			})
			err := vm.Run()
			if err == nil {
				t.Error("error expected on incorrect program")
			}
		})
	}
}

func TestVMCallbacks(t *testing.T) {
	prog := program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "abc"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
			{Opcode: asm.ShowLine, Arg: asm.Value{}},
			{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "abc"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	vm, _ := New(prog)
	err := vm.Run()
	if err != nil {
		t.Error("No error expected on correct program")
	}

	prog = program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 13}},
			{Opcode: asm.ShowLine, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	vm, _ = New(prog)
	err = vm.Run()
	if err == nil {
		t.Error("error expected on incorrect program")
	}

	prog = program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "abc"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
			{Opcode: asm.ShowLine, Arg: asm.Value{}},
			{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "abc"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	enterNode, showLine, exitNode := false, false, false
	vm, _ = New(
		prog,
		HandleEnterNode(func(v *VM, s string) ExecutionType {
			enterNode = true
			if s != "abc" {
				t.Errorf("Expected %v got %v", "abc", s)
			}
			return PauseExecution
		}),
		HandleExitNode(func(v *VM, s string) ExecutionType {
			exitNode = true
			if s != "abc" {
				t.Errorf("Expected %v got %v", "abc", s)
			}
			return PauseExecution
		}),
		HandleShowLine(func(v *VM, s string) ExecutionType {
			showLine = true
			if s != "abc" {
				t.Errorf("Expected %v got %v", "abc", s)
			}
			return PauseExecution
		}),
	)
	err = vm.Run()
	if err != nil {
		t.Error("No error expected on correct program")
	}
	if !enterNode {
		t.Error("Expected EnterNode handler to be called")
	}
	if vm.runState != suspendedState {
		t.Errorf("Expected runstate %v got %v", suspendedState, vm.runState)
	}
	err = vm.Resume()
	if err != nil {
		t.Error("No error expected on correct program")
	}
	if !showLine {
		t.Error("Expected ShowLine handler to be called")
	}
	if vm.runState != suspendedState {
		t.Errorf("Expected runstate %v got %v", suspendedState, vm.runState)
	}
	err = vm.Resume()
	if err != nil {
		t.Error("No error expected on correct program")
	}
	if !exitNode {
		t.Error("Expected ExitNode handler to be called")
	}
	if vm.runState != suspendedState {
		t.Errorf("Expected runstate %v got %v", suspendedState, vm.runState)
	}
	err = vm.Resume()
	if err != nil {
		t.Error("No error expected on correct program")
	}
	if vm.runState != stoppedState {
		t.Errorf("Expected runstate %v got %v", stoppedState, vm.runState)
	}
}

func TestVMPushChoice(t *testing.T) {
	prog := program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Choice 1"}},
			{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 7}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Choice 2"}},
			{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 8}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Choice 3"}},
			{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 9}},
			{Opcode: asm.ShowChoice, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	choiceHandled := false
	choiceHandler := func(v *VM, s []string) {
		choiceHandled = true
		expected := []string{"Choice 1", "Choice 2", "Choice 3"}
		if !compareStrings(s, expected) {
			t.Errorf("Expected %v, got %v", expected, s)
		}
	}

	vm, _ := New(prog, HandleShowChoice(choiceHandler))
	vm.Run()
	if !choiceHandled {
		t.Error("ShowChoice handler not called")
	}

	vm.ChooseAndResume(0)
	if vm.pc != 8 {
		t.Errorf("Expected VM to halt on pc=%v, got %v", 8, vm.pc)
	}

	vm, _ = New(prog)
	vm.Run()
	vm.ChooseAndResume(1)
	if vm.pc != 9 {
		t.Errorf("Expected VM to halt on pc=%v, got %v", 9, vm.pc)
	}

	vm.Run()
	vm.ChooseAndResume(1)
	if vm.pc != 9 {
		t.Errorf("Expected VM to halt on pc=%v, got %v", 10, vm.pc)
	}
}

func TestVmCall(t *testing.T) {
	prog := program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: true}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
			{Opcode: asm.PushNull, Arg: asm.Value{}},
			{Opcode: asm.Call, Arg: asm.Value{Type: asm.SymbolType, Val: "func"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		Funcs: map[string][]asm.Type{},
	}

	tests := []struct {
		f                Function
		prototype        []asm.Type
		callbackName     string
		expectedRunState runState
		expectedErr      bool
	}{
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					if len(args) != 0 {
						t.Errorf("Wrong number of params, expected %v got %v", 0, len(args))
					}
					return ContinueExecution
				},
			},
			[]asm.Type{},
			"func",
			stoppedState,
			false,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					if len(args) != 0 {
						t.Errorf("Wrong number of params, expected %v got %v", 0, len(args))
					}
					return ContinueExecution
				},
			},
			[]asm.Type{},
			"badfunc",
			errorState,
			true,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					if len(args) != 0 {
						t.Errorf("Wrong number of params, expected %v got %v", 0, len(args))
					}
					return PauseExecution
				},
			},
			[]asm.Type{},
			"func",
			suspendedState,
			false,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					t.Errorf("Did not expect callback invoked")
					return ContinueExecution
				},
			},
			[]asm.Type{asm.NumberType, asm.NumberType, asm.NumberType, asm.NumberType, asm.NumberType},
			"func",
			errorState,
			true,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					t.Errorf("Did not expect callback invoked")
					return ContinueExecution
				},
			},
			[]asm.Type{asm.NumberType, asm.NumberType, asm.NumberType, asm.NumberType},
			"func",
			errorState,
			true,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					t.Errorf("Did not expect callback invoked")
					return ContinueExecution
				},
			},
			[]asm.Type{asm.NumberType, asm.NumberType, asm.NumberType, asm.NumberType},
			"func1",
			errorState,
			true,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					expectedArgs := []asm.Value{asm.Null}
					if !compareValue(expectedArgs, args) {
						t.Errorf("Wrong params, expected %v, got %v", expectedArgs, args)
					}
					return ContinueExecution
				},
			},
			[]asm.Type{asm.NullType},
			"func",
			stoppedState,
			false,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					expectedArgs := []asm.Value{{Type: asm.NumberType, Val: 17}, asm.Null}
					if !compareValue(expectedArgs, args) {
						t.Errorf("Wrong params, expected %v, got %v", expectedArgs, args)
					}
					return ContinueExecution
				},
			},
			[]asm.Type{asm.NumberType, asm.NullType},
			"func",
			stoppedState,
			false,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					expectedArgs := []asm.Value{{Type: asm.BooleanType, Val: true}, {Type: asm.NumberType, Val: 17}, asm.Null}
					if !compareValue(expectedArgs, args) {
						t.Errorf("Wrong params, expected %v, got %v", expectedArgs, args)
					}
					return ContinueExecution
				},
			},
			[]asm.Type{asm.BooleanType, asm.NumberType, asm.NullType},
			"func",
			stoppedState,
			false,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					expectedArgs := []asm.Value{{Type: asm.StringType, Val: "abc"}, {Type: asm.BooleanType, Val: true}, {Type: asm.NumberType, Val: 17}, asm.Null}
					if !compareValue(expectedArgs, args) {
						t.Errorf("Wrong params, expected %v, got %v", expectedArgs, args)
					}
					return ContinueExecution
				},
			},
			[]asm.Type{asm.StringType, asm.BooleanType, asm.NumberType, asm.NullType},
			"func",
			stoppedState,
			false,
		},
		{
			Function{
				func(vm *VM, args ...asm.Value) ExecutionType {
					expectedArgs := []asm.Value{{Type: asm.StringType, Val: "abc"}, {Type: asm.BooleanType, Val: true}, {Type: asm.NumberType, Val: 17}, asm.Null}
					if !compareValue(expectedArgs, args) {
						t.Errorf("Wrong params, expected %v, got %v", expectedArgs, args)
					}
					return PauseExecution
				},
			},
			[]asm.Type{asm.StringType, asm.BooleanType, asm.NumberType, asm.NullType},
			"func",
			suspendedState,
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			testProg := prog
			testProg.Funcs = map[string][]asm.Type{}
			testProg.Funcs[test.callbackName] = test.prototype
			vm, _ := New(testProg, RegisterCallback(test.f))
			err := vm.Run()
			if (err != nil) != test.expectedErr {
				t.Errorf("%d: Unexpected error state", i)
			}
			if vm.runState != test.expectedRunState {
				t.Errorf("%d: vm runstate expected %v got %v", i, test.expectedRunState, vm.runState)
			}
		})
	}
}
