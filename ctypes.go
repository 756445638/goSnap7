// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs ctypes/ctypes.go

package snap7go

type TS7CpuInfo struct {
	ModuleTypeName	[33]int8
	SerialNumber	[25]int8
	ASName		[25]int8
	Copyright	[27]int8
	ModuleName	[25]int8
}
type Tm struct {
	Sec	int32
	Min	int32
	Hour	int32
	Mday	int32
	Mon	int32
	Year	int32
	Wday	int32
	Yday	int32
	Isdst	int32
}

type TS7DataItem struct {
	Area		int32
	WordLen		int32
	Result		int32
	DBNumber	int32
	Start		int32
	Amount		int32
	Pdata		*byte
}
type TS7BlocksList struct {
	OBCount		int32
	FBCount		int32
	FCCount		int32
	SFBCount	int32
	SFCCount	int32
	DBCount		int32
	SDBCount	int32
}
type TS7BlockInfo struct {
	BlkType		int32
	BlkNumber	int32
	BlkLang		int32
	BlkFlags	int32
	MC7Size		int32
	LoadSize	int32
	LocalData	int32
	SBBLength	int32
	CheckSum	int32
	Version		int32
	CodeDate	[11]int8
	IntfDate	[11]int8
	Author		[9]int8
	Family		[9]int8
	Header		[9]int8
}
type TSrvEvent struct {
	EvtTime		int64
	EvtSender	int32
	EvtCode		uint32
	EvtRetCode	uint16
	EvtParam1	uint16
	EvtParam2	uint16
	EvtParam3	uint16
	EvtParam4	uint16
}

type TS7BlocksOfType [8192]uint16

type PS7Tag struct {
	Area		int32
	DBNumber	int32
	Start		int32
	Size		int32
	WordLen		int32
}

