package jsonbase

func (f *Jsonbase) Exists() bool {
	return f.ReferError() == nil
}
func (f *Jsonbase) ReferError() error {
	_, err := f.GetInterfacePt()
	return err
}

func (f *Jsonbase) IsString() bool {
	_, ok := (f.Interface()).(string)
	return ok
}
func (f *Jsonbase) IsBool() bool {
	_, ok := (f.Interface()).(bool)
	return ok
}
func (f *Jsonbase) IsInt() bool {
	_, ok := (f.Interface()).(int)
	return ok
}
func (f *Jsonbase) IsUint() bool {
	_, ok := (f.Interface()).(uint)
	return ok
}
func (f *Jsonbase) IsFloat() bool {
	_, ok := (f.Interface()).(float64)
	return ok
}
func (f *Jsonbase) IsNull() bool {
	if !f.Exists() {
		return false
	}
	return f.Interface() == nil
}
func (f *Jsonbase) IsArray() bool {
	_, ok := (f.Interface()).([]interface{})
	return ok
}
func (f *Jsonbase) IsMap() bool {
	_, ok := (f.Interface()).(map[string]interface{})
	return ok
}

func (f *Jsonbase) HasChild(a interface{}) bool {
	return f.Child(a).Exists()
}
