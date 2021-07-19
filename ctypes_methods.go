package snap7go

/*
#include <time.h>
#include <stdio.h>
#include <stdlib.h>
typedef struct tm* Tm;
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

func convertInt8SliceToString(s []int8) string {
	ss := make([]byte, len(s))
	for k, v := range s {
		ss[k] = byte(v)
	}
	return string(ss)
}

func (t *TS7CpuInfo) GetModuleTypeName() string {
	return convertInt8SliceToString(t.ModuleTypeName[:])
}
func (t *TS7CpuInfo) GetSerialNumber() string {
	return convertInt8SliceToString(t.SerialNumber[:])
}
func (t *TS7CpuInfo) GetASName() string {
	return convertInt8SliceToString(t.ASName[:])
}
func (t *TS7CpuInfo) GetCopyright() string {
	return convertInt8SliceToString(t.Copyright[:])
}
func (t *TS7CpuInfo) GetModuleName() string {
	return convertInt8SliceToString(t.ModuleName[:])
}

func (t *TS7Protection) GetProtectionString() string {
	return fmt.Sprintf("%+v", t)
}

/*
https://www.runoob.com/cprogramming/c-standard-library-time-h.html
*/
type time_t = C.time_t

func (t *Tm) FromTime(goTime time.Time) {
	var time_t = time_t(goTime.Unix())
	*t = *(*Tm)(unsafe.Pointer(C.localtime(&time_t)))
}

func (t *Tm) ToTime() time.Time {
	x := (C.Tm)(unsafe.Pointer(t))
	var s = C.mktime(x)
	return time.Unix(int64(s), 0)
}
