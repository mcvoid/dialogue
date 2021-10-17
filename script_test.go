package dialogue

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"testing"

	"github.com/mcvoid/dialogue/internal/program"
	"github.com/mcvoid/dialogue/internal/types/asm"
)

func TestFromReader(t *testing.T) {
	expected := program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "start"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque ut nibh eleifend eros varius malesuada eget nec nulla. Mauris vitae sem non nibh posuere suscipit sed non quam. Ut varius diam eros, in pulvinar diam gravida eu."}},
			{Opcode: asm.ShowLine},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Nullam volutpat nisl quis congue maximus. Nam eget imperdiet metus, ac rutrum est. Nam ac venenatis mi. Nam vehicula neque a porta ornare."}},
			{Opcode: asm.ShowLine},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Vestibulum sapien felis, pharetra sed tellus sed, feugiat auctor justo. Donec tristique non mi at hendrerit. Duis vel est sodales, finibus turpis eget, finibus augue."}},
			{Opcode: asm.ShowLine},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Donec molestie vulputate vehicula. Suspendisse efficitur, neque eu lacinia vestibulum, odio neque ullamcorper nibh, quis pretium nisi diam sit amet arcu."}},
			{Opcode: asm.ShowLine},
			{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "start"}},
			{Opcode: asm.EndDialogue},
		},
	}
	// test uncompressed
	b := bytes.Buffer{}
	expected.WriteTo(&b)
	actual, err := FromReader(
		ScriptInput(&b),
		IsUncompressed(),
	)
	if err != nil {
		t.Fatalf("no error expected, got %v", err)
	}
	if len(expected.Code) != len(actual.program.Code) {
		t.Fatalf("expected %v got %v", expected, actual.program)
	}

	for i := range expected.Code {
		if expected.Code[i] != actual.program.Code[i] {
			t.Fatalf("expected %v got %v", expected, actual.program)
		}
	}

	// test compressed
	b = bytes.Buffer{}
	wc := zlib.NewWriter(&b)
	expected.WriteTo(wc)
	wc.Close()

	actual, err = FromReader(
		ScriptInput(&b),
	)
	if err != nil {
		t.Fatalf("no error expected, got %v", err)
	}

	if len(expected.Code) != len(actual.program.Code) {
		t.Fatalf("expected %v got %v", expected, actual.program)
	}

	for i := range expected.Code {
		if expected.Code[i] != actual.program.Code[i] {
			t.Fatalf("expected %v got %v", expected, actual.program)
		}
	}

	// read error compressed
	b = bytes.Buffer{}
	_, err = FromReader(
		ScriptInput(&b),
	)
	if err == nil {
		t.Fatalf("error expected, got nil")
	}

	// read error uncompressed
	b = bytes.Buffer{}
	_, err = FromReader(
		ScriptInput(&b),
		IsUncompressed(),
	)
	if err == nil {
		t.Fatalf("error expected, got nil")
	}

	// read error compressed
	b = bytes.Buffer{}
	wc = zlib.NewWriter(&b)
	fmt.Fprintf(wc, "nonsense 21346798yuihejkdfsx90uijk@$#QERWDFSG^RTEY")
	wc.Close()
	_, err = FromReader(
		ScriptInput(&b),
	)
	if err == nil {
		t.Fatalf("error expected, got nil")
	}
}

type mockHandler struct {
	lastNodeEntered string
	numChoices      int
	lastNodeExited  string
	lastLineShown   string
	scriptEnded     bool
	returnValue     ExecutionType
}

func TestNew(t *testing.T) {
	script := Script{
		program: program.Program{
			Start: 0,
			Code: []asm.Instruction{
				{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "start"}},
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Lorem ipsum dolor sit amet"}},
				{Opcode: asm.ShowLine},
				{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "start"}},
				{Opcode: asm.EndDialogue},
			},
		},
	}

	h := mockHandler{}

	var hf HandlerFunc = func(m Message) ExecutionType {
		switch m.Type {
		case ShowLineType:
			h.lastLineShown = m.Line
		case ShowChoiceType:
			h.numChoices = len(m.Options)
		case EnterNodeType:
			h.lastNodeEntered = m.NodeEntered
		case ExitNodeType:
			h.lastNodeExited = m.NodeExited
		case EndScriptType:
			h.scriptEnded = true
		}
		return h.returnValue
	}

	_, err := script.New(nil)
	if err == nil {
		t.Fatalf("error expected got nil")
	}

	h = mockHandler{returnValue: Continue}
	proc, err := script.New(hf)
	if err != nil {
		t.Errorf("no error expected, got %v", err)
	}
	proc.Start()

	if h.lastLineShown != "Lorem ipsum dolor sit amet" {
		t.Error("wrong last line shown")
	}
	if h.lastNodeEntered != "start" {
		t.Errorf("wrong last node entered")
	}
	if h.lastNodeExited != "start" {
		t.Errorf("wrong last node exited")
	}
	if !h.scriptEnded {
		t.Errorf("Expected script exited")
	}

	script = Script{
		program: program.Program{
			Start: 0,
			Code: []asm.Instruction{
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
				{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 7}},
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
				{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 7}},
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
				{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 7}},
				{Opcode: asm.ShowChoice},
				{Opcode: asm.EndDialogue},
			},
		},
	}
	h = mockHandler{returnValue: Continue}
	proc, err = script.New(hf)
	if err != nil {
		t.Errorf("no error expected, got %v", err)
	}
	proc.Start()
	if h.numChoices != 3 {
		t.Errorf("expected show choice handler invoked")
	}
	if h.scriptEnded {
		t.Errorf("vm should not continue on show choice")
	}
	proc.ChooseAndResume(0)

	if !h.scriptEnded {
		t.Errorf("vm did not continue on choose and resume")
	}

	script = Script{
		program: program.Program{
			Start: 0,
			Code: []asm.Instruction{
				{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "abc"}},
				{Opcode: asm.ShowLine},
				{Opcode: asm.EndDialogue},
			},
		},
	}
	h = mockHandler{returnValue: Pause}
	proc, err = script.New(hf)
	if err != nil {
		t.Errorf("no error expected, got %v", err)
	}
	proc.Start()
	if h.lastLineShown != "abc" {
		t.Error("wrong last line shown")
	}
	if h.scriptEnded {
		t.Errorf("vm should not continue on suspend")
	}
	proc.Resume()

	if !h.scriptEnded {
		t.Errorf("vm did not continue on resume")
	}
}
