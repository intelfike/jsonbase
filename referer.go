package filebase

import (
	"fmt"
	"regexp"
)

// Get json root.
func (f Filebase) Root() *Filebase {
	f.path = []interface{}{}
	return &f
}

// Child(...interface{}.(type) == string or int)
func (f Filebase) Child(path ...interface{}) *Filebase {
	for _, v := range path {
		switch v.(type) {
		case string, int:
			f.path = append(f.path, v)
		default:
			panic("Child(...interface{}.(type) == string or int)")
		}
	}
	return &f
}

var splitReg = regexp.MustCompile("[/\\\\]")

// "a/b/c" => {"a":{"b":{"c":???}}}
//
// "/" Split.
func (f Filebase) ChildPath(path ...string) *Filebase {
	for _, v := range path {
		for _, v2 := range splitReg.Split(v, -1) {
			f.path = append(f.path, v2)
		}
	}
	return &f
}
func (f Filebase) ChildPathf(format string, a ...interface{}) *Filebase {
	return f.ChildPath(fmt.Sprintf(format, a...))
}

// Get json parent.
func (f Filebase) Parent() *Filebase {
	return f.Ancestor(1)
}

// Ancestor(1) == Parent()
func (f Filebase) Ancestor(i int) *Filebase {
	if i < 0 {
		panic("Ancestor() argument can't set under 0")
	}
	anc := len(f.path) - i
	if 0 > anc {
		panic("JSON root has not parent.")
	}
	f.path = f.path[:anc]
	return &f
}
