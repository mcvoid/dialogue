package vm

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/program"
	"github.com/mcvoid/dialogue/internal/types/asm"
)

func TestNewVM(t *testing.T) {
	_, err := New(emptyProgram)
	if err != nil {
		t.Error("No error expected in valid trivial program")
	}
}

func TestHandlers(t *testing.T) {
	p := emptyProgram
	_, err := New(
		p,
		HandleEnterNode(nil),
	)
	if err == nil {
		t.Error("Expected error when supplied invalid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleEnterNode(func(v *VM, s string) ExecutionType { return ContinueExecution }),
	)
	if err != nil {
		t.Error("No error expected when supplied valid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleExitNode(nil),
	)
	if err == nil {
		t.Error("Expected error when supplied invalid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleExitNode(func(v *VM, s string) ExecutionType { return ContinueExecution }),
	)
	if err != nil {
		t.Error("No error expected when supplied valid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleShowLine(nil),
	)
	if err == nil {
		t.Error("Expected error when supplied invalid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleShowLine(func(v *VM, s string) ExecutionType { return ContinueExecution }),
	)
	if err != nil {
		t.Error("No error expected when supplied valid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleEndDialogue(nil),
	)
	if err == nil {
		t.Error("Expected error when supplied invalid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleEndDialogue(func(v *VM) {}),
	)
	if err != nil {
		t.Error("No error expected when supplied valid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleShowChoice(nil),
	)
	if err == nil {
		t.Error("Expected error when supplied invalid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		HandleShowChoice(func(v *VM, s []string) {}),
	)
	if err != nil {
		t.Error("No error expected when supplied valid handler")
	}

	p = emptyProgram
	_, err = New(
		p,
		RegisterCallback(Function{nil}),
	)
	if err == nil {
		t.Error("Expected error when supplied invalid handler")
	}

	p = emptyProgram
	p.Funcs["callback"] = []asm.Type{}
	_, err = New(
		p,
		RegisterCallback(Function{func(vm *VM, args ...asm.Value) ExecutionType { return ContinueExecution }}),
	)
	if err != nil {
		t.Error("No error expected when supplied valid handler")
	}

	p = emptyProgram
	p.Funcs["callback"] = []asm.Type{}
	_, err = New(
		p,
		HandleEnterNode(func(v *VM, s string) ExecutionType { return ContinueExecution }),
		HandleExitNode(func(v *VM, s string) ExecutionType { return ContinueExecution }),
		HandleShowLine(func(v *VM, s string) ExecutionType { return ContinueExecution }),
		HandleEndDialogue(nil),
		HandleShowChoice(func(v *VM, s []string) {}),
		RegisterCallback(Function{func(vm *VM, args ...asm.Value) ExecutionType { return ContinueExecution }}),
	)
	if err == nil {
		t.Error("Expected error when supplied invalid handler")
	}

	p = emptyProgram
	p.Funcs["callback"] = []asm.Type{}
	_, err = New(
		p,
		HandleEnterNode(func(v *VM, s string) ExecutionType { return ContinueExecution }),
		HandleExitNode(func(v *VM, s string) ExecutionType { return ContinueExecution }),
		HandleShowLine(func(v *VM, s string) ExecutionType { return ContinueExecution }),
		HandleEndDialogue(func(v *VM) {}),
		HandleShowChoice(func(v *VM, s []string) {}),
		RegisterCallback(Function{func(vm *VM, args ...asm.Value) ExecutionType { return ContinueExecution }}),
	)
	if err != nil {
		t.Error("No error expected when supplied valid handler")
	}
}

func TestVmRun(t *testing.T) {
	vm, _ := New(emptyProgram)
	vm.runState = runningState

	err := vm.Run()
	if err == nil {
		t.Error("Expected error when running VM in invalid state")
	}

	vm, _ = New(emptyProgram)
	vm.runState = suspendedState

	err = vm.Run()
	if err == nil {
		t.Error("Expected error when running VM in invalid state")
	}

	vm, _ = New(emptyProgram)
	vm.runState = waitingForInputState

	err = vm.Run()
	if err == nil {
		t.Error("Expected error when running VM in invalid state")
	}

	vm, _ = New(emptyProgram)
	vm.runState = errorState

	err = vm.Run()
	if err == nil {
		t.Error("Expected error when running VM in invalid state")
	}

	vm, _ = New(emptyProgram)

	err = vm.Run()
	if err != nil {
		t.Error("No error expected when running VM in a valid state")
	}

	vm, _ = New(emptyProgram)
	vm.runState = stoppedState

	err = vm.Run()
	if err != nil {
		t.Error("No error expected when running VM in a valid state")
	}

	vm.code = []asm.Instruction{}
	err = vm.Run()
	if err == nil {
		t.Error("Expected error when running VM in invalid state")
	}

	vm.code = []asm.Instruction{
		{Opcode: asm.Add, Arg: asm.Value{}},
		{Opcode: asm.EndDialogue, Arg: asm.Value{}},
	}
	vm.runState = stoppedState
	vm.stack = []asm.Value{}

	err = vm.Run()
	if err == nil {
		t.Error("Expected error when running VM in invalid state")
	}
}

func TestVmResume(t *testing.T) {
	vm, _ := New(emptyProgram)

	err := vm.Resume()
	if err == nil {
		t.Error("Expected error when resuming non-suspended VM")
	}

	vm, _ = New(emptyProgram)
	vm.runState = runningState

	err = vm.Resume()
	if err == nil {
		t.Error("Expected error when resuming non-suspended VM")
	}

	vm, _ = New(emptyProgram)
	vm.runState = waitingForInputState

	err = vm.Resume()
	if err == nil {
		t.Error("Expected error when resuming non-suspended VM")
	}

	vm, _ = New(emptyProgram)
	vm.runState = errorState

	err = vm.Resume()
	if err == nil {
		t.Error("Expected error when resuming non-suspended VM")
	}

	vm, _ = New(emptyProgram)
	vm.runState = stoppedState

	err = vm.Resume()
	if err == nil {
		t.Error("Expected error when resuming non-suspended VM")
	}

	vm, _ = New(emptyProgram)
	vm.runState = suspendedState

	err = vm.Resume()
	if err != nil {
		t.Error("No error expected when resuming suspended VM")
	}
}

func TestVmReset(t *testing.T) {
	vm, _ := New(emptyProgram)
	vm.Reset()
	if vm.runState != stoppedState {
		t.Error("Expected reset to put vm in stopped state")
	}

	vm.runState = runningState
	vm.Reset()
	if vm.runState != stoppedState {
		t.Error("Expected reset to put vm in stopped state")
	}

	vm.runState = suspendedState
	vm.Reset()
	if vm.runState != stoppedState {
		t.Error("Expected reset to put vm in stopped state")
	}

	vm.runState = waitingForInputState
	vm.Reset()
	if vm.runState != stoppedState {
		t.Error("Expected reset to put vm in stopped state")
	}

	vm.runState = errorState
	vm.Reset()
	if vm.runState != stoppedState {
		t.Error("Expected reset to put vm in stopped state")
	}

	vm.SetVariableBoolean("abc", true)
	vm.SetVariableNull("def")
	vm.SetVariableNumber("bcd", 15)
	vm.SetVariableString("cde", "cde")
	vm.Reset()
	if len(vm.variables) != 0 {
		t.Error("Expected Reset to clear variables")
	}
}

func TestVmChoose(t *testing.T) {
	prog := program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	vm, _ := New(prog)
	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = waitingForInputState
	err := vm.ChooseAndResume(0)
	if err != nil {
		t.Error("No error expected on valid choice")
	}
	if len(vm.choices) != 0 {
		t.Error("Choosing option did not clear pushed choices")
	}

	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = waitingForInputState
	err = vm.ChooseAndResume(4)
	if err != nil {
		t.Error("No error expected on valid choice")
	}

	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = waitingForInputState
	err = vm.ChooseAndResume(-1)
	if err == nil {
		t.Error("Expected error on invalid choice")
	}

	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = waitingForInputState
	err = vm.ChooseAndResume(5)
	if err == nil {
		t.Error("Expected error on invalid choice")
	}

	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = suspendedState
	err = vm.ChooseAndResume(5)
	if err == nil {
		t.Error("Expected error on invalid choice")
	}

	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = runningState
	err = vm.ChooseAndResume(5)
	if err == nil {
		t.Error("Expected error on invalid choice")
	}

	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = stoppedState
	err = vm.ChooseAndResume(5)
	if err == nil {
		t.Error("Expected error on invalid choice")
	}

	vm.choices = []choice{
		{asm.Value{Type: asm.StringType, Val: "choice0"}, asm.Value{Type: asm.NumberType, Val: 0}},
		{asm.Value{Type: asm.StringType, Val: "choice1"}, asm.Value{Type: asm.NumberType, Val: 1}},
		{asm.Value{Type: asm.StringType, Val: "choice2"}, asm.Value{Type: asm.NumberType, Val: 2}},
		{asm.Value{Type: asm.StringType, Val: "choice3"}, asm.Value{Type: asm.NumberType, Val: 3}},
		{asm.Value{Type: asm.StringType, Val: "choice4"}, asm.Value{Type: asm.NumberType, Val: 4}},
	}
	vm.runState = errorState
	err = vm.ChooseAndResume(5)
	if err == nil {
		t.Error("Expected error on invalid choice")
	}
}

func TestVmGetVariable(t *testing.T) {
	vm, _ := New(emptyProgram)
	val, ok := vm.GetVariable("abc")
	if val != asm.Null {
		t.Error("Expected nonexitent variable to be null")
	}
	if ok {
		t.Error("Expected nonexistent variable to return false")
	}

	abc := asm.Value{Type: asm.NumberType, Val: 1}
	vm.variables[asm.Value{Type: asm.SymbolType, Val: "abc"}] = abc
	val, ok = vm.GetVariable("abc")
	if val != abc {
		t.Errorf("Expected %v got %v", abc, val)
	}
	if !ok {
		t.Errorf("Expected var to exist")
	}
	val, ok = vm.GetVariable("def")
	if val != asm.Null {
		t.Error("Expected nonexitent variable to be null")
	}
	if ok {
		t.Error("Expected nonexistent variable to return false")
	}
}

func TestVmSetVariableNumber(t *testing.T) {
	vm := VM{}
	vm.variables = map[asm.Value]asm.Value{}

	vm.SetVariableNumber("abc", 5)
	val, ok := vm.variables[asm.Value{Type: asm.SymbolType, Val: "abc"}]
	if !ok {
		t.Error("Expected added value to be in variables")
	}
	expected := asm.Value{Type: asm.NumberType, Val: 5}
	if val != expected {
		t.Errorf("Expected %v got %v", expected, val)
	}
}

func TestVmSetVariableBoolean(t *testing.T) {
	vm := VM{}
	vm.variables = map[asm.Value]asm.Value{}

	vm.SetVariableBoolean("abc", true)
	val, ok := vm.variables[asm.Value{Type: asm.SymbolType, Val: "abc"}]
	if !ok {
		t.Error("Expected added value to be in variables")
	}
	expected := asm.Value{Type: asm.BooleanType, Val: true}
	if val != expected {
		t.Errorf("Expected %v got %v", expected, val)
	}

	vm.SetVariableBoolean("abc", false)
	val, ok = vm.variables[asm.Value{Type: asm.SymbolType, Val: "abc"}]
	if !ok {
		t.Error("Expected added value to be in variables")
	}
	expected = asm.Value{Type: asm.BooleanType, Val: false}
	if val != expected {
		t.Errorf("Expected %v got %v", expected, val)
	}
}

func TestVmSetVariableString(t *testing.T) {
	vm := VM{}
	vm.variables = map[asm.Value]asm.Value{}

	vm.SetVariableString("abc", "5")
	val, ok := vm.variables[asm.Value{Type: asm.SymbolType, Val: "abc"}]
	if !ok {
		t.Error("Expected added value to be in variables")
	}
	expected := asm.Value{Type: asm.StringType, Val: "5"}
	if val != expected {
		t.Errorf("Expected %v got %v", expected, val)
	}
}

func TestVmSetVariableNull(t *testing.T) {
	vm := VM{}
	vm.variables = map[asm.Value]asm.Value{}

	vm.SetVariableNull("abc")
	val, ok := vm.variables[asm.Value{Type: asm.SymbolType, Val: "abc"}]
	if !ok {
		t.Error("Expected added value to be in variables")
	}
	expected := asm.Null
	if val != expected {
		t.Errorf("Expected %v got %v", expected, val)
	}
}
func TestVmLoadVariable(t *testing.T) {
	prog := program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.LoadVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "abc"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	vm, _ := New(prog)
	vm.SetVariableNumber("abc", 15)
	vm.Run()
	expectedStack := []asm.Value{{Type: asm.NumberType, Val: 15}}
	if !compareValue(expectedStack, vm.stack) {
		t.Errorf("Expected %v got %v", expectedStack, vm.stack)
	}
	vm.Reset()
	vm.Run()
	expectedStack = []asm.Value{asm.Null}
	if !compareValue(expectedStack, vm.stack) {
		t.Errorf("Expected %v got %v", expectedStack, vm.stack)
	}
}

func TestVmStoreVariable(t *testing.T) {
	prog := program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "a"}},
			{Opcode: asm.PushNull, Arg: asm.Value{}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "b"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "c"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "d"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
		Funcs: map[string][]asm.Type{},
	}
	vm, _ := New(prog)
	vm.Run()
	actual, _ := vm.GetVariable("a")
	expected := asm.Value{Type: asm.BooleanType, Val: false}
	if actual != expected {
		t.Errorf("expected %v got %v", expected, actual)
	}
	actual, _ = vm.GetVariable("b")
	expected = asm.Null
	if actual != expected {
		t.Errorf("expected %v got %v", expected, actual)
	}
	actual, _ = vm.GetVariable("c")
	expected = asm.Value{Type: asm.StringType, Val: "abc"}
	if actual != expected {
		t.Errorf("expected %v got %v", expected, actual)
	}
	actual, _ = vm.GetVariable("d")
	expected = asm.Value{Type: asm.NumberType, Val: 17}
	if actual != expected {
		t.Errorf("expected %v got %v", expected, actual)
	}

	prog = program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.PushBool, Arg: asm.Value{Type: asm.BooleanType, Val: false}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "a"}},
			{Opcode: asm.PushNull, Arg: asm.Value{}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "a"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "a"}},
			{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 17}},
			{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "a"}},
			{Opcode: asm.EndDialogue, Arg: asm.Value{}},
		},
	}
	vm, _ = New(prog)
	vm.Run()
	actual, _ = vm.GetVariable("a")
	expected = asm.Value{Type: asm.NumberType, Val: 17}
	if actual != expected {
		t.Errorf("expected %v got %v", expected, actual)
	}
}
