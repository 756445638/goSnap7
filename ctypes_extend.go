package snap7go

import "unsafe"

//#include <malloc.h>
import "C"

type TS7DataItemGo struct {
	Area     int32
	WordLen  int32
	Result   int32
	DBNumber int32
	Start    int32
	Amount   int32

	/*
		读取和写值的数据
		如果是读这个字段不需要初始化
	*/

	Pdata []byte
}

/*
	Pdata没有处理
	读和写需要单独处理
*/
func (g *TS7DataItemGo) ToC() TS7DataItem {
	return TS7DataItem{
		Area:     g.Area,
		WordLen:  g.WordLen,
		Result:   g.Result,
		DBNumber: g.DBNumber,
		Start:    g.Start,
		Amount:   g.Amount,
	}
}

/*

 */
func (g *TS7DataItem) GetBytes() (data []byte) {
	p := uintptr(unsafe.Pointer(g.Pdata))
	length := dataLength(S7WL(g.WordLen), g.Amount)

	data = GetBytesFromC(p, int(length))
	return data
}

func (g *TS7DataItemGo) CopyPdata(to *TS7DataItem) {
	length := dataLength(S7WL(g.WordLen), g.Amount)
	to.Pdata = (*byte)(C.malloc(C.size_t(length)))
	up := uintptr(unsafe.Pointer(to.Pdata))
	CopyToC(g.Pdata, up)
}

func CopyToC(bs []byte, c uintptr) {
	length := len(bs)
	for i := 0; i < int(length); i++ {
		*((*byte)(unsafe.Pointer(c + uintptr(i)))) = bs[i]
	}
}

func GetBytesFromC(userData uintptr, length int) (data []byte) {
	data = make([]byte, length)
	for i := 0; i < int(length); i++ {
		data[i] = *(*byte)(unsafe.Pointer(userData + uintptr(i)))
	}
	return
}
