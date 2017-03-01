package filebase

// if has child then
import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
			// _, ok := pathv.(string)
			// if ok{

			// }
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
func (f *Filebase) Remove() error {
	if len(f.path) == 0 {
		return errors.New("Root can't remove. Use Empty().")
	}
	path := f.path[len(f.path)-1]
	i := f.Parent().Interface()
	// switch t := path.(type) {
	// case string:
	// 	m, ok := i.(map[string]interface{})
	// 	if !ok {
	// 		return errors.New("JSON node is not Map.")
	// 	}
	// 	delete(m, t)
	// case int:
	// 	arr := i.([]interface{})
	// 	if 0 > t || t >= len(arr) {
	// 		return errors.New("Array index out of range.")
	// 	}
	// 	f.Parent().Set(append(arr[:t], arr[t+1:]...))
	// }
	switch t := i.(type) {
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
		f.Parent().Set(append(t[:it], t[it+1:]...))
	default:
		return errors.New("JSON node not exists. ")
	}
	return nil
}
func (f *Filebase) Empty() error {
	if f.IsArray() {
		f.Set([]interface{}{})
	} else if f.IsMap() {
		f.Set(map[string]interface{}{})
	} else {
		return errors.New("JSON node is not Array or Map. or node not exists.")
	}
	return nil
}
