package program

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/mcvoid/dialogue/internal/types/asm"
)

type Program struct {
	Start int                   `json:"start"`
	Code  []asm.Instruction     `json:"code"`
	Funcs map[string][]asm.Type `json:"funcs,omitempty"`
}

func (p *Program) ReadFrom(r io.Reader) (n int64, err error) {
	b, err := ioutil.ReadAll(r)
	bytesRead := int64(len(b))
	if err != nil {
		return bytesRead, err
	}
	err = json.Unmarshal(b, p)
	if p.Funcs == nil {
		p.Funcs = map[string][]asm.Type{}
	}
	return bytesRead, err
}

func (p *Program) WriteTo(w io.Writer) (n int64, err error) {
	b, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}

	num, err := w.Write(b)

	return int64(num), err
}
