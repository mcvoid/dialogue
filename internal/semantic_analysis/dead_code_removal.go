package semantic_analysis

import "github.com/mcvoid/dialogue/internal/types/ast"

func PruneScript(script ast.Script) ast.Script {
	prunedScript := ast.Script{
		Functions: map[string][]ast.Type{},
		Nodes:     []ast.Node{},
	}

	for _, node := range script.Nodes {
		prunedScript.Nodes = append(prunedScript.Nodes, PruneNode(node))
	}

	prunedScript.Nodes = PruneUnreachableNodes(prunedScript.Nodes)

	return prunedScript
}

func PruneUnreachableNodes(nodes []ast.Node) []ast.Node {
	if len(nodes) == 0 {
		return nodes
	}

	minScript := []ast.Node{nodes[0]}
	lastLen := 0
	for len(minScript) > lastLen {
		lastLen = len(minScript)
		minScript = growReachableNodes(nodes, minScript)
	}

	return minScript
}

func growReachableNodes(fullProgram, minimalProgram []ast.Node) []ast.Node {
	seenNodes := map[ast.Symbol]bool{}
	for _, node := range minimalProgram {
		seenNodes[node.Name] = true
		for _, block := range node.Body {
			for _, name := range FindNodeNamesInBlock(block) {
				seenNodes[name] = true
			}
		}
	}

	newMin := []ast.Node{}
	for _, node := range fullProgram {
		if seenNodes[node.Name] {
			newMin = append(newMin, node)
		}
	}
	return newMin
}

func FindNodeNamesInBlock(b ast.BlockElement) []ast.Symbol {
	names := []ast.Symbol{}
	switch b := b.(type) {
	case ast.Link:
		names = append(names, b.Dest)
	case ast.Option:
		for _, link := range b {
			names = append(names, link.Dest)
		}
	case ast.CodeBlock:
		for _, stmt := range b.Code {
			names = append(names, FindNodeNamesinStatement(stmt)...)
		}
	}
	return names
}

func FindNodeNamesinStatement(s ast.Statement) []ast.Symbol {
	names := []ast.Symbol{}
	switch s := s.(type) {
	case ast.StatementBlock:
		for _, stmt := range s {
			names = append(names, FindNodeNamesinStatement(stmt)...)
		}
	case ast.GotoNode:
		names = append(names, s.Name)
	case ast.Conditional:
		names = append(names, FindNodeNamesinStatement(s.Consequent)...)
		names = append(names, FindNodeNamesinStatement(s.Alternate)...)
	case ast.Loop:
		names = append(names, FindNodeNamesinStatement(s.Consequent)...)
	}
	return names
}

func PruneNode(node ast.Node) ast.Node {
	prunedBlocks := []ast.BlockElement{}

loop:
	for _, block := range node.Body {
		switch block := block.(type) {
		case ast.Link:
			// anything in a node after a link is unreachable
			prunedBlocks = append(prunedBlocks, block)
			break loop
		case ast.Option:
			// anything in a node after an option is unreachable
			prunedBlocks = append(prunedBlocks, block)
			break loop
		case ast.CodeBlock:
			prunedBlock, endsNode := PruneCodeBlock(block)
			prunedBlocks = append(prunedBlocks, prunedBlock)
			if endsNode {
				// if the code block ends in a goto
				// (and it will if the pruned block contains a top-level goto)
				// then everything else is unreachable
				break loop
			}
		default:
			prunedBlocks = append(prunedBlocks, block)
		}
	}

	return ast.Node{
		Name: node.Name,
		Body: prunedBlocks,
	}
}

func PruneCodeBlock(block ast.CodeBlock) (prunedBlock ast.CodeBlock, endsNode bool) {
	prunedStatements := []ast.Statement{}

	for _, stmt := range block.Code {
		prunedStatement, endsNode := PruneStatement(stmt)
		prunedStatements = append(prunedStatements, prunedStatement)
		if endsNode {
			return ast.CodeBlock{
				Code: prunedStatements,
			}, true
		}
	}

	return ast.CodeBlock{
		Code: prunedStatements,
	}, false
}

func PruneStatement(stmt ast.Statement) (prunedStatement ast.Statement, endsNode bool) {
	switch stmt := stmt.(type) {
	case ast.GotoNode:
		return stmt, true
	case ast.Conditional:
		{
			cons, consEndsNode := PruneStatement(stmt.Consequent)
			alt, altEndsNode := PruneStatement(stmt.Alternate)

			return ast.Conditional{
				Cond:       stmt.Cond,
				Consequent: cons,
				Alternate:  alt,
				// if the node ends no matter which way you take, the rest is unreachable
			}, consEndsNode && altEndsNode
		}
	case ast.Loop:
		{
			if stmt.Cond.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: true}) {
				cons, endsNode := PruneStatement(stmt.Consequent)
				return ast.InfiniteLoop{
					Consequent: cons,
				}, endsNode
			}

			if stmt.Cond.CompareExpression(ast.Literal{Type: ast.BooleanType, Val: false}) {
				return ast.StatementBlock{}, false
			}

			cons, endsNode := PruneStatement(stmt.Consequent)
			return ast.Loop{
				Cond:       stmt.Cond,
				Consequent: cons,
			}, endsNode

		}
	case ast.StatementBlock:
		{
			stmts := ast.StatementBlock{}

			for _, s := range stmt {
				prunedStatement, endsNode := PruneStatement(s)
				stmts = append(stmts, prunedStatement)
				if endsNode {
					// stop adding statements here - the rest is unreachable
					return stmts, true
				}
			}
			return stmts, false
		}
	default:
		return stmt, false

	}
}
