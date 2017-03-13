package jsonbase

// null & nil?
import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// If you want to do type switch then use this.
// Do not use it much.
//
// You can do type switch with regexp too.
// Refer to String().
//
// This function get interface{} pinter.
func (f *Jsonbase) GetInterfacePt() (*interface{}, error) {
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
			case string:
				i, err := strconv.Atoi(pt)
				if err != nil {
					return nil, errors.New("JSON node is Array. But you tried to refer by string key.")
				}
				if 0 > i || i >= len(mas) {
					return nil, errors.New("Array index out of range.")
				}
				*cur = mas[i]
			}
		default:
			return nil, errors.New("JSON node not found.")
		}
	}
	return cur, nil
}

func New() *Jsonbase {
	f := new(Jsonbase)
	f.master = new(interface{})
	return f
}

func (f *Jsonbase) WriteTo(w io.Writer) error {
	pt, err := f.GetInterfacePt()
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(*pt)
}

// gzip file ".gz".
func (f *Jsonbase) WriteToFile(filename string) error {
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

//  fmt.Stringer interface
func (f Jsonbase) String() string {
	var b []byte
	if len(f.Indent) == 0 {
		b, _ = f.Bytes()
	} else {
		b, _ = f.BytesIndent()
	}
	return string(b)
}

func (f *Jsonbase) Bytes() ([]byte, error) {
	return json.Marshal(f.Interface())
}
func (f *Jsonbase) BytesIndent() ([]byte, error) {
	return json.MarshalIndent(f.Interface(), "", f.Indent)
}

// Assert string.
func (f *Jsonbase) ToString() string {
	return f.Interface().(string)
}

func (f *Jsonbase) ToBytes() []byte {
	return []byte(f.Interface().(string))
}

func (f *Jsonbase) ToBool() bool {
	return f.Interface().(bool)
}

func (f *Jsonbase) ToInt() int64 {
	return f.Interface().(int64)
}

func (f *Jsonbase) ToUint() uint64 {
	return f.Interface().(uint64)
}

func (f *Jsonbase) ToFloat() float64 {
	return f.Interface().(float64)
}

func (f *Jsonbase) ToArray() []*Jsonbase {
	arr := f.Interface().([]interface{})
	rv := make([]*Jsonbase, len(arr))
	for n, _ := range arr {
		rv[n] = f.Child(n)
	}
	return rv
}

func (f *Jsonbase) ToMap() map[string]*Jsonbase {
	m := map[string]*Jsonbase{}
	for k, _ := range f.Interface().(map[string]interface{}) {
		m[k] = f.Child(k)
	}
	return m
}

func (f *Jsonbase) Interface() interface{} {
	i, err := f.GetInterfacePt()
	if err != nil {
		return nil
	}
	return *i
}

// If json node is map then return key list & nil.
//
// else then return nil & error.
func (f *Jsonbase) Keys() []string {
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
func (f *Jsonbase) Len() int {
	ar := f.Interface().([]interface{})
	return len(ar)
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
