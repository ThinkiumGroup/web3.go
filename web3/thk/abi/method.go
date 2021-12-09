package abi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ThinkiumGroup/web3.go/common"
	"reflect"
	"strings"
)

const (
	ErrorMethodId = "0x08c379a0"
)

var ErrorMethod Method

func init() {
	typeString, _ := NewType("string", []ArgumentMarshaling{{
		Name:       "string",
		Type:       "string",
		Components: nil,
		Indexed:    false,
	}})
	ErrorMethod = Method{
		Name:   "error",
		Const:  false,
		Inputs: nil,
		Outputs: []Argument{{
			Name:    "revertReason",
			Type:    typeString,
			Indexed: false,
		}},
	}
}

type Method struct {
	Name    string
	Const   bool
	Inputs  Arguments
	Outputs Arguments
}

func (method Method) Sig() string {
	types := make([]string, len(method.Inputs))
	for i, input := range method.Inputs {
		types[i] = input.Type.String()
	}
	return fmt.Sprintf("%v(%v)", method.Name, strings.Join(types, ","))
}

func (method Method) String() string {
	inputs := make([]string, len(method.Inputs))
	for i, input := range method.Inputs {
		inputs[i] = fmt.Sprintf("%v %v", input.Type, input.Name)
	}
	outputs := make([]string, len(method.Outputs))
	for i, output := range method.Outputs {
		outputs[i] = output.Type.String()
		if len(output.Name) > 0 {
			outputs[i] += fmt.Sprintf(" %v", output.Name)
		}
	}
	constant := ""
	if method.Const {
		constant = "constant "
	}
	return fmt.Sprintf("function %v(%v) %sreturns(%v)", method.Name, strings.Join(inputs, ", "), constant, strings.Join(outputs, ", "))
}

func (method Method) Id() []byte {
	return common.SystemHash256([]byte(method.Sig()))[:4]
}

func (method Method) singleInputUnpack(v interface{}, input []byte) error {
	valueOf := reflect.ValueOf(v)
	if reflect.Ptr != valueOf.Kind() {
		s := fmt.Sprintf("abi: Unpack(non-pointer %T)", v)
		return errors.New(s)
	}

	value := valueOf.Elem()
	marshalledValue, err := toGoType(0, method.Inputs[0].Type, input)
	if err != nil {
		return err
	}

	if err := myset(value, reflect.ValueOf(marshalledValue), method.Inputs[0]); err != nil {
		return err
	}

	return nil
}

func (method Method) multInputUnpack(v []interface{}, input []byte) error {
	j := 0
	for i := 0; i < len(method.Inputs); i++ {
		valueOf := reflect.ValueOf(v[i])
		if reflect.Ptr != valueOf.Kind() {
			s := fmt.Sprintf("abi: Unpack(non-pointer %T)", v)
			return errors.New(s)
		}

		toUnpack := method.Inputs[i]
		if toUnpack.Type.T == ArrayTy {
			j += toUnpack.Type.Size
		}

		marshalledValue, err := toGoType((i+j)*32, toUnpack.Type, input)
		if err != nil {
			return err
		}

		if err := myset(valueOf.Elem(), reflect.ValueOf(marshalledValue), method.Inputs[i]); err != nil {
			return err
		}
	}
	return nil
}

func IsErrorOutput(output string) bool {
	return strings.HasPrefix(output, ErrorMethodId)
}

type ErrorOutput struct {
	RevertReason string `json:"revertReason"`
}

func ExtractRevertReason(output string) (*ErrorOutput, error) {
	if !IsErrorOutput(output) {
		return nil, fmt.Errorf("not error output")
	}

	output = common.CleanHexPrefix(output)

	decodeBytes, err := hex.DecodeString(output)
	if err != nil {
		return nil, err
	}

	res := new(ErrorOutput)
	if err := ErrorMethod.Outputs.Unpack(res, decodeBytes[4:]); err != nil {
		return nil, err
	}
	return res, nil
}
