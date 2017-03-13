# filebase ver1.2

update/refer to json like firebase(web).
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
fb.Each(func(f *filebase.Jsonbase) {
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
    type Jsonbase struct{
        Indent string
    }
```


[checker.go]
```
func (f *Jsonbase) Exists() bool {
func (f *Jsonbase) ReferError() error {
func (f *Jsonbase) IsString() bool {
func (f *Jsonbase) IsBool() bool {
func (f *Jsonbase) IsInt() bool {
func (f *Jsonbase) IsUint() bool {
func (f *Jsonbase) IsFloat() bool {
func (f *Jsonbase) IsNull() bool {
func (f *Jsonbase) IsArray() bool {
func (f *Jsonbase) IsMap() bool {
func (f *Jsonbase) HasChild(a interface{}) bool {
```

[getter.go]
```
func (f *Jsonbase) GetInterfacePt() (*interface{}, error) {
func New() *Jsonbase {
func (f *Jsonbase) WriteTo(w io.Writer) error {
func (f *Jsonbase) WriteToFile(filename string) error {
func (f Jsonbase) String() string {
func (f *Jsonbase) Bytes() ([]byte, error) {
func (f *Jsonbase) BytesIndent() ([]byte, error) {
func (f *Jsonbase) ToString() string {
func (f *Jsonbase) ToBytes() []byte {
func (f *Jsonbase) ToBool() bool {
func (f *Jsonbase) ToInt() int64 {
func (f *Jsonbase) ToUint() uint64 {
func (f *Jsonbase) ToFloat() float64 {
func (f *Jsonbase) ToArray() []*Jsonbase {
func (f *Jsonbase) ToMap() map[string]*Jsonbase {
func (f *Jsonbase) Interface() interface{} {
func (f *Jsonbase) Keys() []string {
func (f *Jsonbase) Len() int {
func (f *Jsonbase) Path() []interface{} {
func (f *Jsonbase) BottomPath() interface{} {
```

[golebase.go]
```
func (f *Jsonbase) Each(fn func(*Jsonbase)) {
func (f *Jsonbase) Clone() (*Jsonbase, error) {
```

[referer.go]
```
func (f Jsonbase) Root() *Jsonbase {
func (f Jsonbase) Child(path ...interface{}) *Jsonbase {
func (f Jsonbase) ChildPath(path ...string) *Jsonbase {
func (f Jsonbase) ChildPathf(format string, a ...interface{}) *Jsonbase {
func (f Jsonbase) Parent() *Jsonbase {
func (f Jsonbase) Ancestor(i int) *Jsonbase {
```

[setter.go]
```
func (f *Jsonbase) Set(i interface{}) error {
func mapNest(m map[string]interface{}, val interface{}, depth int, s ...string) {
func (f *Jsonbase) SetFB(fb *Jsonbase) error {
func (f *Jsonbase) SetReader(r io.Reader) error {
func (f *Jsonbase) SetReadFile(filename string) error {
func (f *Jsonbase) SetStr(s string) error {
func (f *Jsonbase) SetPrint(a ...interface{}) error {
func (f *Jsonbase) SetPrintf(format string, a ...interface{}) error {
func (f *Jsonbase) Push(a interface{}) error {
func (f *Jsonbase) PushFB(fb *Jsonbase) error {
func (f *Jsonbase) PushReader(r io.Reader) error {
func (f *Jsonbase) PushReadFile(filename string) error {
func (f *Jsonbase) PushStr(s string) error {
func (f *Jsonbase) PushPrint(a ...interface{}) error {
func (f *Jsonbase) PushPrintf(format string, a ...interface{}) error {
func (f *Jsonbase) Remove() error {
func (f *Jsonbase) Empty() error {
```

## Licence
MIT