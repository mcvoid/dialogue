package asm

import (
	"encoding/json"
	"fmt"
)

type (
	Opcode string
	Type   string
	Value  struct {
		Type Type
		Val  interface{}
	}
	Instruction struct {
		Opcode Opcode
		Arg    Value
	}
)

const (
	PushString         Opcode = "PushString"
	PushNumber         Opcode = "PushNumber"
	PushBool           Opcode = "PushBool"
	PushNull           Opcode = "PushNull"
	PopValue           Opcode = "PopValue"
	Concat             Opcode = "Concat"
	And                Opcode = "And"
	Or                 Opcode = "Or"
	Not                Opcode = "Not"
	Equal              Opcode = "Equal"
	NotEqual           Opcode = "NotEqual"
	GreaterThan        Opcode = "GreaterThan"
	Lessthan           Opcode = "Lessthan"
	GreaterThanOrEqual Opcode = "GreaterThanOrEqual"
	LessthanOrEqual    Opcode = "LessthanOrEqual"
	Negative           Opcode = "Negative"
	Add                Opcode = "Add"
	Subtract           Opcode = "Subtract"
	Multiply           Opcode = "Multiply"
	Divide             Opcode = "Divide"
	Modulo             Opcode = "Modulo"
	Increment          Opcode = "Increment"
	Decrement          Opcode = "Decrement"
	LoadVariable       Opcode = "LoadVariable"
	StoreVariable      Opcode = "StoreVariable"
	ShowLine           Opcode = "ShowLine"
	Jump               Opcode = "Jump"
	JumpIfFalse        Opcode = "JumpIfFalse"
	PushChoice         Opcode = "PushChoice"
	ShowChoice         Opcode = "ShowChoice"
	EnterNode          Opcode = "EnterNode"
	ExitNode           Opcode = "ExitNode"
	EndDialogue        Opcode = "EndDialogue"
	Call               Opcode = "Call"
)

const (
	BooleanType Type = "boolean"
	NumberType  Type = "number"
	StringType  Type = "string"
	NullType    Type = "null"
	SymbolType  Type = "symbol"
)

var (
	Null  = Value{NullType, nil}
	True  = Value{BooleanType, true}
	False = Value{BooleanType, false}
)

var (
	unaryOpcodes = map[Opcode]bool{
		PushString:    true,
		PushNumber:    true,
		PushBool:      true,
		Jump:          true,
		JumpIfFalse:   true,
		LoadVariable:  true,
		StoreVariable: true,
		PushChoice:    true,
		EnterNode:     true,
		ExitNode:      true,
		Call:          true,
	}
	nullaryOpcodes = map[Opcode]bool{
		PushNull:    true,
		PopValue:    true,
		Concat:      true,
		And:         true,
		Or:          true,
		Not:         true,
		Equal:       true,
		NotEqual:    true,
		GreaterThan: true,
		Lessthan:    true,
		Negative:    true,
		Add:         true,
		Subtract:    true,
		Multiply:    true,
		Divide:      true,
		Modulo:      true,
		Increment:   true,
		Decrement:   true,
		ShowLine:    true,
		ShowChoice:  true,
		EndDialogue: true,
	}
)

func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{v.Type, v.Val})
}

func (v *Value) UnmarshalJSON(b []byte) error {
	template := []interface{}{&v.Type, &v.Val}
	if err := json.Unmarshal(b, &template); err != nil {
		return err
	}
	switch v.Type {
	case BooleanType:
		if _, ok := v.Val.(bool); !ok {
			return fmt.Errorf("type boolean must be encoded in JSON boolean value")
		}
	case StringType:
		if _, ok := v.Val.(string); !ok {
			return fmt.Errorf("type string must be encoded in JSON string value")
		}
	case SymbolType:
		if _, ok := v.Val.(string); !ok {
			return fmt.Errorf("type symbol must be encoded in JSON string value")
		}
	case NumberType:
		num, ok := v.Val.(float64)
		if !ok {
			return fmt.Errorf("type number must be encoded in JSON number value")
		}
		v.Val = int(num)
	default:
		return fmt.Errorf("invalid type %v", v.Type)
	}
	return nil
}

func (i Instruction) MarshalJSON() ([]byte, error) {
	if nullaryOpcodes[i.Opcode] {
		return json.Marshal([]interface{}{i.Opcode})
	}
	if unaryOpcodes[i.Opcode] {
		return json.Marshal([]interface{}{i.Opcode, i.Arg})
	}
	return nil, fmt.Errorf("invalid opcode")
}

func (i *Instruction) UnmarshalJSON(b []byte) error {
	template := []interface{}{&i.Opcode}
	if err := json.Unmarshal(b, &template); err != nil {
		return err
	}

	if nullaryOpcodes[i.Opcode] {
		return nil
	}
	if unaryOpcodes[i.Opcode] {
		template = []interface{}{&i.Opcode, &i.Arg}
		err := json.Unmarshal(b, &template)
		if err != nil {
			return err
		}
		if len(template) < 2 {
			return fmt.Errorf("wrong arity")
		}
		return nil
	}
	return fmt.Errorf("invalid opcode")
}
