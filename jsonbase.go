// Gopherのためのjson操作パッケージ var1.1
package jsonbase

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Jsonbase struct {
	master *interface{}
	path   []interface{}
	Indent string
}

var _ fmt.Stringer = Jsonbase{}

func New() *Jsonbase {
	f := new(Jsonbase)
	f.master = new(interface{})
	return f
}

// loop map or array
func (f *Jsonbase) Each(fn func(*Jsonbase)) {
	if f.IsArray() {
		for n := 0; n < f.Len(); n++ {
			fn(f.Child(n))
		}
	}
	if f.IsMap() {
		for _, key := range f.Keys() {
			fn(f.Child(key))
		}
	}
}

// f location become to new json root
func (f *Jsonbase) Clone() (*Jsonbase, error) {
	if !f.Exists() {
		return nil, errors.New("JSON node not exists.")
	}
	newfb := new(Jsonbase)
	newfb.Indent = f.Indent
	newfb.master = new(interface{})
	b, err := f.Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &newfb.master)
	return newfb, err
}
