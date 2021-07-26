package snap7go

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
/*
	extern void GlobalEventsCallback(void *usrPtr, PSrvEvent PEvent, int Size);
	extern void GlobalReadEventsCallback(void *usrPtr, PSrvEvent PEvent, int Size);

	extern int GlobalRWAreaCallback(void *usrPtr, int Sender, int Operation, PS7Tag PTag, void *pUsrData);

*/
import "C"
import "unsafe"

type Pfn_SrvEventCallBack = func(uintptr, *TSrvEvent)

type Pfn_RWAreaCallBack = func(uintptr, int32, Operation, *PS7Tag, uintptr) int32

var (
	svrEventCallbacks      = make(map[uintptr]Pfn_SrvEventCallBack)
	svrReadEventsCallbacks = make(map[uintptr]Pfn_SrvEventCallBack)
	svrRWAreaCallbacks     = make(map[uintptr]Pfn_RWAreaCallBack)
)

// int S7API Srv_SetEventsCallback(S7Object Server, pfn_SrvCallBack pCallback, void *usrPtr);
func Srv_SetEventsCallback(svr S7Object, handle Pfn_SrvEventCallBack, usrPtr uintptr) error {
	var code C.int = C.Srv_SetEventsCallback(svr, (*[0]byte)(C.GlobalEventsCallback), unsafe.Pointer(usrPtr))
	err := Srv_ErrorText(code)
	if err != nil {
		return err
	}
	svrEventCallbacks[usrPtr] = handle
	return nil
}

//int S7API Srv_SetReadEventsCallback(S7Object Server, pfn_SrvCallBack pCallback, void *usrPtr);
func Srv_SetReadEventsCallback(svr S7Object, handle Pfn_SrvEventCallBack, usrPtr uintptr) error {
	var code C.int = C.Srv_SetReadEventsCallback(svr, (*[0]byte)(C.GlobalReadEventsCallback), unsafe.Pointer(usrPtr))
	err := Srv_ErrorText(code)
	if err != nil {
		return err
	}
	svrReadEventsCallbacks[usrPtr] = handle
	return nil
}

//int S7API Srv_SetRWAreaCallback(S7Object Server, pfn_RWAreaCallBack pCallback, void *usrPtr);
func Srv_SetRWAreaCallback(svr S7Object, handle Pfn_RWAreaCallBack, usrPtr uintptr) error {
	var code C.int = C.Srv_SetRWAreaCallback(svr, (*[0]byte)(C.GlobalRWAreaCallback), unsafe.Pointer(usrPtr))
	err := Srv_ErrorText(code)
	if err != nil {
		return err
	}
	svrRWAreaCallbacks[usrPtr] = handle
	return nil
}

func getSrvEventFromC(e C.PSrvEvent) (ego TSrvEvent) {
	ego = *((*TSrvEvent)(unsafe.Pointer(e)))
	return
}

//export GlobalEventsCallback
func GlobalEventsCallback(usrPtr *C.void, event C.PSrvEvent, size C.int) {
	// xxx(uintptr(event))
	up := uintptr(unsafe.Pointer(usrPtr))
	callback := svrEventCallbacks[up]
	if callback == nil {
		return
	}
	e := getSrvEventFromC(event)
	callback(up, &e)
}

//export GlobalReadEventsCallback
func GlobalReadEventsCallback(usrPtr *C.void, event C.PSrvEvent, size C.int) {
	up := uintptr(unsafe.Pointer(usrPtr))
	callback := svrReadEventsCallbacks[up]
	if callback == nil {
		return
	}
	e := getSrvEventFromC(event)
	callback(up, &e)
}
func getPS7TagFormC(t C.PS7Tag) (et PS7Tag) {
	et = *((*PS7Tag)(unsafe.Pointer(t)))
	return
}

//export GlobalRWAreaCallback
func GlobalRWAreaCallback(usrPtr *C.void, sender C.int, operation C.int, pTag C.PS7Tag, pUserData *C.void) C.int {
	up := uintptr(unsafe.Pointer(usrPtr))
	callback := svrRWAreaCallbacks[up]
	if callback == nil {
		return 0 // no callback no error
	}

	pt := getPS7TagFormC(pTag)
	return C.int(callback(
		up,
		int32(sender),
		Operation(operation),
		&pt,
		uintptr(unsafe.Pointer(pUserData))))
}

/*

todo GlobalEventsCallback 最后一个参数Size是？？？
用demo跑了一下 看起来都是26 文档里面说了为了以后兼容

yuyang@yuyang-PC:~/projects/snap7-full-1.4.2/snap7-full-1.4.2/examples/plain-c/x86_64-linux$ sudo ./server
2021-07-13 10:39:14 Server started
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Client added
 Size:26
2021-07-13 10:39:17 [127.0.0.1] The client requires a PDU size of 480 bytes
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read SZL request, ID:0x0011 INDEX:0x0000 --> OK
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read SZL request, ID:0x001c INDEX:0x0000 --> OK
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read SZL request, ID:0x0131 INDEX:0x0001 --> OK
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read SZL request, ID:0x0424 INDEX:0x0000 --> OK
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read SZL request, ID:0x0011 INDEX:0x0000 --> OK
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Block upload requested --> NOT PERFORMED (due to invalid security level)
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Block upload requested --> NOT PERFORMED (due to invalid security level)
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Block upload requested --> NOT PERFORMED (due to invalid security level)
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Block upload requested --> NOT PERFORMED (due to invalid security level)
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read request, Area : MK, Start : 0, Size : 0 --> Area not found
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read request, Area : PE, Start : 0, Size : 0 --> Area not found
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read request, Area : PA, Start : 0, Size : 0 --> Area not found
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read request, Area : TM, Start : 0, Size : 0 --> Area not found
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Read request, Area : CT, Start : 0, Size : 0 --> Area not found
 Size:26
2021-07-13 10:39:17 [127.0.0.1] Client disconnected by peer
 Size:26


*/
