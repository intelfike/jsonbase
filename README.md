# jsonbase ver1.0

update/refer to json like firebase(web).
You do not need to directly manipulate complex nested interface{}<br>
firebase(web)っぽくjsonを加工・参照できるgolangのパッケージです。<br>
複雑なinterface{}を直接操作する必要はありません。<br>

## Installation
command

```go get github.com/intelfike/jsonbase```

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
jb := jsonbase.New()
jb.Indent = "\t"
jb.Set().JsonText(jsonData)

// Add Data
item, _ := jb.Child(0).Clone() // value copy
item.Child("id").Set().Value(6)
item.Child("name").Set().Value("トクガワ")
jb.Push().JB(item) // like append()

// Display [class == "A"]
jb.Each(func(f *jsonbase.Jsonbase) {
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
    fmt.Println(jb.Child("a", "b").Exists()) // [false]
    jb.Child("a", "b").Set().Value(10)
    fmt.Println(jb.Child("a", "b").Exists()) // [true]
```

### Array append?

```
    jb.Push().Value(11)
```

### Make and use template?

```
    template := jsonbase.New()
    template.Set().JsonText(`{"key":"value"}`)
    jb.Push().JB(template) // jb.Child(0) == template
    jb.Push().JB(template) // jb.Child(1) == template
```

### Get child from filepath or other path?

```
    jb.ChildPath("a/b/c").Set().Set(1) // or
    jb.ChildPath("a\\b\\b").Set().Set(2) // or
    jb.ChildPath("a/b\\c").Set().Set(3) // or
    jb.ChildPath("a/b/c", "e/f", "g").Set().Set(4)
```

## type and func list

### Jsonbase

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
```

[jsonbase.go]
```
func New() *Jsonbase {
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
func (f *Jsonbase) Path() []interface{} {
func (f *Jsonbase) BottomPath() interface{} {
func (f *Jsonbase) Depth() int {
```

[setter.go]
```
func (j *Jsonbase) Set() *Setter {
func (j *Jsonbase) Push() *Setter {
func (f *Jsonbase) Remove() error {
func (f *Jsonbase) Empty() error {
```

### Setter

```
    type Setter struct{...}
```

[setter.go]
```
func (j *Jsonbase) Set() *Setter {
func (j *Jsonbase) Push() *Setter {
func (s *Setter) Value(i interface{}) error {
func (s *Setter) Reader(r io.Reader) error {
func (s *Setter) ReadFile(filename string) error {
func (s *Setter) JB(jb *Jsonbase) error {
func (s *Setter) JsonText(str string) error {
func (s *Setter) Print(a ...interface{}) error {
func (s *Setter) Printf(format string, a ...interface{}) error {
```

## Licence
MIT