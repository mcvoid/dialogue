package vm

import (
	"fmt"

	"github.com/mcvoid/dialogue/internal/types/asm"
)

type (
	runState int
	choice   struct {
		text asm.Value
		dest asm.Value
	}
)

const (
	runningState runState = iota
	suspendedState
	waitingForInputState
	stoppedState
	errorState
)

var stackNeeded = map[asm.Opcode]int{
	asm.PopValue:      1,
	asm.PushBool:      0,
	asm.PushNull:      0,
	asm.PushNumber:    0,
	asm.PushString:    0,
	asm.GreaterThan:   2,
	asm.Lessthan:      2,
	asm.Concat:        2,
	asm.And:           2,
	asm.Or:            2,
	asm.Not:           1,
	asm.Equal:         2,
	asm.Add:           2,
	asm.Subtract:      2,
	asm.Multiply:      2,
	asm.Divide:        2,
	asm.Increment:     1,
	asm.Decrement:     1,
	asm.LoadVariable:  0,
	asm.StoreVariable: 1,
	asm.ShowLine:      1,
	asm.Jump:          0,
	asm.JumpIfFalse:   1,
	asm.PushChoice:    1,
	asm.ShowChoice:    0,
	asm.EnterNode:     0,
	asm.ExitNode:      0,
	asm.EndDialogue:   0,
	asm.Call:          0,
}

func run(vm *VM) error {
	for vm.runState == runningState {
		if err := singleStep(vm); err != nil {
			vm.runState = errorState
			return err
		}
	}
	return nil
}

func push(vm *VM, val asm.Value) {
	vm.stack = append(vm.stack, val)
}

func pop(vm *VM) asm.Value {
	l := len(vm.stack) - 1
	val := vm.stack[l]
	vm.stack = vm.stack[:l]
	return val
}

func singleStep(vm *VM) error {
	if vm.pc < 0 || vm.pc >= len(vm.code) {
		return fmt.Errorf("%d: jumped to out of bounds location", vm.pc)
	}
	instr := vm.code[vm.pc]
	vm.pc++

	if size := stackNeeded[instr.Opcode]; size > len(vm.stack) {
		return fmt.Errorf("%d: vm stack underflow", vm.pc)
	}

	switch instr.Opcode {
	case asm.Call:
		{
			funcName := instr.Arg
			callback, ok := vm.functions[funcName]
			prototype, protoOk := vm.prototypes[funcName]
			if !ok || !protoOk {
				return fmt.Errorf("%d: callback not found: %v", vm.pc, funcName)
			}
			if len(vm.stack) < len(prototype) {
				return fmt.Errorf("%d: vm stack underflow", vm.pc)
			}
			args := []asm.Value{}
			for range prototype {
				arg := pop(vm)
				args = append([]asm.Value{arg}, args...)
			}
			for i, paramType := range prototype {
				if args[i].Type != paramType {
					return fmt.Errorf("%d: callback %v param %d type error, expected %v got %v", vm.pc, funcName, i, paramType, args[i].Type)
				}
			}

			if executionType := callback.Func(vm, args...); executionType == PauseExecution {
				vm.runState = suspendedState
			}
		}
	case asm.EndDialogue:
		vm.handleEndDialogue(vm)
		vm.runState = stoppedState
	case asm.EnterNode:
		{
			nodeName := instr.Arg.Val.(string)
			if executionType := vm.handleEnterNode(vm, nodeName); executionType == PauseExecution {
				vm.runState = suspendedState
			}
		}
	case asm.ExitNode:
		{
			nodeName := instr.Arg.Val.(string)
			if executionType := vm.handleExitNode(vm, nodeName); executionType == PauseExecution {
				vm.runState = suspendedState
			}
		}
	case asm.ShowLine:
		{
			line := pop(vm)
			lineText, ok := line.Val.(string)
			if line.Type != asm.StringType || !ok {
				return fmt.Errorf("%d: value %v is not of type string", vm.pc, line)
			}
			if executionType := vm.handleShowLine(vm, lineText); executionType == PauseExecution {
				vm.runState = suspendedState
			}
		}
	case asm.LoadVariable:
		{
			symbol := instr.Arg
			val, ok := vm.variables[symbol]
			if !ok {
				val = asm.Null
			}
			push(vm, val)
		}
	case asm.StoreVariable:
		{
			symbol := instr.Arg
			val := pop(vm)
			vm.variables[symbol] = val
		}
	case asm.PushChoice:
		{
			dest := instr.Arg
			str := pop(vm)
			if str.Type != asm.StringType {
				return fmt.Errorf("value %v is not of type String", str)
			}
			if dest.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", dest)
			}
			vm.choices = append(vm.choices, choice{text: str, dest: dest})
		}
	case asm.ShowChoice:
		{
			optionText := []string{}
			for _, choice := range vm.choices {
				optionText = append(optionText, string(choice.text.Val.(string)))
			}
			vm.runState = waitingForInputState
			vm.handleShowChoice(vm, optionText)
		}
	case asm.Jump:
		{
			dest := instr.Arg.Val.(int)
			vm.pc = dest
		}
	case asm.JumpIfFalse:
		{
			dest := instr.Arg.Val.(int)
			val := pop(vm)
			if val.Type != asm.BooleanType {
				return fmt.Errorf("%d: value %v is not of type Boolean", vm.pc, val)
			}
			if val == asm.False {
				vm.pc = dest
			}
		}
	case asm.PushNumber:
		{
			val := instr.Arg
			if val.Type != asm.NumberType {
				return fmt.Errorf("%d: value %v is not of type Number", vm.pc, val)
			}
			push(vm, val)
		}
	case asm.PushString:
		{
			val := instr.Arg
			if val.Type != asm.StringType {
				return fmt.Errorf("%d: value %v is not of type String", vm.pc, val)
			}
			push(vm, val)
		}
	case asm.PushBool:
		{
			val := instr.Arg
			if val.Type != asm.BooleanType {
				return fmt.Errorf("%d: value %v is not of type Boolean", vm.pc, val)
			}
			push(vm, val)
		}
	case asm.PushNull:
		{
			push(vm, asm.Null)
		}
	case asm.PopValue:
		pop(vm)
	case asm.Negative:
		{
			val := pop(vm)
			if val.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val)
			}
			push(vm, asm.Value{Type: asm.NumberType, Val: -(val.Val.(int))})
		}
	case asm.Concat:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1 == asm.Null {
				val1.Val = "null"
			}
			if val2 == asm.Null {
				val2.Val = "null"
			}
			str := fmt.Sprintf("%v%v", val1.Val, val2.Val)
			push(vm, asm.Value{Type: asm.StringType, Val: str})
		}
	case asm.Add:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.NumberType, Val: num1 + num2})
		}
	case asm.Subtract:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.NumberType, Val: num1 - num2})
		}
	case asm.Multiply:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.NumberType, Val: num1 * num2})
		}
	case asm.Divide:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.NumberType, Val: num1 / num2})
		}
	case asm.Modulo:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.NumberType, Val: num1 % num2})
		}
	case asm.Equal:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			push(vm, asm.Value{Type: asm.BooleanType, Val: val1 == val2})
		}
	case asm.NotEqual:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			push(vm, asm.Value{Type: asm.BooleanType, Val: val1 != val2})
		}
	case asm.GreaterThan:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.BooleanType, Val: num1 > num2})
		}
	case asm.Lessthan:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.BooleanType, Val: num1 < num2})
		}
	case asm.GreaterThanOrEqual:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.BooleanType, Val: num1 >= num2})
		}
	case asm.LessthanOrEqual:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val1)
			}
			if val2.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val2)
			}
			num1 := val1.Val.(int)
			num2 := val2.Val.(int)
			push(vm, asm.Value{Type: asm.BooleanType, Val: num1 <= num2})
		}
	case asm.And:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.BooleanType {
				return fmt.Errorf("value %v is not of type Boolean", val1)
			}
			if val2.Type != asm.BooleanType {
				return fmt.Errorf("value %v is not of type Boolean", val2)
			}
			bool1 := val1.Val.(bool)
			bool2 := val2.Val.(bool)
			push(vm, asm.Value{Type: asm.BooleanType, Val: bool1 && bool2})
		}
	case asm.Or:
		{
			val2 := pop(vm)
			val1 := pop(vm)
			if val1.Type != asm.BooleanType {
				return fmt.Errorf("value %v is not of type Boolean", val1)
			}
			if val2.Type != asm.BooleanType {
				return fmt.Errorf("value %v is not of type Boolean", val2)
			}
			bool1 := val1.Val.(bool)
			bool2 := val2.Val.(bool)
			push(vm, asm.Value{Type: asm.BooleanType, Val: bool1 || bool2})
		}
	case asm.Not:
		{
			val := pop(vm)
			if val.Type != asm.BooleanType {
				return fmt.Errorf("value %v is not of type Boolean", val)
			}
			bool1 := val.Val.(bool)
			push(vm, asm.Value{Type: asm.BooleanType, Val: !bool1})
		}
	case asm.Increment:
		{
			val := pop(vm)
			if val.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val)
			}
			num := val.Val.(int)
			push(vm, asm.Value{Type: asm.NumberType, Val: num + 1})
		}
	case asm.Decrement:
		{
			val := pop(vm)
			if val.Type != asm.NumberType {
				return fmt.Errorf("value %v is not of type Number", val)
			}
			num := val.Val.(int)
			push(vm, asm.Value{Type: asm.NumberType, Val: num - 1})
		}
	default:
		return fmt.Errorf("%d: invalid instruction: %v", vm.pc, instr)
	}

	return nil
}
