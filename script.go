package dialogue

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"

	"github.com/mcvoid/dialogue/internal/program"
	"github.com/mcvoid/dialogue/internal/vm"
)

type (
	// ExecutionType as a handler's return value is a signal to
	// the Process whether or not to carry on executing once a
	// callback is handled.
	ExecutionType vm.ExecutionType

	// Script is a compiled script. Running this spawns a Process,
	// which will run the dialogue logic.
	Script struct {
		program program.Program
	}

	scriptOptions struct {
		reader       io.Reader
		uncompressed bool
	}

	ScriptOption struct {
		apply func(*scriptOptions)
	}

	// Handler is the Process's interface to the rest of the program.
	// This method is used to handle messages coming from the process.
	Handler interface {
		Handle(m Message) ExecutionType
	}

	HandlerFunc func(m Message) ExecutionType

	MessageType int

	Message struct {
		Type MessageType
		ShowLine
		EnterNode
		ExitNode
		EndScript
		ShowChoice
		FunctionCall
	}

	ShowLine struct {
		Line string
	}

	EnterNode struct {
		NodeEntered string
	}

	ExitNode struct {
		NodeExited string
	}

	EndScript struct {
	}

	ShowChoice struct {
		Options []string
	}

	FunctionCall struct {
		Name string
		Args []interface{}
	}

	// Process is an instance of a script to execute. Run the script by
	// invoking the Start() method.
	Process struct {
		vm *vm.VM
	}
)

const (
	EndScriptType MessageType = iota
	ShowLineType
	EnterNodeType
	ExitNodeType
	ShowChoiceType
	FunctionCallType
)

const (
	// Returning Pause will cause the Process to suspend until its
	// Continue() method is invoked.
	Pause ExecutionType = ExecutionType(vm.PauseExecution)
	// Returning Continue will cause the Process to continue executing
	// and will not stop until either the dialogue ends, an option is
	// shown, or one of the other ScriptHandler methods returns Pause
	Continue ExecutionType = ExecutionType(vm.ContinueExecution)
)

func ScriptInput(r io.Reader) ScriptOption {
	return ScriptOption{
		apply: func(so *scriptOptions) {
			so.reader = r
		},
	}
}

func IsUncompressed() ScriptOption {
	return ScriptOption{
		apply: func(so *scriptOptions) {
			so.uncompressed = true
		},
	}
}

// FromReader loads a compiled script from a file or other reader.
func FromReader(opts ...ScriptOption) (*Script, error) {
	args := scriptOptions{
		reader:       os.Stdin,
		uncompressed: false,
	}
	var err error

	for _, opt := range opts {
		opt.apply(&args)
	}

	p := program.Program{}
	if args.uncompressed {
		_, err = (&p).ReadFrom(args.reader)
		if err != nil {
			return nil, err
		}
	} else {
		rc, err := zlib.NewReader(args.reader)
		if err != nil {
			return nil, err
		}
		_, err = (&p).ReadFrom(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
	}
	return &Script{p}, nil
}

func (h HandlerFunc) Handle(m Message) ExecutionType {
	return h(m)
}

// New spawns a new Process which runs this particular script.
// The Process interacts with the rest of the program by
// invoking various callbacks supplied by the ScriptHandler.
// A new Process is created for each invocation of New.
func (s *Script) New(h Handler) (*Process, error) {
	if h == nil {
		return nil, fmt.Errorf("cannot have nil handler")
	}
	v, _ := vm.New(s.program,
		vm.HandleShowLine(func(v *vm.VM, s string) vm.ExecutionType {
			return vm.ExecutionType(h.Handle(Message{
				Type:     ShowLineType,
				ShowLine: ShowLine{Line: s},
			}))
		}),
		vm.HandleEndDialogue(func(v *vm.VM) {
			h.Handle(Message{
				Type:      EndScriptType,
				EndScript: EndScript{},
			})
		}),
		vm.HandleShowChoice(func(v *vm.VM, s []string) {
			h.Handle(Message{
				Type:       ShowChoiceType,
				ShowChoice: ShowChoice{Options: s},
			})
		}),
		vm.HandleEnterNode(func(v *vm.VM, s string) vm.ExecutionType {
			return vm.ExecutionType(h.Handle(Message{
				Type:      EnterNodeType,
				EnterNode: EnterNode{NodeEntered: s},
			}))
		}),
		vm.HandleExitNode(func(v *vm.VM, s string) vm.ExecutionType {
			return vm.ExecutionType(h.Handle(Message{
				Type:     ExitNodeType,
				ExitNode: ExitNode{NodeExited: s},
			}))
		}),
	)
	return &Process{v}, nil
}

// Start begins execution on a script. Calling this on an in-progress
// Process, even if it's suspended, results in an error.
func (p *Process) Start() error {
	return p.vm.Run()
}

// Resume continues execution of a suspended Process. Calling this on
// a Process which isn't suspended, or is waiting for user input, will
// result in an error.
func (p *Process) Resume() error {
	return p.vm.Resume()
}

// ChooseAndResume continues a Process which is waiting for user input.
// Calling this on a Process which isn't waiting for user input
// will result in an error.
func (p *Process) ChooseAndResume(choice int) error {
	return p.vm.ChooseAndResume(choice)
}
