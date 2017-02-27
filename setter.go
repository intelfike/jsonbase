package filebase

// if has child then
import (
	"encoding/json"
	"errors"
	"fmt"
)

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
					s, ok := v.(string)
					if !ok {
						return errors.New("Child(...) can't make or append array. You can make or append array with Push() or Fpush().")
					}
					paths = append(paths, s)
				}
				m, ok := prev.(map[string]interface{})
				if !ok {
					return errors.New("JSON node is not map. IsMap() == true?")
				}
				mapNest(m, i, 0, paths...)
				return nil
			}
			prevkey = pt
		case int:
			mas, ok := (*cur).([]interface{})
			if !ok {
				return errors.New("JSON node is not array. IsArray() == true?")
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

// Set *Filebase.
func (f *Filebase) SetFB(fb *Filebase) error {
	b, err := fb.Bytes()
	if err != nil {
		return err
	}
	return f.SetStr(string(b))
}

// Set string.
func (f *Filebase) SetStr(s string) error {
	i := new(interface{})
	err := json.Unmarshal([]byte(s), i)
	if err != nil {
		return err
	}
	return f.Set(*i)
}

// SetStr() + fmt.Sprint().
func (f *Filebase) SetPrint(a ...interface{}) error {
	return f.SetStr(fmt.Sprint(a...))
}

// SetStr() + fmt.Sprintf().
func (f *Filebase) SetPrintf(format string, a ...interface{}) error {
	return f.SetStr(fmt.Sprintf(format, a...))
}

// It like append().
//
// If json node is array then add i.
//
// else then set []interface{i} (initialize for array).
func (f *Filebase) Push(a interface{}) error {
	if !f.IsArray() {
		f.Set([]interface{}{a})
		return nil
	}
	i := f.Interface()
	ar := i.([]interface{})
	return f.Set(append(ar, a))
}

// Push *Filebase.
func (f *Filebase) PushFB(fb *Filebase) error {
	b, err := fb.Bytes()
	if err != nil {
		return err
	}
	return f.PushStr(string(b))
}

// push string.
func (f *Filebase) PushStr(s string) error {
	i := new(interface{})
	err := json.Unmarshal([]byte(s), i)
	if err != nil {
		return err
	}
	return f.Push(*i)
}

// PushStr() + fmt.Sprint().
func (f *Filebase) PushPrint(a ...interface{}) error {
	return f.PushStr(fmt.Sprint(a...))
}

// PushStr() + fmt.Sprintf().
func (f *Filebase) PushPrintf(format string, a ...interface{}) error {
	return f.PushStr(fmt.Sprintf(format, a...))
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
