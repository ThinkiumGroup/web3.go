package compiler

import (
	"strings"
)

func CompileContract(file ...string) (map[string]*Contract, error) {
	resMap, err := CompileSolidity("", file...)
	if err != nil {
		return nil, err
	}
	for key, v := range resMap {
		if strings.Contains(key, "<stdin>:") {
			newKey := strings.Replace(key, "<stdin>:", "", -1)
			delete(resMap, key)
			resMap[newKey] = v
		}
	}
	return resMap, err
}
