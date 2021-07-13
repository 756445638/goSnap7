package snap7go

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"
)

var serverErrorsTable = map[int]error{
	0x001: errors.New("The server cannot be started."),
	0x002: errors.New("A null was passed as area pointer."),
	0x003: errors.New("Trying to re-registering an area."),
	0x004: errors.New("Area code unknown."),
	0x005: errors.New("Invalid param(s) supplied to the current function."),
	0x006: errors.New("Trying to registering too many DB (>2048)."),
	0x007: errors.New("Invalid param number suppilied to Get/SetParam."),
	0x008: errors.New("Cannot change parameter because running."),
}

func Srv_Create() (server S7Object) {
	server = C.Srv_Create()
	return
}

func Srv_Destroy(server S7Object) {
	C.Srv_Destroy((*C.S7Object)(unsafe.Pointer(&server)))
	return
}

/*
	ParamNumber 为P_u16_LocalPort的时候 value的数据是uint16 其他情况类似的
*/
//int S7API Srv_GetParam(S7Object Server, int ParamNumber, void *pValue);
func Srv_GetParam(Server S7Object, paraNumber ParamNumber) (value interface{}, err error) {
	var pValue unsafe.Pointer
	switch paraNumber {
	case P_u16_LocalPort:
		pValue = unsafe.Pointer(new(uint16))
	case P_u16_RemotePort:
		pValue = unsafe.Pointer(new(uint16))
	case P_i32_PingTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_SendTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_RecvTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_WorkInterval:
		pValue = unsafe.Pointer(new(int32))
	case P_u16_SrcRef:
		pValue = unsafe.Pointer(new(uint16))
	case P_u16_DstRef:
		pValue = unsafe.Pointer(new(uint16))
	case P_u16_SrcTSap:
		pValue = unsafe.Pointer(new(uint16))
	case P_i32_PDURequest:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_MaxClients:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_BSendTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_BRecvTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_u32_RecoveryTime:
		pValue = unsafe.Pointer(new(uint32))
	case P_u32_KeepAliveTime:
		pValue = unsafe.Pointer(new(uint32))
	}
	var code C.int = C.Srv_GetParam(Server, C.int(paraNumber), pValue)
	err = serverErrorsTable[int(code)]
	if err != nil {
		return
	}
	switch paraNumber {
	case P_u16_LocalPort:
		value = *(*uint16)(pValue)
	case P_u16_RemotePort:
		value = *(*uint16)(pValue)
	case P_i32_PingTimeout:
		value = *(*int32)(pValue)
	case P_i32_SendTimeout:
		value = *(*int32)(pValue)
	case P_i32_RecvTimeout:
		value = *(*int32)(pValue)
	case P_i32_WorkInterval:
		value = *(*int32)(pValue)
	case P_u16_SrcRef:
		value = *(*uint16)(pValue)
	case P_u16_DstRef:
		value = *(*uint16)(pValue)
	case P_u16_SrcTSap:
		value = *(*uint16)(pValue)
	case P_i32_PDURequest:
		value = *(*int32)(pValue)
	case P_i32_MaxClients:
		value = *(*int32)(pValue)
	case P_i32_BSendTimeout:
		value = *(*int32)(pValue)
	case P_i32_BRecvTimeout:
		value = *(*int32)(pValue)
	case P_u32_RecoveryTime:
		value = *(*uint32)(pValue)
	case P_u32_KeepAliveTime:
		value = *(*uint32)(pValue)
	}
	return
}

func Srv_SetParam(Server S7Object, paraNumber ParamNumber, value interface{}) (err error) {
	var pValue unsafe.Pointer
	pValue = Value_Pvalue(paraNumber,value)
	var code C.int = C.Srv_SetParam(Server, C.int(paraNumber), pValue)
	err = serverErrorsTable[int(code)]
	return
}

// int S7API Srv_StartTo(S7Object Server, const char *Address);
func Srv_StartTo(Server S7Object, Address string) (err error) {
	address := C.CString(Address)
	defer func() {
		C.free(unsafe.Pointer(address))
	}()
	var code C.int = C.Srv_StartTo(Server, address)
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_Start(S7Object Server);
func Srv_Start(Server S7Object) (err error) {
	var code C.int = C.Srv_Start(Server)
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_Stop(S7Object Server)
func Srv_Stop(Server S7Object) (err error) {
	var code C.int = C.Srv_Stop(Server)
	err = serverErrorsTable[int(code)]
	return
}

//typedef uint16_t   word;
// func Srv_RegisterArea(S7Object Server, int AreaCode, word Index, void *pUsrData, int Size)
func Srv_RegisterArea(Server S7Object, AreaCode int, Index uint16, pUsrData []byte, Size int) (err error) {
	var code C.int = C.Srv_RegisterArea(Server,C.int(AreaCode),C.uint16_t(Index),unsafe.Pointer(&pUsrData[0]),C.int(Size))
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_UnregisterArea(S7Object Server, int AreaCode, word Index);
func Srv_UnregisterArea(Server S7Object, AreaCode int, Index uint16) (err error) {
	var code C.int = C.Srv_UnregisterArea(Server,C.int(AreaCode),C.uint16_t(Index))
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_LockArea(S7Object Server, int AreaCode, word Index);
func Srv_LockArea(Server S7Object, AreaCode int, Index uint16) (err error) {
	var code C.int = C.Srv_LockArea(Server,C.int(AreaCode),C.uint16_t(Index))
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_UnlockArea(S7Object Server, int AreaCode, word Index);
func Srv_UnlockArea(Server S7Object, AreaCode int, Index uint16) (err error) {
	var code C.int = C.Srv_UnlockArea(Server,C.int(AreaCode),C.uint16_t(Index))
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_GetStatus(S7Object Server, int *ServerStatus, int *CpuStatus, int *ClientsCount);
func Srv_GetStatus( Server S7Object, CpuStatus int, ClientsCount int) (ServerStatus int,err error) {
	var code C.int = C.Srv_GetStatus(Server,(*C.int)(unsafe.Pointer(&ServerStatus)),(*C.int)(unsafe.Pointer(&CpuStatus)),(*C.int)(unsafe.Pointer(&ClientsCount)))
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_SetCpuStatus(S7Object Server, int CpuStatus);
func Srv_SetCpuStatus(Server S7Object , CpuStatus int) (err error) {
	var code C.int = C.Srv_SetCpuStatus(Server,C.int(CpuStatus))
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_ClearEvents(S7Object Server);
func Srv_ClearEvents( Server S7Object) (err error) {
	var code C.int = C.Srv_ClearEvents(Server)
	err = serverErrorsTable[int(code)]
	return
}

// func Srv_PickEvent(S7Object Server, TSrvEvent *pEvent, int *EvtReady);
func Srv_PickEvent(Server S7Object, pEvent TSrvEvent , EvtReady int )  (err error){
	var code C.int = C.Srv_PickEvent(Server,(*C.TSrvEvent)(unsafe.Pointer(&pEvent)),(*C.int)(unsafe.Pointer(&EvtReady)))
	 err = serverErrorsTable[int(code)]
	return
}

// func Srv_GetMask(S7Object Server, int MaskKind, longword *Mask);  uint32_t
func Srv_GetMask(Server S7Object ,MaskKind int) ( Mask uint32,err error) {
	var code C.int = C.Srv_GetMask(Server,C.int(MaskKind),(*C.uint32_t)((unsafe.Pointer(&Mask))))
	err = serverErrorsTable[int(code)]
	return
}
// func Srv_SetMask(S7Object Server, int MaskKind, longword Mask);
func Srv_SetMask( Server S7Object, MaskKind int , Mask uint32 ) (err error) {
	var code C.int = C.Srv_SetMask(Server,C.int(MaskKind),C.uint32_t(Mask))
	err = serverErrorsTable[int(code)]
	return
}


// func Srv_SetEventsCallback(S7Object Server, pfn_SrvCallBack pCallback, void *usrPtr);
// func Srv_SetReadEventsCallback(S7Object Server, pfn_SrvCallBack pCallback, void *usrPtr);
// func Srv_SetRWAreaCallback(S7Object Server, pfn_RWAreaCallBack pCallback, void *usrPtr);


// func Srv_EventText(TSrvEvent *Event, char *Text, int TextLen)
func Srv_EventText(Event TSrvEvent, Text string, TextLen int) (err error){
	text := C.CString(Text)
	defer func() {
		C.free(unsafe.Pointer(text))
	}()
	var code C.int = C.Srv_EventText((*C.TSrvEvent)(unsafe.Pointer(&Event)),text,C.int(TextLen))
	err = serverErrorsTable[int(code)]
	return
}
// func Srv_ErrorText(int Error, char *Text, int TextLen);
func Srv_ErrorText( Error int, Text string,  TextLen int)  (err error) {
	text := C.CString(Text)
	defer func() {
		C.free(unsafe.Pointer(text))
	}()
	var code C.int = C.Srv_ErrorText(C.int(Error),text,C.int(TextLen))
	err = serverErrorsTable[int(code)]
	return
}