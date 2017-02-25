package filebase

func (f *Filebase) Exists() bool {
	_, err := f.GetInterfacePt()
	return err == nil
}
func (f *Filebase) IsString() bool {
	_, ok := (f.Interface()).(string)
	return ok
}
func (f *Filebase) IsBool() bool {
	_, ok := (f.Interface()).(bool)
	return ok
}
func (f *Filebase) IsInt() bool {
	_, ok := (f.Interface()).(int)
	return ok
}
func (f *Filebase) IsUint() bool {
	_, ok := (f.Interface()).(uint)
	return ok
}
func (f *Filebase) IsFloat() bool {
	_, ok := (f.Interface()).(float64)
	return ok
}
func (f *Filebase) IsNull() bool {
	if !f.Exists() {
		return false
	}
	return f.Interface() == nil
}
func (f *Filebase) IsArray() bool {
	_, ok := (f.Interface()).([]interface{})
	return ok
}
func (f *Filebase) IsMap() bool {
	_, ok := (f.Interface()).(map[string]interface{})
	return ok
}
