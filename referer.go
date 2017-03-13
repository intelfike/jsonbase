package jsonbase

import (
	"fmt"
	"regexp"
)

// Get json root.
func (f Jsonbase) Root() *Jsonbase {
	f.path = []interface{}{}
	return &f
}

// Child(...interface{}.(type) == string or int)
func (f Jsonbase) Child(path ...interface{}) *Jsonbase {
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
func (f Jsonbase) ChildPath(path ...string) *Jsonbase {
	for _, v := range path {
		for _, v2 := range splitReg.Split(v, -1) {
			f.path = append(f.path, v2)
		}
	}
	return &f
}
func (f Jsonbase) ChildPathf(format string, a ...interface{}) *Jsonbase {
	return f.ChildPath(fmt.Sprintf(format, a...))
}

// Get json parent.
func (f Jsonbase) Parent() *Jsonbase {
	return f.Ancestor(1)
}

// Ancestor(1) == Parent()
func (f Jsonbase) Ancestor(i int) *Jsonbase {
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

func (f *Jsonbase) Path() []interface{} {
	return f.path
}
func (f *Jsonbase) BottomPath() interface{} {
	if len(f.path) == 0 {
		return nil
	}
	return f.path[len(f.path)-1]
}

func (f *Jsonbase) Depth() int {
	return len(f.path)
}
