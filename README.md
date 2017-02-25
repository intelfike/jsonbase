# 現在大幅な仕様変更中で、README.mdは当てになりません

# update/refer to json like firebase(web).
You do not need to directly manipulate complex nested interface{}<br>
firebase(web)っぽくjsonを加工・参照できるgolangのパッケージです。<br>
あなたが複雑なinterface{}を直接操作する必要はありません。<br>

## install
command

```go get github.com/intelfike/lib/filebase```

## usage

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
fb, _ := filebase.New([]byte(jsonData))

// Add Data
item, _ := fb.Child(0).Clone() // value copy
item.Child("id").Set(6)
item.Child("name").Set("トクガワ")
fb.Fpush(item) // like append()

// Display [class == "A"]
fb.Each(func(f *filebase.Filebase) {
    if f.Child("class").String() == `"A"` {
        fmt.Println(f) // ↓[output]↓
    }
})
```

output

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

## type and func list

### type

```
    type Filebase struct{...}
```

### Maker func

```
    func New(b []byte) (*Filebase, error)
    func NewByFile(name string) (*Filebase, error)
    func NewByReader(reader io.Reader) (*Filebase, error)
```

### Referer func

```
    func (f Filebase) Child(path ...interface{}) *Filebase
    func (f Filebase) Parent() *Filebase
    func (f Filebase) Root() *Filebase
```
Child(...interface{} => string or int) <br>
string => refer map (has not child => return nil/make child) <br>
int => refer array (overflow => panic()/panic()) <br>

### Getter func

```
    func (f Filebase) GetInterface() (*interface{}, error)
    func (f Filebase) String() string
    func (f Filebase) Keys() ([]string, error)
    func (f Filebase) Len() (int, error)
```

GetInterface() => If you want to do type switch then use this.<br>
But do not often use it for eliminate mistake because hard to use.<br>
<br>
String() => You can do type switch with regexp too.<br>
Auto indent by tab.<br>

```
[regexp(string value) => type] 
    ".*" => string 
    [1-9][0-9]* => int 
    [1-9][0-9]*.[0-9]*[1-9] => float 
    (true|false) => bool 
    null => null 
    nil => [NotHasChild] 
```
<br>
Keys() => map keys (not map => Error!) <br>
Len() => array length (not array => Error!) <br>

### Setter func

```
    func (f *Filebase) Fpush(fb *Filebase)
    func (f *Filebase) Fset(fb *Filebase)
    func (f *Filebase) Push(i interface{})
    func (f *Filebase) Set(i interface{}) error
    func (f *Filebase) Remove()
```
Set() => append map or set value<br>
Push() => append array <br>

### Other func
```
    func (f Filebase) Clone() (*Filebase, error) 
    func (f Filebase) Each(fn func(*Filebase))
```
Clone() => value copy. <br>
"f" location become to new json root.<br>
Each() => loop map or array.<br>

## Licence
MIT(適当)