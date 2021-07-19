package snap7go

func convertInt8SliceToString(s []int8) string {
	ss := make([]byte, len(s))
	for k, v := range s {
		ss[k] = byte(v)
	}
	return string(ss)
}

func (t TS7CpuInfo) GetModuleTypeName() string {
	return convertInt8SliceToString(t.ModuleTypeName[:])
}
func (t TS7CpuInfo) GetSerialNumber() string {
	return convertInt8SliceToString(t.SerialNumber[:])
}
func (t TS7CpuInfo) GetASName() string {
	return convertInt8SliceToString(t.ASName[:])
}
func (t TS7CpuInfo) GetCopyright() string {
	return convertInt8SliceToString(t.Copyright[:])
}
func (t TS7CpuInfo) GetModuleName() string {
	return convertInt8SliceToString(t.ModuleName[:])
}

//func (t TS7Protection) String() string {
//
//}
