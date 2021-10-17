package semantic_analysis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mcvoid/dialogue/internal/types/ast"
)

var ws = regexp.MustCompile("[ \t\n]+")

func ConstantFoldScript(node ast.Script) ast.Script {
	foldedScript := ast.Script{
		Functions: node.Functions,
		Nodes:     []ast.Node{},
	}
	for _, n := range node.Nodes {
		foldedScript.Nodes = append(foldedScript.Nodes, ConstantFoldNode(n))
	}
	return foldedScript
}

func ConstantFoldNode(node ast.Node) ast.Node {
	foldedBlocks := []ast.BlockElement{}

	for _, block := range node.Body {
		switch block := block.(type) {
		case ast.Paragraph:
			foldedBlocks = append(foldedBlocks, ConstantFoldParagraph(block))
		case ast.CodeBlock:
			foldedBlocks = append(foldedBlocks, ConstantFoldCodeBlock(block))
		default:
			foldedBlocks = append(foldedBlocks, block)
		}
	}

	return ast.Node{
		Name: node.Name,
		Body: foldedBlocks,
	}
}

func ConstantFoldCodeBlock(node ast.CodeBlock) ast.CodeBlock {
	foldedStmts := []ast.Statement{}

	for _, stmt := range node.Code {
		foldedStmts = append(foldedStmts, ConstantFoldStatement(stmt))
	}

	return ast.CodeBlock{
		Code: foldedStmts,
	}
}

func ConstantFoldStatement(node ast.Statement) ast.Statement {
	switch node := node.(type) {
	case ast.Assignment:
		{
			expr, _ := ConstantFoldExpression(node.Val)
			return ast.Assignment{
				Name: node.Name,
				Val:  expr,
			}
		}
	case ast.StatementBlock:
		{
			newBlock := ast.StatementBlock{}
			for _, stmt := range node {
				newBlock = append(newBlock, ConstantFoldStatement(stmt))
			}
			return newBlock
		}
	case ast.Conditional:
		{
			expr, _ := ConstantFoldExpression(node.Cond)
			folded := ast.Conditional{
				Cond:       expr,
				Consequent: ConstantFoldStatement(node.Consequent),
				Alternate:  ConstantFoldStatement(node.Alternate),
			}

			if block, ok := folded.Consequent.(ast.StatementBlock); ok && len(block) == 0 {
				expr, _ = ConstantFoldExpression(ast.UnaryOp{
					Operator: ast.NotOp,
					Arg:      expr,
				})
				folded = ast.Conditional{
					Cond:       expr,
					Consequent: folded.Alternate,
					Alternate:  ast.StatementBlock{},
				}
			}

			return folded
		}
	case ast.Loop:
		{
			expr, _ := ConstantFoldExpression(node.Cond)

			return ast.Loop{
				Cond:       expr,
				Consequent: ConstantFoldStatement(node.Consequent),
			}
		}
	case ast.FunctionCall:
		{
			foldedArgs := []ast.Expression{}
			for _, arg := range node.Params {
				expr, _ := ConstantFoldExpression(arg)
				foldedArgs = append(foldedArgs, expr)
			}
			return ast.FunctionCall{
				Name:   node.Name,
				Params: foldedArgs,
			}
		}
	}
	return node
}

func ConstantFoldParagraph(node ast.Paragraph) (foldedNode ast.Paragraph) {
	lastFoldedIndex := 0
	for i, inline := range node {
		if i == 0 {
			// no previous inline to fold into
			foldedNode = ast.Paragraph{inline}
			lastFoldedIndex = 0
			continue
		}

		var thisConst ast.Text
		switch inline := inline.(type) {
		case ast.InlineCode:
			foldedExpr, isConst := ConstantFoldExpression(inline.Expr)
			if !isConst {
				// don't fold
				foldedNode = append(foldedNode, ast.InlineCode{Expr: foldedExpr})
				lastFoldedIndex++
				continue
			}
			thisConst = ast.Text(fmt.Sprintf("%v", foldedExpr.(ast.Literal).Val))
		case ast.Text:
			thisConst = inline
		}

		// if we reached this point, the current inline is a constant.
		// all previously added constants are text
		// so if the previous node is text, it's const
		// and we can concatenate the two
		switch prev := foldedNode[lastFoldedIndex].(type) {
		case ast.InlineCode:
			// previous node isn't const, just add the inline
			foldedNode = append(foldedNode, thisConst)
			lastFoldedIndex++
			continue
		case ast.Text:
			// both previous and current nodes are constant and text
			// combine them

			// also let's collapse extra whitespace like HTML
			newStr := string(prev + thisConst)
			newStr = ws.ReplaceAllString(newStr, " ")
			if i == len(node)-1 {
				newStr = strings.TrimRight(newStr, " ")
			}
			foldedNode[lastFoldedIndex] = ast.Text(newStr)
		}
	}

	return foldedNode
}

func ConstantFoldExpression(node ast.Expression) (foldedNode ast.Expression, isConstExpr bool) {
	switch node := node.(type) {
	case ast.BinaryOp:
		return ConstantFoldBinaryOperation(node)
	case ast.UnaryOp:
		return ConstantFoldUnaryOperation(node)
	case ast.Literal:
		return ConstantFoldLiteral(node)
	}
	return nil, false
}

func ConstantFoldLiteral(node ast.Literal) (foldedNode ast.Literal, isConstExpr bool) {
	// variables aren't constants
	if node.Type == ast.SymbolType {
		return node, false
	}
	// constants are constants
	return node, true
}

func ConstantFoldUnaryOperation(node ast.UnaryOp) (foldedNode ast.Expression, isConstExpr bool) {
	arg, argIsConst := ConstantFoldExpression(node.Arg)

	if argIsConst {
		switch node.Operator {
		case ast.IncOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  arg.(ast.Literal).Val.(int) + 1,
			}, true
		case ast.DecOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  arg.(ast.Literal).Val.(int) - 1,
			}, true
		case ast.NotOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  !arg.(ast.Literal).Val.(bool),
			}, true
		case ast.NegOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  -(arg.(ast.Literal).Val.(int)),
			}, true
		}
	}

	// not affects equality and other nots
	if node.Operator == ast.NotOp {
		if binarg, ok := arg.(ast.BinaryOp); ok {
			switch binarg.Operator {
			case ast.EqOp:
				return ast.BinaryOp{
					Operator: ast.NeqOp,
					LeftArg:  binarg.LeftArg,
					RightArg: binarg.RightArg,
				}, false
			case ast.GtOp:
				return ast.BinaryOp{
					Operator: ast.LteOp,
					LeftArg:  binarg.LeftArg,
					RightArg: binarg.RightArg,
				}, false
			case ast.LtOp:
				return ast.BinaryOp{
					Operator: ast.GteOp,
					LeftArg:  binarg.LeftArg,
					RightArg: binarg.RightArg,
				}, false
			case ast.GteOp:
				return ast.BinaryOp{
					Operator: ast.LtOp,
					LeftArg:  binarg.LeftArg,
					RightArg: binarg.RightArg,
				}, false
			case ast.LteOp:
				return ast.BinaryOp{
					Operator: ast.GtOp,
					LeftArg:  binarg.LeftArg,
					RightArg: binarg.RightArg,
				}, false
			}
		}
		if unarg, ok := arg.(ast.UnaryOp); ok {
			switch unarg.Operator {
			case ast.NotOp:
				return unarg.Arg, false
			}
		}
	}

	// neg affects other negs
	if node.Operator == ast.NegOp {
		if unarg, ok := arg.(ast.UnaryOp); ok {
			switch unarg.Operator {
			case ast.NegOp:
				return unarg.Arg, false
			}
		}
	}

	return ast.UnaryOp{
		Operator: node.Operator,
		Arg:      arg,
	}, false
}

func ConstantFoldBinaryOperation(node ast.BinaryOp) (foldedNode ast.Expression, isConstExpr bool) {
	left, leftIsConst := ConstantFoldExpression(node.LeftArg)
	right, rightIsConst := ConstantFoldExpression(node.RightArg)

	if leftIsConst && rightIsConst {
		switch node.Operator {
		case ast.AddOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  left.(ast.Literal).Val.(int) + right.(ast.Literal).Val.(int),
			}, true
		case ast.SubOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  left.(ast.Literal).Val.(int) - right.(ast.Literal).Val.(int),
			}, true
		case ast.MulOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  left.(ast.Literal).Val.(int) * right.(ast.Literal).Val.(int),
			}, true
		case ast.DivOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  left.(ast.Literal).Val.(int) / right.(ast.Literal).Val.(int),
			}, true
		case ast.ModOp:
			return ast.Literal{
				Type: ast.NumberType,
				Val:  left.(ast.Literal).Val.(int) % right.(ast.Literal).Val.(int),
			}, true
		case ast.GtOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(int) > right.(ast.Literal).Val.(int),
			}, true
		case ast.GteOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(int) >= right.(ast.Literal).Val.(int),
			}, true
		case ast.LtOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(int) < right.(ast.Literal).Val.(int),
			}, true
		case ast.LteOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(int) <= right.(ast.Literal).Val.(int),
			}, true
		case ast.EqOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(int) == right.(ast.Literal).Val.(int),
			}, true
		case ast.NeqOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(int) != right.(ast.Literal).Val.(int),
			}, true
		case ast.AndOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(bool) && right.(ast.Literal).Val.(bool),
			}, true
		case ast.OrOp:
			return ast.Literal{
				Type: ast.BooleanType,
				Val:  left.(ast.Literal).Val.(bool) || right.(ast.Literal).Val.(bool),
			}, true
		case ast.ConcatOp:
			return ast.Literal{
				Type: ast.StringType,
				Val:  fmt.Sprintf("%v%v", left.(ast.Literal).Val, right.(ast.Literal).Val),
			}, true
		}
	}

	// comparisons of a variable to itself
	// are known at compile time
	if leftLiteral, isLeftSym := left.(ast.Literal); isLeftSym {
		if rightLiteral, isRightSym := right.(ast.Literal); isRightSym {
			if leftLiteral.Type == ast.SymbolType &&
				rightLiteral.Type == ast.SymbolType &&
				leftLiteral.Val == rightLiteral.Val {
				switch node.Operator {
				case ast.EqOp:
					// x == x => always true
					return ast.Literal{Type: ast.BooleanType, Val: true}, true
				case ast.NeqOp:
					// x != x => always false
					return ast.Literal{Type: ast.BooleanType, Val: false}, true
				case ast.GtOp:
					// x > x => always false
					return ast.Literal{Type: ast.BooleanType, Val: false}, true
				case ast.LtOp:
					// x < x => always false
					return ast.Literal{Type: ast.BooleanType, Val: false}, true
				case ast.GteOp:
					// x >= x => always true
					return ast.Literal{Type: ast.BooleanType, Val: true}, true
				case ast.LteOp:
					// x <= x => always true
					return ast.Literal{Type: ast.BooleanType, Val: true}, true
				}
			}
		}
	}

	// simplify arithmetic identities
	switch node.Operator {
	case ast.AddOp:
		// 0 + x = x
		if left.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 0}) {
			return right, false
		}
		// x + 0 = x
		if right.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 0}) {
			return left, false
		}
		// 1 + x = ++x
		if left.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 1}) {
			return ast.UnaryOp{Operator: ast.IncOp, Arg: right}, false
		}
		// x + 1 = ++x
		if right.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 1}) {
			return ast.UnaryOp{Operator: ast.IncOp, Arg: left}, false
		}
	case ast.SubOp:
		// 0 - x = -x
		if left.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 0}) {
			return ast.UnaryOp{Operator: ast.NegOp, Arg: right}, false
		}
		// x - 0 = x
		if right.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 0}) {
			return left, false
		}
		// x - 1 = --x
		if right.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 1}) {
			return ast.UnaryOp{Operator: ast.DecOp, Arg: left}, false
		}
	case ast.MulOp:
		// 0 * x = 0
		if left.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 0}) {
			return ast.Literal{Type: ast.NumberType, Val: 0}, true
		}
		// x * 0 = 0
		if right.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 0}) {
			return ast.Literal{Type: ast.NumberType, Val: 0}, true
		}
		// 1 * x = x
		if left.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 1}) {
			return right, false
		}
		// x * 1 = x
		if right.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 1}) {
			return left, false
		}
	case ast.DivOp:
		// 0 / x = 0
		if left.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 0}) {
			return ast.Literal{Type: ast.NumberType, Val: 0}, true
		}
		// x / 1 = x
		if right.CompareExpression(ast.Literal{Type: ast.NumberType, Val: 1}) {
			return left, false
		}

		// simplify logical identities
	case ast.AndOp:
		// T && x = x
		if left.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: true}) {
			return right, false
		}
		// x && T = x
		if right.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: true}) {
			return left, false
		}
		// F && x = F
		if left.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: false}) {
			return ast.Literal{Type: ast.BooleanType, Val: false}, true
		}
		// x && F = F
		if right.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: false}) {
			return ast.Literal{Type: ast.BooleanType, Val: false}, true
		}

		// simplify logical identities
	case ast.OrOp:
		// T || x = T
		if left.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: true}) {
			return ast.Literal{Type: ast.BooleanType, Val: true}, true
		}
		// x || T = T
		if right.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: true}) {
			return ast.Literal{Type: ast.BooleanType, Val: true}, true
		}
		// F || x = x
		if left.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: false}) {
			return right, false
		}
		// x || F = x
		if right.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: false}) {
			return left, false
		}
	}

	return ast.BinaryOp{
		Operator: node.Operator,
		LeftArg:  left,
		RightArg: right,
	}, false
}
