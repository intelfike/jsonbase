// Gopherのためのjson操作パッケージ var1.1
package filebase

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Filebase struct {
	master *interface{}
	path   []interface{}
	Indent string
}

var _ fmt.Stringer = Filebase{}

var (
	Array = NewByUnmarshaled([]interface{}{})
	Map   = NewByUnmarshaled(map[string]interface{}{})
)

// file name to *Filebase.
// gzip file ".gz".
func NewByFile(filename string) (*Filebase, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if strings.HasSuffix(filename, ".gz") {
		r, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		return NewByReader(r)
	}

	return NewByReader(file)
}

func NewByReader(reader io.Reader) (*Filebase, error) {
	fb := new(Filebase)
	fb.master = new(interface{})
	r := json.NewDecoder(reader)
	err := r.Decode(fb.master)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// Byte data to *Filebase
func New(s string) (*Filebase, error) {
	fb := new(Filebase)
	fb.master = new(interface{})
	err := json.Unmarshal([]byte(s), fb.master)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

func MustNew(s string) *Filebase {
	fb, err := New(s)
	if err != nil {
		panic(err)
	}
	return fb
}

func NewByUnmarshaled(i interface{}) *Filebase {
	fb := new(Filebase)
	fb.master = &i
	return fb
}

func (f *Filebase) WriteTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(f.Interface())
}

// gzip file ".gz".
func (f *Filebase) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if strings.HasSuffix(filename, ".gz") {
		w := gzip.NewWriter(file)
		defer w.Close()
		return f.WriteTo(w)
	}
	return f.WriteTo(file)
}

// loop map or array
func (f *Filebase) Each(fn func(*Filebase)) {
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
func (f *Filebase) Clone() (*Filebase, error) {
	if !f.Exists() {
		return nil, errors.New("JSON node not exists.")
	}
	newfb := new(Filebase)
	newfb.Indent = f.Indent
	newfb.master = new(interface{})
	b, err := f.Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, newfb.master)
	return newfb, err
}
