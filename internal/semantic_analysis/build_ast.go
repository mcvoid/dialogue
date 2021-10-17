package semantic_analysis

import (
	"encoding/json"

	"github.com/mcvoid/dialogue/internal/types/ast"
	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

func BuildScriptAst(src parsetree.Script) ast.Script {
	dest := ast.Script{
		Nodes:     []ast.Node{},
		Functions: make(map[string][]ast.Type),
	}

	for _, node := range src.Nodes {
		dest.Nodes = append(dest.Nodes, BuildNodeAst(node))
	}

	return dest
}

func BuildNodeAst(src parsetree.Node) ast.Node {
	dest := ast.Node{}

	dest.Name = ast.Symbol(src.Header.Name.Val)
	dest.Body = []ast.BlockElement{}
	for _, block := range src.Blocks {
		dest.Body = append(dest.Body, BuildBlockAst(block))
	}

	return dest
}

func BuildBlockAst(src parsetree.Block) ast.BlockElement {
	switch src := src.(type) {

	case parsetree.Paragraph:
		{
			dest := ast.Paragraph{}
			for _, line := range src.Lines {
				for _, inline := range line.Items {
					switch inline := inline.(type) {
					case parsetree.Text:
						{
							text := inline.Text.Val
							dest = append(dest, ast.Text(text))
						}
					case parsetree.InlineCode:
						{
							expr := BuildExpressionAst(inline.Code)
							dest = append(dest, ast.InlineCode{Expr: expr})
						}
					}
				}
				dest = append(dest, ast.Text("\n"))
			}
			return dest
		}

	case parsetree.CodeBlock:
		{
			dest := ast.CodeBlock{}
			for _, stmt := range src.Code {
				dest.Code = append(dest.Code, BuildStatementAst(stmt))
			}
			return dest
		}

	case parsetree.List:
		{
			dest := ast.Option{}
			for _, listItem := range src.Links {
				dest = append(dest, BuildLinkAst(listItem.Link))
			}
			return dest
		}

	case parsetree.LinkBlock:
		{
			return BuildLinkAst(src.Link)
		}
	}

	return nil
}

func BuildExpressionAst(src parsetree.Expression) ast.Expression {
	switch src := src.(type) {
	case parsetree.NestedExpression:
		return BuildExpressionAst(src.Expr)
	case parsetree.BinaryExpression:
		return BuildBinaryOperationAst(src)
	case parsetree.UnaryExpression:
		return BuildUnaryOperationAst(src)
	case parsetree.Literal:
		return BuildLiteralAst(src)
	}
	return nil
}

func BuildLiteralAst(src parsetree.Literal) ast.Literal {
	switch src.Value.Type {
	case lexeme.Null:
		{
			return ast.Literal{Type: ast.NullType, Val: nil}
		}
	case lexeme.Symbol:
		{
			return ast.Literal{Type: ast.SymbolType, Val: src.Value.Val}
		}
	case lexeme.Number:
		{
			num := 0
			json.Unmarshal([]byte(src.Value.Val), &num)
			return ast.Literal{Type: ast.NumberType, Val: num}
		}
	case lexeme.String:
		{
			str := ""
			json.Unmarshal([]byte(src.Value.Val), &str)
			return ast.Literal{Type: ast.StringType, Val: str}
		}
	case lexeme.Boolean:
		{
			b := false
			json.Unmarshal([]byte(src.Value.Val), &b)
			return ast.Literal{Type: ast.BooleanType, Val: b}
		}
	}
	return ast.Literal{}
}

func BuildUnaryOperationAst(src parsetree.UnaryExpression) ast.UnaryOp {
	dest := ast.UnaryOp{}

	dest.Arg = BuildExpressionAst(src.Operand)
	dest.Operator = map[lexeme.ItemType]ast.UnaryOperator{
		lexeme.Not:   ast.NotOp,
		lexeme.Inc:   ast.IncOp,
		lexeme.Dec:   ast.DecOp,
		lexeme.Minus: ast.NegOp,
	}[src.Operator.Type]

	return dest
}

func BuildBinaryOperationAst(src parsetree.BinaryExpression) ast.BinaryOp {
	dest := ast.BinaryOp{}

	dest.LeftArg = BuildExpressionAst(src.LeftOperand)
	dest.RightArg = BuildExpressionAst(src.RightOperand)
	dest.Operator = map[lexeme.ItemType]ast.BinaryOperator{
		lexeme.Plus:     ast.AddOp,
		lexeme.Minus:    ast.SubOp,
		lexeme.Star:     ast.MulOp,
		lexeme.Slash:    ast.DivOp,
		lexeme.Percent:  ast.ModOp,
		lexeme.DoubleEq: ast.EqOp,
		lexeme.Neq:      ast.NeqOp,
		lexeme.Dot:      ast.ConcatOp,
		lexeme.And:      ast.AndOp,
		lexeme.Or:       ast.OrOp,
		lexeme.Gt:       ast.GtOp,
		lexeme.Gte:      ast.GteOp,
		lexeme.Lt:       ast.LtOp,
		lexeme.Lte:      ast.LteOp,
	}[src.Operator.Type]

	return dest
}

func BuildStatementAst(src parsetree.Statement) ast.Statement {
	switch src := src.(type) {
	case parsetree.Conditional:
		return ast.Conditional{
			Cond:       BuildExpressionAst(src.Cond),
			Consequent: BuildStatementAst(src.Consequent),
			Alternate:  ast.StatementBlock{},
		}
	case parsetree.ConditionalWithElse:
		return ast.Conditional{
			Cond:       BuildExpressionAst(src.Cond),
			Consequent: BuildStatementAst(src.Consequent),
			Alternate:  BuildStatementAst(src.Alternate),
		}
	case parsetree.StatementBlock:
		{
			b := ast.StatementBlock{}
			for _, stmt := range src.Statements {
				b = append(b, BuildStatementAst(stmt))
			}
			return b
		}
	case parsetree.Assignment:
		return ast.Assignment{Name: ast.Symbol(src.Symbol.Val), Val: BuildExpressionAst(src.Value)}
	case parsetree.Loop:
		return ast.Loop{Cond: BuildExpressionAst(src.Cond), Consequent: BuildStatementAst(src.Body)}
	case parsetree.Goto:
		return ast.GotoNode{Name: ast.Symbol(src.Symbol.Val)}
	case parsetree.FunctionCall:
		{
			args := []ast.Expression{}
			for _, expr := range src.Args {
				args = append(args, BuildExpressionAst(expr))
			}
			return ast.FunctionCall{
				Name:   ast.Symbol(src.Symbol.Val),
				Params: args,
			}
		}
	default:
		return nil
	}
}

func BuildLinkAst(src parsetree.Link) ast.Link {
	return ast.Link{
		Dest: ast.Symbol(src.Symbol.Val),
		Text: ast.Text(src.Text.(parsetree.Text).Text.Val),
	}
}
