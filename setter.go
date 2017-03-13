package jsonbase

// if has child then
import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// update child.
//
// else
// make child.
//
// Can't make array.
func (f *Jsonbase) Set(i interface{}) error {
	if len(f.path) == 0 {
		*f.master = i
		return nil
	} else {
		// topがマップや配列じゃなかったら作る
		switch (*f.master).(type) {
		case []interface{}, map[string]interface{}:
		default:
			*f.master = map[string]interface{}{}
		}
	}
	cur := new(interface{})
	*cur = *f.master
	for n, pathv := range f.path {
		switch mas := (*cur).(type) {
		case map[string]interface{}:
			pt, ok := pathv.(string)
			if !ok {
				pt = strconv.Itoa(pathv.(int))
			}
			*cur, ok = mas[pt]
			if !ok {
				paths := []string{}
				for _, v := range f.path[n:] {
					s, ok := v.(string)
					if !ok {
						s = strconv.Itoa(v.(int))
					}
					paths = append(paths, s)
				}
				mapNest(mas, i, 0, paths...)
				return nil
			}
			if n == len(f.path)-1 {
				mas[pt] = i
			}
		case []interface{}:
			switch pt := pathv.(type) {
			case string:
				f.Parent().Set(map[string]interface{}{})
				f.Set(i)
			case int:
				if 0 > pt || pt >= len(mas) {
					return errors.New("Array index out of range.")
				}
				if n == len(f.path)-1 { // 最後の要素なら
					mas[pt] = i
					return nil
				}
				*cur = mas[pt]
			}
		default:

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

// Set *Jsonbase.
func (f *Jsonbase) SetFB(fb *Jsonbase) error {
	b, err := fb.Bytes()
	if err != nil {
		return err
	}
	return f.SetReader(bytes.NewReader(b))
}

func (f *Jsonbase) SetReader(r io.Reader) error {
	i := new(interface{})
	dec := json.NewDecoder(r)
	err := dec.Decode(i)
	if err != nil {
		return err
	}
	return f.Set(*i)
}

func (f *Jsonbase) SetReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if strings.HasSuffix(filename, ".gz") {
		r, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer r.Close()
		return f.SetReader(r)
	}

	return f.SetReader(file)
}

// Set string.
func (f *Jsonbase) SetStr(s string) error {
	return f.SetReader(strings.NewReader(s))
}

// SetStr() + fmt.Sprint().
func (f *Jsonbase) SetPrint(a ...interface{}) error {
	return f.SetStr(fmt.Sprint(a...))
}

// SetStr() + fmt.Sprintf().
func (f *Jsonbase) SetPrintf(format string, a ...interface{}) error {
	return f.SetStr(fmt.Sprintf(format, a...))
}

// It like append().
//
// If json node is array then add i.
//
// else then set []interface{i} (initialize for array).
func (f *Jsonbase) Push(a interface{}) error {
	if !f.IsArray() {
		f.Set([]interface{}{a})
		return nil
	}
	pt, err := f.GetInterfacePt()
	if err != nil {
		return err
	}
	ar := (*pt).([]interface{})
	return f.Set(append(ar, a))
}

// Push *Jsonbase.
func (f *Jsonbase) PushFB(fb *Jsonbase) error {
	b, err := fb.Bytes()
	if err != nil {
		return err
	}
	return f.PushReader(bytes.NewReader(b))
}

func (f *Jsonbase) PushReader(r io.Reader) error {
	i := new(interface{})
	dec := json.NewDecoder(r)
	err := dec.Decode(i)
	if err != nil {
		return err
	}
	return f.Push(*i)
}

func (f *Jsonbase) PushReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if strings.HasSuffix(filename, ".gz") {
		r, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer r.Close()
		return f.PushReader(r)
	}

	return f.PushReader(file)
}

// push string.
func (f *Jsonbase) PushStr(s string) error {
	return f.PushReader(strings.NewReader(s))
}

// PushStr() + fmt.Sprint().
func (f *Jsonbase) PushPrint(a ...interface{}) error {
	return f.PushStr(fmt.Sprint(a...))
}

// PushStr() + fmt.Sprintf().
func (f *Jsonbase) PushPrintf(format string, a ...interface{}) error {
	return f.PushStr(fmt.Sprintf(format, a...))
}

// Remove() remove map or array element
func (f *Jsonbase) Remove() error {
	if len(f.path) == 0 {
		return errors.New("Root can't remove. Use Empty().")
	}
	path := f.path[len(f.path)-1]
	pt, err := f.Parent().GetInterfacePt()
	if err != nil {
		return err
	}
	switch t := (*pt).(type) {
	case map[string]interface{}:
		st, ok := path.(string)
		if !ok {
			st = strconv.Itoa(path.(int))
		}
		_, ok = t[st]
		if !ok {
			return errors.New("JSON node not exists.")
		}
		delete(t, st)
	case []interface{}:
		it, ok := path.(int)
		if !ok {
			return errors.New("JSON node is not Map().")
		}
		if 0 > it || it >= len(t) {
			return errors.New("Array index out of range.")
		}
		return f.Parent().Set(append(t[:it], t[it+1:]...))
	default:
		return errors.New("JSON node not exists. ")
	}
	return nil
}
func (f *Jsonbase) Empty() error {
	if f.IsArray() {
		f.Set([]interface{}{})
	} else if f.IsMap() {
		f.Set(map[string]interface{}{})
	} else {
		return errors.New("JSON node is not Array or Map. or node not exists.")
	}
	return nil
}
