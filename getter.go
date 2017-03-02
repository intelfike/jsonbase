package filebase

// null & nil?
import (
	"encoding/json"
	"errors"
	"strconv"
)

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
		switch mas := (*cur).(type) {
		case map[string]interface{}:
			spt, ok := pathv.(string)
			if !ok {
				spt = strconv.Itoa(pathv.(int))
			}
			*cur, ok = mas[spt]
			if !ok {
				return nil, errors.New("JSON node not exists.")
			}
		case []interface{}:
			switch pt := pathv.(type) {
			case int:
				if 0 > pt || pt >= len(mas) {
					return nil, errors.New("Array index out of range.")
				}
				*cur = mas[pt]
			default:
				return nil, errors.New("JSON node is not Map.")
			}
		default:
			return nil, errors.New("JSON node not found.")
		}
	}
	return cur, nil
}

//  fmt.Stringer interface
func (f Filebase) String() string {
	var b []byte
	if len(f.Indent) == 0 {
		b, _ = f.Bytes()
	} else {
		b, _ = f.BytesIndent()
	}
	return string(b)
}

func (f *Filebase) Bytes() ([]byte, error) {
	return json.Marshal(f.Interface())
}
func (f *Filebase) BytesIndent() ([]byte, error) {
	return json.MarshalIndent(f.Interface(), "", f.Indent)
}

// Assert string.
func (f *Filebase) ToString() string {
	return f.Interface().(string)
}

func (f *Filebase) ToBytes() []byte {
	return []byte(f.Interface().(string))
}

func (f *Filebase) ToBool() bool {
	return f.Interface().(bool)
}

func (f *Filebase) ToInt() int64 {
	return f.Interface().(int64)
}

func (f *Filebase) ToUint() uint64 {
	return f.Interface().(uint64)
}

func (f *Filebase) ToFloat() float64 {
	return f.Interface().(float64)
}

func (f *Filebase) ToArray() []*Filebase {
	arr := f.Interface().([]interface{})
	rv := make([]*Filebase, len(arr))
	for n, _ := range arr {
		rv[n] = f.Child(n)
	}
	return rv
}

func (f *Filebase) ToMap() map[string]*Filebase {
	m := map[string]*Filebase{}
	for k, _ := range f.Interface().(map[string]interface{}) {
		m[k] = f.Child(k)
	}
	return m
}

func (f *Filebase) Interface() interface{} {
	i, err := f.GetInterfacePt()
	if err != nil {
		return nil
	}
	return *i
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
	ar := f.Interface().([]interface{})
	return len(ar)
}

func (f *Filebase) Path() []interface{} {
	return f.path
}
func (f *Filebase) BottomPath() interface{} {
	if len(f.path) == 0 {
		return nil
	}
	return f.path[len(f.path)-1]
}
