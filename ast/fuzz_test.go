package ast

import "testing"

var crasher = `package _ s={[[]|g;0]:{}}`

func TestFuzzX(t *testing.T) {
	_, _, err := ParseStatements("", crasher)

	if err == nil {
		CompileModules(map[string]string{"": crasher})
	}
}
