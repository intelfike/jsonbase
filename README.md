# update/refer to json like firebase(web).
You do not need to directly manipulate complex nested interface{}<br>
firebase(web)っぽくjsonを加工・参照できるgolangのパッケージです。<br>
複雑なinterface{}を直接操作する必要はありません。<br>

## Installation
command

```go get github.com/intelfike/lib/filebase```

## Usage

Append array <br>
↓<br>
Display json node if [class == A]<br>

```
jsonData := `
[
    {"id": 1,"name": "タカハシ","class": "A" },
    {"id": 2,"name": "スズキ","class": "B" },
    {"id": 3,"name": "タナカ","class": "B"},
    {"id": 4,"name": "イシバシ","class": "C"},
    {"id": 5,"name": "ナカヤマ","class": "C"} 
]
`
fb, _ := filebase.New(jsonData)
fb.Indent = "\t"

// Add Data
item, _ := fb.Child(0).Clone() // value copy
item.Child("id").Set(6)
item.Child("name").Set("トクガワ")
fb.PushFB(item) // like append()

// Display [class == "A"]
fb.Each(func(f *filebase.Filebase) {
    if f.Child("class").ToString() == "A" {
        fmt.Println(f) // ↓[output]↓
    }
})
```

### output

```
{
        "class": "A",
        "id": 1,
        "name": "タカハシ"
}
{
        "class": "A",
        "id": 6,
        "name": "トクガワ"
}
```

## How to

### Map append?

```
    fmt.Print(fb.Child("a", "b").Exists()) // [false]
    fb.Child("a", "b").Set(10)
    fmt.Print(fb.Child("a", "b").Exists()) // [true]
```

### Array append?

```
    fb.Push(11)
```

### Make and use template?

```
    template := filebase.MustNew(`{"id":0, "option":"e"}`)
    fb.PushFB(template) // fb.Child(0) == template
    fb.PushFB(template) // fb.Child(1) == template
```

### Get child from filepath or other?

```
    fb.ChildPath("a/b/c").Set(1) // or
    fb.ChildPath("a\\b\\b").Set(2) // or
    fb.ChildPath("a/b\\c").Set(3) // or
    fb.ChildPath("a/b/c", "e/f", "g").Set(4)
```

## type and func list

### type

```
    type Filebase struct{
        Indent string
    }
```

### Maker func

```
    func MustNew(b []byte) *Filebase
    func New(b []byte) (*Filebase, error)
    func NewByFile(filename) (*Filebase, error)
    func NewByReader(reader io.Reader) (*Filebase, error)
```
NewByFile() => gzip file ".gz"

### Writer func

```
    func (f *Filebase) WriteTo(w io.Writer) error
    func (f *Filebase) WriteToFile(filename string) error
```
WriteToFile() => gzip file ".gz"

### Referer func

```
    func (f Filebase) Child(path ...interface{}) *Filebase
    func (f Filebase) ChildPath(path ...string) *Filebase
    func (f Filebase) Parent() *Filebase
    func (f Filebase) Root() *Filebase
```
Child(...interface{} => string or int) <br>
string => refer map (has not child => return nil/make child) <br>
int => refer array (overflow => panic()/panic()) <br>

### Getter func

```
    func (f Filebase) String() string
    func (f *Filebase) Bytes() ([]byte, error)
    func (f *Filebase) BytesIndent() ([]byte, error)

    func (f *Filebase) ToArray() []*Filebase
    func (f *Filebase) ToBool() bool
    func (f *Filebase) ToBytes() []byte
    func (f *Filebase) ToFloat() float64
    func (f *Filebase) ToInt() int64
    func (f *Filebase) ToMap() map[string]*Filebase
    func (f *Filebase) ToString() string
    func (f *Filebase) ToUint() uint64

    func (f *Filebase) Interface() interface{}
    func (f *Filebase) GetInterface() (*interface{}, error)

    func (f *Filebase) Keys() []string
    func (f *Filebase) Len() int
```

To*() => Interface() wrapper<br>
<br>
GetInterface() => If you want to do type switch then use this.<br>
But do not often use it for eliminate mistake because hard to use.<br>
<br>
String() => fmt.Stringer<br>
<br>
Keys() => map keys <br>
Len() => array length <br>

### Setter func

```
    func (f *Filebase) Push(a interface{}) error
    func (f *Filebase) PushFB(fb *Filebase) error
    func (f *Filebase) PushPrint(i ...interface{}) error
    func (f *Filebase) PushPrintf(s string, i ...interface{}) error
    func (f *Filebase) PushStr(s string) error

    func (f *Filebase) Set(i interface{}) error
    func (f *Filebase) SetFB(fb *Filebase) error
    func (f *Filebase) SetPrint(i ...interface{}) error
    func (f *Filebase) SetPrintf(s string, i ...interface{}) error
    func (f *Filebase) SetStr(s string) error
    
    func (f *Filebase) Remove()
```
Set() => append map or set value<br>
Push() => append array <br>

### Checker func

```
    func (f *Filebase) Exists() bool
    func (f *Filebase) HasKey(s string) bool

    func (f *Filebase) IsArray() bool
    func (f *Filebase) IsBool() bool
    func (f *Filebase) IsFloat() bool
    func (f *Filebase) IsInt() bool
    func (f *Filebase) IsMap() bool
    func (f *Filebase) IsNull() bool
    func (f *Filebase) IsString() bool
    func (f *Filebase) IsUint() bool
```

### Other func
```
    func (f *Filebase) Clone() (*Filebase, error) 
    func (f *Filebase) Each(fn func(*Filebase))
```
Clone() => value copy. <br>
"f" location become to new json root.<br>
Each() => loop map or array.<br>

## Licence
MIT