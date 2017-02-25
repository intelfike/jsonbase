// Gopherのためのjson操作パッケージ
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

// file name to *Filebase.
func NewByFile(filename, indent string) (*Filebase, error) {
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
		return NewByReader(r, indent)
	}

	return NewByReader(file, indent)
}

func NewByReader(reader io.Reader, indent string) (*Filebase, error) {
	fb := new(Filebase)
	fb.Indent = indent
	fb.master = new(interface{})
	r := json.NewDecoder(reader)
	err := r.Decode(fb.master)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// Byte data to *Filebase
func New(b []byte, indent string) (*Filebase, error) {
	fb := new(Filebase)
	fb.Indent = indent
	fb.master = new(interface{})
	err := json.Unmarshal(b, fb.master)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

func MustNew(b []byte, indent string) *Filebase {
	fb, err := New(b, indent)
	if err != nil {
		panic(err)
	}
	return fb
}

func (f *Filebase) WriteTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(f.Interface())
}

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

// if has child then
// update child.
//
// else
// make child.
//
// Can't make array.
func (f *Filebase) Set(i interface{}) error {
	if len(f.path) == 0 {
		*f.master = i
		return nil
	}
	cur := new(interface{})
	*cur = *f.master
	prev := *cur
	prevkey := ""
	_ = prevkey
	for n, pathv := range f.path {
		switch pt := pathv.(type) {
		case string:
			switch mas := (*cur).(type) {
			case map[string]interface{}:
				var ok bool
				prev = *cur
				*cur, ok = mas[pt]
				if !ok {
					mas[pt] = nil //map[string]interface{}{pt: i}
				}
				if n == len(f.path)-1 {
					mas[pt] = i
				}
			default:
				paths := []string{prevkey}
				for _, v := range f.path[n:] {
					paths = append(paths, v.(string))
				}
				m, ok := prev.(map[string]interface{})
				if !ok {

					panic("can't map make")
				}
				mapNest(m, i, 0, paths...)
				return nil
			}
			prevkey = pt
		case int:
			mas, ok := (*cur).([]interface{})
			if !ok {
				return errors.New("JSON node is not array.")
			}
			if n == len(f.path)-1 { // 最後の要素なら
				mas[pt] = i
				return nil
			}
			*cur = mas[pt]
		}
	}
	return nil
}
func mapNest(m map[string]interface{}, val interface{}, depth int, s ...string) {
	if depth == len(s)-1 {
		m[s[depth]] = val
		return
	}
	mm := map[string]interface{}{s[depth+1]: nil}
	m[s[depth]] = mm
	mapNest(mm, val, depth+1, s...)
}

// It like append().
//
// If json node is array then add i.
//
// else then set []interface{i} (initialize for array).
func (f *Filebase) Push(a interface{}) {
	if !f.IsArray() {
		f.Set([]interface{}{a})
		return
	}
	i := f.Interface()
	ar := i.([]interface{})
	f.Set(append(ar, a))
}

// Set() => Filebase to Filebase.
func (f *Filebase) Fset(fb *Filebase) {
	i := fb.Interface()
	f.Set(i)
}

// Push() => Filebase to Filebase.
func (f *Filebase) Fpush(fb *Filebase) {
	i := fb.Interface()
	f.Push(i)
}

// Remove() remove map or array element
func (f *Filebase) Remove() {
	if len(f.path) == 0 {
		if f.IsArray() {
			f.Set([]interface{}{})
		} else if f.IsMap() {
			f.Set(map[string]interface{}{})
		}
		return
	}
	path := f.path[len(f.path)-1]
	i := f.Parent().Interface()
	switch t := path.(type) {
	case string:
		delete(i.(map[string]interface{}), t)
	case int:
		arr := i.([]interface{})
		f.Parent().Set(append(arr[:t], arr[t+1:]...))
	}
}

// If you want to do type switch then use this.
// Do not use it much.
//
// You can do type switch with regexp too.
// Refer to String().
//
// This function get interface{} pinter.
func (f *Filebase) GetInterfacePt() (*interface{}, error) {
	cur := new(interface{})
	*cur = *f.master
	for _, pathv := range f.path {
		switch pt := pathv.(type) {
		case string:
			mas, ok := (*cur).(map[string]interface{})
			if !ok {
				return nil, errors.New("JSON node is not map")
			}
			*cur, ok = mas[pt]
			if !ok {
				return nil, errors.New("JSON node not found.")
			}
		case int:
			mas, ok := (*cur).([]interface{})
			if !ok {
				return nil, errors.New("JSON node is not array")
			}
			*cur = mas[pt]
		}
	}
	return cur, nil
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
	json.Unmarshal(f.Bytes(), newfb.master)
	return newfb, nil
}

// If json node is map then return key list & nil.
//
// else then return nil & error.
func (f *Filebase) Keys() []string {
	i := f.Interface()
	s := []string{}
	for key, _ := range i.(map[string]interface{}) {
		s = append(s, key)
	}
	return s
}

//This get len, check if array.
//
// If json node is array then return len(array) & nil.
//
// else then return -1 & error.
func (f *Filebase) Len() int {
	i := f.Interface()
	return len(i.([]interface{}))
}
