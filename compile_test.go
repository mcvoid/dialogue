package dialogue

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"
	"testing"

	"github.com/mcvoid/dialogue/internal/program"
	"github.com/mcvoid/dialogue/internal/types/asm"
)

type mockReader struct{}

func (r mockReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func TestCompile(t *testing.T) {
	input := "```\n" +
		"# start\n" +
		"\n" +
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque ut nibh eleifend eros\n" +
		"varius malesuada eget nec nulla. Mauris vitae sem non nibh posuere suscipit sed non\n" +
		"quam. Ut varius diam eros, in pulvinar diam gravida eu.\n" +
		"\n" +
		"Nullam volutpat nisl quis congue maximus. Nam eget imperdiet metus, ac rutrum est. Nam\n" +
		"ac venenatis mi. Nam vehicula neque a porta ornare.\n" +
		"\n" +
		"Vestibulum sapien felis, pharetra sed tellus sed, feugiat auctor justo. Donec tristique\n" +
		"non mi at hendrerit. Duis vel est sodales, finibus turpis eget, finibus augue.\n" +
		"\n" +
		"Donec molestie vulputate vehicula. Suspendisse efficitur, neque eu lacinia\n" +
		"vestibulum, odio neque ullamcorper nibh, quis pretium nisi diam sit amet arcu.\n" +
		"\n"

	r := strings.NewReader(input)
	var b bytes.Buffer

	err := Compile(
		CompilerInput(r),
		CompilerOutput(&b),
	)
	if err != nil {
		t.Errorf("no error expected, got %v", err)
	}
	actual := b.String()
	b = bytes.Buffer{}

	expectedProg := program.Program{
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
	zlibWriter := zlib.NewWriter(&b)
	expectedProg.WriteTo(zlibWriter)
	zlibWriter.Close()
	expected := b.String()

	if expected != actual {
		t.Errorf("expected:\n%v\ngot:\n%v", expected, actual)
	}

	r = strings.NewReader(input)
	b = bytes.Buffer{}

	err = Compile(
		NoCodeFolding,
		NoDeadCodeElimination,
		CompilerInput(r),
		CompilerOutput(&b),
	)
	if err != nil {
		t.Errorf("no error expected, got %v", err)
	}
	actual = b.String()
	b = bytes.Buffer{}

	expectedProg = program.Program{
		Start: 0,
		Code: []asm.Instruction{
			{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "start"}},

			// Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque ut nibh eleifend eros
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque ut nibh eleifend eros"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			// varius malesuada eget nec nulla. Mauris vitae sem non nibh posuere suscipit sed non
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "varius malesuada eget nec nulla. Mauris vitae sem non nibh posuere suscipit sed non"}},
			{Opcode: asm.Concat},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			// quam. Ut varius diam eros, in pulvinar diam gravida eu.
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "quam. Ut varius diam eros, in pulvinar diam gravida eu."}},
			{Opcode: asm.Concat},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			{Opcode: asm.ShowLine},

			// Nullam volutpat nisl quis congue maximus. Nam eget imperdiet metus, ac rutrum est. Nam
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Nullam volutpat nisl quis congue maximus. Nam eget imperdiet metus, ac rutrum est. Nam"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			// ac venenatis mi. Nam vehicula neque a porta ornare.
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "ac venenatis mi. Nam vehicula neque a porta ornare."}},
			{Opcode: asm.Concat},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			{Opcode: asm.ShowLine},

			// Vestibulum sapien felis, pharetra sed tellus sed, feugiat auctor justo. Donec tristique
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Vestibulum sapien felis, pharetra sed tellus sed, feugiat auctor justo. Donec tristique"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			// non mi at hendrerit. Duis vel est sodales, finibus turpis eget, finibus augue.
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "non mi at hendrerit. Duis vel est sodales, finibus turpis eget, finibus augue."}},
			{Opcode: asm.Concat},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			{Opcode: asm.ShowLine},

			// Donec molestie vulputate vehicula. Suspendisse efficitur, neque eu lacinia
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Donec molestie vulputate vehicula. Suspendisse efficitur, neque eu lacinia"}},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},

			// vestibulum, odio neque ullamcorper nibh, quis pretium nisi diam sit amet arcu.
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "vestibulum, odio neque ullamcorper nibh, quis pretium nisi diam sit amet arcu."}},
			{Opcode: asm.Concat},
			{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "\n"}},
			{Opcode: asm.Concat},
			{Opcode: asm.ShowLine},

			{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "start"}},
			{Opcode: asm.EndDialogue},
		},
		Funcs: map[string][]asm.Type{},
	}
	expectedProg.WriteTo(&b)
	expected = b.String()

	if expected != actual {
		t.Errorf("expected:\n%v\ngot:\n%v", expected, actual)
	}

	input = "```\n" +
		"# start\n" +
		"\n" +
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque ut nibh eleifend eros\n" +
		"varius malesuada eget nec nulla. Mauris vitae sem non nibh posuere suscipit sed non\n" +
		"quam. Ut varius diam eros, in pulvinar diam gravida eu.\n" +
		"\n" +
		"Nullam volutpat nisl quis congue maximus. Nam eget imperdiet metus, ac rutrum est. Nam\n" +
		"ac venenatis mi. Nam vehicula neque a porta ornare.\n" +
		"\n" +
		"Vestibulum sapien felis, pharetra sed tellus sed, feugiat auctor justo. Donec tristique\n" +
		"non mi at hendrerit. Duis vel est sodales, finibus turpis eget, finibus augue.\n" +
		"\n" +
		"Donec molestie vulputate vehicula. Suspendisse efficitur, neque eu lacinia\n" +
		"vestibulum, odio neque ullamcorper nibh, quis pretium nisi diam sit amet arcu."

	r = strings.NewReader(input)
	b = bytes.Buffer{}
	err = Compile(
		CompilerInput(r),
		CompilerOutput(&b),
	)
	if err == nil {
		t.Errorf("Expected error got nil")
	}

	b = bytes.Buffer{}
	badReader := mockReader{}
	err = Compile(
		CompilerInput(badReader),
		CompilerOutput(&b),
	)
	if err == nil {
		t.Errorf("Expected error got nil")
	}

	input = "```\n" +
		"# start\n" +
		"\n" +
		"[abc](this is bad text)\n" +
		"\n"

	r = strings.NewReader(input)
	b = bytes.Buffer{}
	err = Compile(
		CompilerInput(r),
		CompilerOutput(&b),
	)
	if err == nil {
		t.Errorf("Expected error got nil")
	}
}
