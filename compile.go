package dialogue

import (
	"compress/zlib"
	"io"
	"io/ioutil"
	"os"

	"github.com/mcvoid/dialogue/internal/codegen"
	"github.com/mcvoid/dialogue/internal/lexer"
	"github.com/mcvoid/dialogue/internal/parser"
	"github.com/mcvoid/dialogue/internal/semantic_analysis"
)

type (
	CompileArg func(*CompileArgs)

	CompileArgs struct {
		codeFolding         bool
		deadCodeElimination bool
		reader              io.Reader
		writer              io.Writer
	}
)

func NoCodeFolding(ca *CompileArgs) {
	ca.codeFolding = false
}

func NoDeadCodeElimination(ca *CompileArgs) {
	ca.deadCodeElimination = false
}

func CompilerInput(r io.Reader) CompileArg {
	return func(ca *CompileArgs) {
		ca.reader = r
	}
}

func CompilerOutput(w io.Writer) CompileArg {
	return func(ca *CompileArgs) {
		ca.writer = w
	}
}

// Compile tranlates a Script from text into its executable form.
// Any reads errors, syntax errors, semantic errors, or write errors
// will result in an error being returned.
func Compile(options ...CompileArg) error {
	args := CompileArgs{
		codeFolding:         true,
		deadCodeElimination: true,
		reader:              os.Stdin,
		writer:              os.Stdout,
	}

	for _, opt := range options {
		opt(&args)
	}

	b, err := ioutil.ReadAll(args.reader)
	if err != nil {
		return err
	}
	l := lexer.New(string(b))
	p := parser.New()
	tree, err := p.Parse(l)
	if err != nil {
		return err
	}

	ast := semantic_analysis.BuildScriptAst(tree)
	if args.codeFolding {
		ast = semantic_analysis.ConstantFoldScript(ast)
	}
	if args.deadCodeElimination {
		ast = semantic_analysis.PruneScript(ast)
	}
	prog, err := codegen.Codegen(ast)
	if err != nil {
		return err
	}

	if args.deadCodeElimination {
		zlibTarget := zlib.NewWriter(args.writer)
		_, err = prog.WriteTo(zlibTarget)
		zlibTarget.Close()
	} else {
		_, err = prog.WriteTo(args.writer)
	}

	return err
}
