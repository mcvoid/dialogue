package semantic_analysis

import "github.com/mcvoid/dialogue/internal/types/ast"

type EffectiveType int

const (
	Error EffectiveType = iota
	Variant
	Number
	Boolean
	String
	Null
	Void
)

var astTypeToEffectiveType = map[ast.Type]EffectiveType{
	ast.StringType:  String,
	ast.NumberType:  Number,
	ast.BooleanType: Boolean,
	ast.NullType:    Null,
	ast.SymbolType:  Variant,
}

func TypeCheckScript(script ast.Script) EffectiveType {
	for _, node := range script.Nodes {
		for _, block := range node.Body {
			switch block := block.(type) {
			case ast.Paragraph:
				for _, inline := range block {
					if inlineCode, ok := inline.(ast.InlineCode); ok {
						if t := TypeCheckExpression(inlineCode.Expr); t == Error {
							return Error
						}
					}
				}
			case ast.CodeBlock:
				for _, stmt := range block.Code {
					if t := TypeCheckStatement(stmt, script); t == Error {
						return Error
					}
				}
			case ast.Link:
				// continue
			case ast.Option:
				// continue
			default:
				return Error
			}
		}
	}
	return Void
}

func TypeCheckStatement(stmt ast.Statement, root ast.Script) EffectiveType {
	switch stmt := stmt.(type) {
	case ast.Assignment:
		{
			if TypeCheckExpression(stmt.Val) == Error {
				return Error
			}
			return Void
		}
	case ast.StatementBlock:
		{
			for _, s := range stmt {
				if t := TypeCheckStatement(s, root); t == Error {
					return Error
				}
			}
			return Void
		}
	case ast.Conditional:
		{
			t := TypeCheckExpression(stmt.Cond)
			if t != Boolean && t != Variant {
				return Error
			}
			t = TypeCheckStatement(stmt.Consequent, root)
			if t != Void {
				return Error
			}
			t = TypeCheckStatement(stmt.Alternate, root)
			if t != Void {
				return Error
			}
			return Void
		}
	case ast.Loop:
		{
			t := TypeCheckExpression(stmt.Cond)
			if t != Boolean && t != Variant {
				return Error
			}
			t = TypeCheckStatement(stmt.Consequent, root)
			if t != Void {
				return Error
			}
			return Void
		}
	case ast.FunctionCall:
		{
			prototype := root.Functions[string(stmt.Name)]
			if len(prototype) != len(stmt.Params) {
				return Error
			}
			for i, param := range stmt.Params {
				t := TypeCheckExpression(param)
				if t == Error {
					return Error
				}
				expectedType := astTypeToEffectiveType[prototype[i]]
				if t != Variant && expectedType != t {
					return Error
				}
			}
			return Void
		}
	case ast.GotoNode:
		return Void
	}
	return Error
}

func TypeCheckExpression(expr ast.Expression) EffectiveType {
	switch expr := expr.(type) {
	case ast.BinaryOp:
		return TypeCheckBinary(expr)
	case ast.UnaryOp:
		return TypeCheckUnary(expr)
	case ast.Literal:
		return TypeCheckLiteral(expr)
	}
	return Error
}

func TypeCheckBinary(op ast.BinaryOp) EffectiveType {
	switch op.Operator {
	case ast.AddOp:
		fallthrough
	case ast.SubOp:
		fallthrough
	case ast.MulOp:
		fallthrough
	case ast.DivOp:
		fallthrough
	case ast.ModOp:
		{
			t := TypeCheckExpression(op.LeftArg)
			if t != Number && t != Variant {
				return Error
			}
			t = TypeCheckExpression(op.RightArg)
			if t != Number && t != Variant {
				return Error
			}
			return Number
		}
	case ast.GtOp:
		fallthrough
	case ast.LtOp:
		fallthrough
	case ast.GteOp:
		fallthrough
	case ast.LteOp:
		{
			t := TypeCheckExpression(op.LeftArg)
			if t != Number && t != Variant {
				return Error
			}
			t = TypeCheckExpression(op.RightArg)
			if t != Number && t != Variant {
				return Error
			}
			return Boolean
		}
	case ast.AndOp:
		fallthrough
	case ast.OrOp:
		{
			t := TypeCheckExpression(op.LeftArg)
			if t != Boolean && t != Variant {
				return Error
			}
			t = TypeCheckExpression(op.RightArg)
			if t != Boolean && t != Variant {
				return Error
			}
			return Boolean
		}
	case ast.EqOp:
		fallthrough
	case ast.NeqOp:
		// eq can take any two types for arguments and always produces boolean
		return Boolean
	case ast.ConcatOp:
		// concat can take any two types for arguments and always produces string
		return String
	}

	return Error
}

func TypeCheckUnary(op ast.UnaryOp) EffectiveType {
	switch op.Operator {
	case ast.IncOp:
		fallthrough
	case ast.DecOp:
		{
			t := TypeCheckExpression(op.Arg)
			if t == Number || t == Variant {
				return Number
			}
			return Error
		}
	case ast.NotOp:
		{
			t := TypeCheckExpression(op.Arg)
			if t == Boolean || t == Variant {
				return Boolean
			}
			return Error
		}
	case ast.NegOp:
		{
			t := TypeCheckExpression(op.Arg)
			if t == Number || t == Variant {
				return Number
			}
			return Error
		}
	}
	return Error
}

func TypeCheckLiteral(lit ast.Literal) EffectiveType {
	switch lit.Type {
	case ast.StringType:
		if _, ok := lit.Val.(string); !ok {
			return Error
		}
		return String
	case ast.NullType:
		if lit.Val != nil {
			return Error
		}
		return Null
	case ast.NumberType:
		if _, ok := lit.Val.(int); !ok {
			return Error
		}
		return Number
	case ast.BooleanType:
		if _, ok := lit.Val.(bool); !ok {
			return Error
		}
		return Boolean
	case ast.SymbolType:
		if _, ok := lit.Val.(string); !ok {
			return Error
		}
		// we don't know the values of variables until runtime
		return Variant
	}
	return Error
}
