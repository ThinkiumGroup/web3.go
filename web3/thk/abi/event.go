package abi

import (
	"fmt"
	"strings"
	"web3.go/common"
)

type Event struct {
	Name      string
	Anonymous bool
	Inputs    Arguments
}

func (e Event) String() string {
	inputs := make([]string, len(e.Inputs))
	for i, input := range e.Inputs {
		inputs[i] = fmt.Sprintf("%v %v", input.Type, input.Name)
		if input.Indexed {
			inputs[i] = fmt.Sprintf("%v indexed %v", input.Type, input.Name)
		}
	}
	return fmt.Sprintf("event %v(%v)", e.Name, strings.Join(inputs, ", "))
}

func (e Event) id() []byte {
	types := make([]string, len(e.Inputs))
	i := 0
	for _, input := range e.Inputs {
		types[i] = input.Type.String()
		i++
	}
	methodSignature := fmt.Sprintf("%v(%v)", e.Name, strings.Join(types, ","))
	return common.Hash256(methodSignature)
}

func (e Event) Id() common.Hash {
	return common.BytesToHash(e.id())
}
