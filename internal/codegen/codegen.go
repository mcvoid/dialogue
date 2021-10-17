package codegen

import (
	"fmt"

	"github.com/mcvoid/dialogue/internal/program"
	"github.com/mcvoid/dialogue/internal/types/asm"
	"github.com/mcvoid/dialogue/internal/types/ast"
)

func Codegen(n ast.Script) (program.Program, error) {
	return generateScript(n)
}

type CodegenContext struct {
	SymbolTable        map[ast.Symbol]int
	BackreferenceTable map[int]ast.Symbol
	Cursor             int
	Code               []asm.Instruction
	CurrentNode        ast.Symbol
}

func (ctx *CodegenContext) AddInstruction(instr asm.Instruction) {
	ctx.Code = append(ctx.Code, instr)
	ctx.Cursor++
}

func (ctx *CodegenContext) AddSymbol(s ast.Symbol) {
	ctx.SymbolTable[s] = ctx.Cursor
	ctx.CurrentNode = s
}

func (ctx *CodegenContext) AddBackRef(s ast.Symbol) {
	ctx.BackreferenceTable[ctx.Cursor] = s
}

func generateScript(n ast.Script) (program.Program, error) {
	ctx := CodegenContext{
		SymbolTable:        map[ast.Symbol]int{},
		BackreferenceTable: map[int]ast.Symbol{},
		Cursor:             0,
		Code:               []asm.Instruction{},
	}

	if len(n.Nodes) == 0 {
		ctx.AddInstruction(asm.Instruction{Opcode: asm.EndDialogue})
		return program.Program{
			Start: 0,
			Code:  ctx.Code,
			Funcs: map[string][]asm.Type{},
		}, nil
	}

	for _, block := range n.Nodes {
		generateBlock(&ctx, block)
	}

	for i, sym := range ctx.BackreferenceTable {
		dest, ok := ctx.SymbolTable[sym]
		if !ok {
			return program.Program{}, fmt.Errorf("unknown symbol %v", sym)
		}
		ctx.Code[i].Arg.Val = dest
	}

	return program.Program{
		Start: 0,
		Code:  ctx.Code,
		Funcs: map[string][]asm.Type{},
	}, nil
}

func generateBlock(ctx *CodegenContext, n ast.Node) {
	ctx.AddSymbol(n.Name)
	var nodeString string = string(n.Name)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.EnterNode,
		Arg:    asm.Value{Type: asm.SymbolType, Val: nodeString},
	})

	for _, el := range n.Body {
		GenerateBlockElement(ctx, el)
	}

	nodeString = string(ctx.CurrentNode)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.ExitNode,
		Arg:    asm.Value{Type: asm.SymbolType, Val: nodeString},
	})
	ctx.AddInstruction(asm.Instruction{Opcode: asm.EndDialogue})
}

func GenerateBlockElement(ctx *CodegenContext, n ast.BlockElement) {
	switch n := n.(type) {
	case ast.Paragraph:
		GenerateParagraph(ctx, n)
	case ast.Link:
		GenerateLink(ctx, n)
	case ast.Option:
		GenerateOption(ctx, n)
	case ast.CodeBlock:
		GenerateCodeBlock(ctx, n)
	}
}

func GenerateParagraph(ctx *CodegenContext, n ast.Paragraph) {
	for i, inline := range n {
		GenerateInline(ctx, inline)
		if i > 0 {
			ctx.AddInstruction(asm.Instruction{Opcode: asm.Concat})
		}
	}
	ctx.AddInstruction(asm.Instruction{Opcode: asm.ShowLine})
}

func GenerateLink(ctx *CodegenContext, n ast.Link) {
	var nodeName string = string(ctx.CurrentNode)
	GenerateInline(ctx, n.Text)
	ctx.AddInstruction(asm.Instruction{Opcode: asm.ShowLine})
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.ExitNode,
		Arg:    asm.Value{Type: asm.SymbolType, Val: nodeName},
	})
	ctx.AddBackRef(n.Dest)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.Jump,
		Arg:    asm.Value{Type: asm.NumberType, Val: 0},
	})
}

func GenerateOption(ctx *CodegenContext, n ast.Option) {
	var nodeName string = string(ctx.CurrentNode)
	for _, link := range n {
		GenerateInline(ctx, link.Text)
		ctx.AddBackRef(link.Dest)
		ctx.AddInstruction(asm.Instruction{
			Opcode: asm.PushChoice,
			Arg:    asm.Value{Type: asm.NumberType, Val: 0},
		})
	}
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.ExitNode,
		Arg:    asm.Value{Type: asm.SymbolType, Val: nodeName},
	})
	ctx.AddInstruction(asm.Instruction{Opcode: asm.ShowChoice})
}

func GenerateCodeBlock(ctx *CodegenContext, n ast.CodeBlock) {
	for _, stmt := range n.Code {
		GenerateStatement(ctx, stmt)
	}
}

func GenerateInline(ctx *CodegenContext, n ast.Inline) {
	switch n := n.(type) {
	case ast.Text:
		GenerateText(ctx, n)
	case ast.InlineCode:
		GenerateInlineCode(ctx, n)
	}
}

func GenerateText(ctx *CodegenContext, n ast.Text) {
	str := string(n)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.PushString,
		Arg:    asm.Value{Type: asm.StringType, Val: str},
	})
}

func GenerateInlineCode(ctx *CodegenContext, n ast.InlineCode) {
	GenerateExpression(ctx, n.Expr)
}

func GenerateStatement(ctx *CodegenContext, n ast.Statement) {
	switch n := n.(type) {
	case ast.Assignment:
		GenerateAssignment(ctx, n)
	case ast.FunctionCall:
		GenerateFunctionCall(ctx, n)
	case ast.StatementBlock:
		GenerateStatementBlock(ctx, n)
	case ast.GotoNode:
		GenerateGotoNode(ctx, n)
	case ast.Conditional:
		GenerateIf(ctx, n)
	case ast.Loop:
		GenerateLoop(ctx, n)
	case ast.InfiniteLoop:
		GenerateInfiniteLoop(ctx, n)
	}
}

func GenerateAssignment(ctx *CodegenContext, n ast.Assignment) {
	GenerateExpression(ctx, n.Val)
	var varName string = string(n.Name)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.StoreVariable,
		Arg:    asm.Value{Type: asm.SymbolType, Val: varName},
	})
}

func GenerateStatementBlock(ctx *CodegenContext, n ast.StatementBlock) {
	for _, statement := range n {
		GenerateStatement(ctx, statement)
	}
}

func GenerateFunctionCall(ctx *CodegenContext, n ast.FunctionCall) {
	for _, arg := range n.Params {
		GenerateExpression(ctx, arg)
	}

	var fnName string = string(n.Name)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.Call,
		Arg:    asm.Value{Type: asm.SymbolType, Val: fnName},
	})
}

func GenerateExpression(ctx *CodegenContext, n ast.Expression) {
	switch n := n.(type) {
	case ast.BinaryOp:
		GenerateBinaryOp(ctx, n)
	case ast.UnaryOp:
		GenerateUnaryOp(ctx, n)
	case ast.Literal:
		GenerateLiteral(ctx, n)
	}
}

func GenerateLiteral(ctx *CodegenContext, n ast.Literal) {
	switch n.Type {
	case ast.NullType:
		ctx.AddInstruction(asm.Instruction{
			Opcode: asm.PushNull,
			Arg:    asm.Null,
		})
	case ast.NumberType:
		ctx.AddInstruction(asm.Instruction{
			Opcode: asm.PushNumber,
			Arg:    asm.Value{Type: asm.NumberType, Val: n.Val},
		})
	case ast.StringType:
		ctx.AddInstruction(asm.Instruction{
			Opcode: asm.PushString,
			Arg:    asm.Value{Type: asm.StringType, Val: n.Val},
		})
	case ast.BooleanType:
		ctx.AddInstruction(asm.Instruction{
			Opcode: asm.PushBool,
			Arg:    asm.Value{Type: asm.BooleanType, Val: n.Val},
		})
	case ast.SymbolType:
		symNode := n.Val.(ast.Symbol)
		var varName string = string(symNode)
		ctx.AddInstruction(asm.Instruction{
			Opcode: asm.LoadVariable,
			Arg:    asm.Value{Type: asm.SymbolType, Val: varName},
		})
	}
}

func GenerateGotoNode(ctx *CodegenContext, n ast.GotoNode) {
	var curString string = string(ctx.CurrentNode)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.ExitNode,
		Arg:    asm.Value{Type: asm.SymbolType, Val: curString},
	})
	ctx.AddBackRef(n.Name)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.Jump,
		Arg:    asm.Value{Type: asm.NumberType, Val: 0},
	})
}

func GenerateIf(ctx *CodegenContext, n ast.Conditional) {
	if block, ok := n.Alternate.(ast.StatementBlock); ok && len(block) == 0 {
		GenerateNakedIf(ctx, n)
		return
	}
	GenerateExpression(ctx, n.Cond)
	cond := ctx.Cursor
	ctx.AddInstruction(asm.Instruction{Opcode: asm.JumpIfFalse})
	GenerateStatement(ctx, n.Consequent)
	consEnd := ctx.Cursor
	ctx.AddInstruction(asm.Instruction{Opcode: asm.Jump})
	ctx.Code[cond] = asm.Instruction{
		Opcode: asm.JumpIfFalse,
		Arg:    asm.Value{Type: asm.NumberType, Val: ctx.Cursor},
	}
	GenerateStatement(ctx, n.Alternate)
	ctx.Code[consEnd] = asm.Instruction{
		Opcode: asm.Jump,
		Arg:    asm.Value{Type: asm.NumberType, Val: ctx.Cursor},
	}
}

func GenerateNakedIf(ctx *CodegenContext, n ast.Conditional) {
	GenerateExpression(ctx, n.Cond)
	cond := ctx.Cursor
	ctx.AddInstruction(asm.Instruction{Opcode: asm.JumpIfFalse})
	GenerateStatement(ctx, n.Consequent)
	ctx.Code[cond] = asm.Instruction{
		Opcode: asm.JumpIfFalse,
		Arg:    asm.Value{Type: asm.NumberType, Val: ctx.Cursor},
	}
}

func GenerateLoop(ctx *CodegenContext, n ast.Loop) {
	loopStart := ctx.Cursor
	GenerateExpression(ctx, n.Cond)
	cond := ctx.Cursor
	ctx.AddInstruction(asm.Instruction{Opcode: asm.JumpIfFalse})
	GenerateStatement(ctx, n.Consequent)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.Jump,
		Arg:    asm.Value{Type: asm.NumberType, Val: loopStart},
	})
	ctx.Code[cond] = asm.Instruction{
		Opcode: asm.JumpIfFalse,
		Arg:    asm.Value{Type: asm.NumberType, Val: ctx.Cursor},
	}
}

func GenerateInfiniteLoop(ctx *CodegenContext, n ast.InfiniteLoop) {
	loopStart := ctx.Cursor
	GenerateStatement(ctx, n.Consequent)
	ctx.AddInstruction(asm.Instruction{
		Opcode: asm.Jump,
		Arg:    asm.Value{Type: asm.NumberType, Val: loopStart},
	})
}

func GenerateUnaryOp(ctx *CodegenContext, n ast.UnaryOp) {
	instr := asm.Instruction{}
	switch n.Operator {
	case ast.NotOp:
		instr = asm.Instruction{Opcode: asm.Not}
	case ast.IncOp:
		instr = asm.Instruction{Opcode: asm.Increment}
	case ast.DecOp:
		instr = asm.Instruction{Opcode: asm.Decrement}
	case ast.NegOp:
		instr = asm.Instruction{Opcode: asm.Negative}
	}

	GenerateExpression(ctx, n.Arg)
	ctx.AddInstruction(instr)
}

func GenerateBinaryOp(ctx *CodegenContext, n ast.BinaryOp) {
	instr := asm.Instruction{}
	switch n.Operator {
	case ast.AddOp:
		instr = asm.Instruction{Opcode: asm.Add}
	case ast.SubOp:
		instr = asm.Instruction{Opcode: asm.Subtract}
	case ast.MulOp:
		instr = asm.Instruction{Opcode: asm.Multiply}
	case ast.DivOp:
		instr = asm.Instruction{Opcode: asm.Divide}
	case ast.ModOp:
		instr = asm.Instruction{Opcode: asm.Modulo}
	case ast.GtOp:
		instr = asm.Instruction{Opcode: asm.GreaterThan}
	case ast.GteOp:
		instr = asm.Instruction{Opcode: asm.GreaterThanOrEqual}
	case ast.LtOp:
		instr = asm.Instruction{Opcode: asm.Lessthan}
	case ast.LteOp:
		instr = asm.Instruction{Opcode: asm.LessthanOrEqual}
	case ast.EqOp:
		instr = asm.Instruction{Opcode: asm.Equal}
	case ast.NeqOp:
		instr = asm.Instruction{Opcode: asm.NotEqual}
	case ast.AndOp:
		instr = asm.Instruction{Opcode: asm.And}
	case ast.OrOp:
		instr = asm.Instruction{Opcode: asm.Or}
	case ast.ConcatOp:
		instr = asm.Instruction{Opcode: asm.Concat}
	}

	GenerateExpression(ctx, n.LeftArg)
	GenerateExpression(ctx, n.RightArg)
	ctx.AddInstruction(instr)
}
