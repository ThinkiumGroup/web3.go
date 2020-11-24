package compiler

import (
	"encoding/json"
	"strings"
)

func CompileContract(name ...string) (res map[string]interface{}, err error) {
	data := new(map[string]interface{})
	ress, err := CompileSolidity("", name...)
	if err != nil {
		return nil, err
	}
	for key, v := range ress {
		if strings.Contains(key, "<stdin>:") {
			newKey := strings.Replace(key, "<stdin>:", "", -1)
			delete(ress, key)
			ress[newKey] = v
		}
	}
	solcres, err := json.Marshal(ress)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(solcres, data); err != nil {
		return nil, err
	}
	return *data, err
}
