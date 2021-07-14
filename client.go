package snap7go

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

type S7Object = C.S7Object

func Cli_Create() (cli S7Object) {
	cli = C.Cli_Create()
	return
}
func Cli_Destroy(cli S7Object) {
	C.Cli_Destroy((*C.S7Object)(unsafe.Pointer(&cli)))
	return
}
func Cli_ConnectTo(cli S7Object, address string, rack int, slot int) (err error) {
	s := C.CString(address)
	defer func() {
		C.free(unsafe.Pointer(s))
	}()
	var code C.int = C.Cli_ConnectTo(cli, s, C.int(rack), C.int(slot))
	err = Cli_ErrorText(code)
	return
}

func Cli_SetConnectionParams(cli S7Object, address string, localTSAP uint16, remoteTSAP uint16) (err error) {
	s := C.CString(address)
	defer func() {
		C.free(unsafe.Pointer(s))
	}()
	var code C.int = C.Cli_SetConnectionParams(cli, s, C.word(localTSAP), C.word(remoteTSAP))
	err = Cli_ErrorText(code)
	return
}
func Cli_SetConnectionType(cli S7Object, connectionType CONNTYPE) (err error) {
	var code C.int = C.Cli_SetConnectionType(cli, C.word(connectionType))
	err = Cli_ErrorText(code)
	return
}
func Cli_Connect(cli S7Object) (err error) {
	var code C.int = C.Cli_Connect(cli)
	err = Cli_ErrorText(code)
	return
}
func Cli_Disconnect(cli S7Object) (err error) {
	var code C.int = C.Cli_Disconnect(cli)
	err = Cli_ErrorText(code)
	return
}

/*
	ParamNumber 为P_u16_LocalPort的时候 value的数据是uint16 其他情况类似的
*/
func Cli_GetParam(cli S7Object, paraNumber ParamNumber) (value interface{}, err error) {
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
	var code C.int = C.Cli_GetParam(cli, C.int(paraNumber), pValue)
	err = Cli_ErrorText(code)
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

/*
	P_u16_LocalPort 设定端口为uint16
*/
func Cli_SetParam(cli S7Object, paraNumber ParamNumber, value interface{}) (err error) {
	pvalue := Value_Pvalue(paraNumber, value)
	var code C.int = C.Cli_SetParam(cli, C.int(paraNumber), pvalue)
	err = Cli_ErrorText(code)
	return
}
func Cli_ReadArea(cli S7Object, area S7Area, dBNumber int, start int, amount int, wordLen S7WL) (pUsrData []byte, err error) {
	pUsrData = make([]byte, dataLength(wordLen, int32(amount), int32(start)))
	var code C.int = C.Cli_ReadArea(cli, C.int(area), C.int(dBNumber), C.int(start), C.int(amount), C.int(wordLen), unsafe.Pointer(&pUsrData[0]))
	err = Cli_ErrorText(code)
	return
}

func Cli_WriteArea(cli S7Object, area S7Area, dBNumber int, start int, amount int, wordLen S7WL, pUsrData []byte) (err error) {
	pUsrData = make([]byte, dataLength(wordLen, int32(amount), int32(start)))
	var code C.int = C.Cli_WriteArea(cli, C.int(area), C.int(dBNumber), C.int(start), C.int(amount), C.int(wordLen), unsafe.Pointer(&pUsrData[0]))
	err = Cli_ErrorText(code)
	return
}

func Cli_ReadMultiVars(cli S7Object, items []TS7DataItemGo) (err error) {
	itemsCount := len(items)
	itemsC := make([]TS7DataItem, itemsCount)

	for k, v := range items {
		t := make([]byte, dataLength(S7WL(v.WordLen), v.Amount, v.Start))
		v.Pdata = t
		itemsC[k] = v.ToC()
	}
	var code C.int = C.Cli_ReadMultiVars(cli, (C.PS7DataItem)(unsafe.Pointer(&itemsC[0])), C.int(itemsCount))
	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	for k, v := range itemsC {
		items[k].Result = v.Result
	}
	return
}
func Cli_WriteMultiVars(cli S7Object, items []TS7DataItemGo) (err error) {
	itemsCount := len(items)
	itemsC := make([]TS7DataItem, itemsCount)
	for k, v := range items {
		itemsC[k] = v.ToC()
	}
	var code C.int = C.Cli_WriteMultiVars(cli, (C.PS7DataItem)(unsafe.Pointer(&itemsC[0])), C.int(itemsCount))
	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	for k, v := range itemsC {
		items[k].Result = v.Result
	}
	return
}

func Cli_DBRead(cli S7Object, dBNumber int, start int, amount int) (pUsrData []byte, err error) {
	return Cli_ReadArea(cli, S7AreaDB, dBNumber, start, amount, S7WLByte)
}
func Cli_DBWrite(cli S7Object, dBNumber int, start int, amount int, pUsrData []byte) (err error) {
	return Cli_WriteArea(cli, S7AreaDB, dBNumber, start, amount, S7WLByte, pUsrData)
}
func Cli_MBRead(cli S7Object, dBNumber int, start int, amount int) (pUsrData []byte, err error) {
	return Cli_ReadArea(cli, S7AreaMK, dBNumber, start, amount, S7WLByte)
}
func Cli_MBWrite(cli S7Object, dBNumber int, start int, amount int, pUsrData []byte) (err error) {
	return Cli_WriteArea(cli, S7AreaMK, dBNumber, start, amount, S7WLByte, pUsrData)
}
func Cli_EBRead(cli S7Object, dBNumber int, start int, amount int) (pUsrData []byte, err error) {
	return Cli_ReadArea(cli, S7AreaPE, dBNumber, start, amount, S7WLByte)
}
func Cli_EBWrite(cli S7Object, dBNumber int, start int, amount int, pUsrData []byte) (err error) {
	return Cli_WriteArea(cli, S7AreaPE, dBNumber, start, amount, S7WLByte, pUsrData)
}
func Cli_ABRead(cli S7Object, dBNumber int, start int, amount int) (pUsrData []byte, err error) {
	return Cli_ReadArea(cli, S7AreaPA, dBNumber, start, amount, S7WLByte)
}
func Cli_ABWrite(cli S7Object, dBNumber int, start int, amount int, pUsrData []byte) (err error) {
	return Cli_WriteArea(cli, S7AreaPA, dBNumber, start, amount, S7WLByte, pUsrData)
}
func Cli_TMRead(cli S7Object, dBNumber int, start int, amount int) (pUsrData []byte, err error) {
	return Cli_ReadArea(cli, S7AreaTM, dBNumber, start, amount, S7WLByte)
}
func Cli_TMWrite(cli S7Object, dBNumber int, start int, amount int, pUsrData []byte) (err error) {
	return Cli_WriteArea(cli, S7AreaTM, dBNumber, start, amount, S7WLByte, pUsrData)
}
func Cli_CTRead(cli S7Object, dBNumber int, start int, amount int) (pUsrData []byte, err error) {
	return Cli_ReadArea(cli, S7AreaCT, dBNumber, start, amount, S7WLByte)
}
func Cli_CTWrite(cli S7Object, dBNumber int, start int, amount int, pUsrData []byte) (err error) {
	return Cli_WriteArea(cli, S7AreaCT, dBNumber, start, amount, S7WLByte, pUsrData)
}
func Cli_ListBlocks(cli S7Object) (pUsrData TS7BlocksList, err error) {
	var code C.int = C.Cli_ListBlocks(cli, (*C.TS7BlocksList)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}
func Cli_GetAgBlockInfo(cli S7Object, blockType Block, blockNum int) (pUsrData TS7BlockInfo, err error) {
	var code C.int = C.Cli_GetAgBlockInfo(cli, C.int(blockType), C.int(blockNum), (*C.TS7BlockInfo)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetPgBlockInfo(S7Object Client, void *pBlock, TS7BlockInfo *pUsrData, int Size);
func Cli_GetPgBlockInfo(cli S7Object, size int) (pBlock []byte, pUsrData TS7BlockInfo, err error) {
	pBlock = make([]byte, size)
	var code C.int = C.Cli_GetPgBlockInfo(cli, unsafe.Pointer(&pBlock[0]), (*C.TS7BlockInfo)(unsafe.Pointer(&pUsrData)), C.int(size))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ListBlocksOfType(S7Object Client, int BlockType, TS7BlocksOfType *pUsrData, int *ItemsCount);
func Cli_ListBlocksOfType(cli S7Object, blockType Block) (pUsrData TS7BlocksOfType, itemsCount int, err error) {
	var code C.int = C.Cli_ListBlocksOfType(cli, C.int(blockType), (*C.TS7BlocksOfType)(unsafe.Pointer(&pUsrData[0])), (*C.int)(unsafe.Pointer(&itemsCount)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_Upload(S7Object Client, int BlockType, int BlockNum, void *pUsrData, int *Size);
func Cli_Upload(cli S7Object, blockType Block, blockNum int, pUsrData []byte) (size int, err error) {
	size = len(pUsrData)
	var code C.int = C.Cli_Upload(cli, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_FullUpload(S7Object Client, int BlockType, int BlockNum, void *pUsrData, int *Size);
func Cli_FullUpload(cli S7Object, blockType Block, blockNum int, pUsrData []byte) (size int, err error) {
	size = len(pUsrData)
	var code C.int = C.Cli_FullUpload(cli, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_Download(S7Object Client, int BlockNum, void *pUsrData, int Size);
func Cli_Download(cli S7Object, blockNum int, pUsrData []byte, size int) (err error) {
	pUsrData = make([]byte, size)
	var code C.int = C.Cli_Download(cli, C.int(blockNum), unsafe.Pointer(&pUsrData[0]), C.int(size))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_Delete(S7Object Client, int BlockType, int BlockNum);
func Cli_Delete(cli S7Object, blockType Block, blockNum int) (err error) {
	var code C.int = C.Cli_Delete(cli, C.int(blockType), C.int(blockNum))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_DBGet(S7Object Client, int DBNumber, void *pUsrData, int *Size);
func Cli_DBGet(cli S7Object, dBNumber int, pUsrData []byte) (size int, err error) {
	return Cli_Upload(cli, Block_DB, dBNumber, pUsrData)
}

//int S7API Cli_DBFill(S7Object Client, int DBNumber, int FillChar);
func Cli_DBFill(cli S7Object, dBNumber int, fillChar int) (err error) {
	var code C.int = C.Cli_DBFill(cli, C.int(dBNumber), C.int(fillChar))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetPlcDateTime(S7Object Client, tm *DateTime);
func Cli_GetPlcDateTime(cli S7Object) (dataTime Tm, err error) {
	var code C.int = C.Cli_GetPlcDateTime(cli, (*C.tm)(unsafe.Pointer(&dataTime)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_SetPlcDateTime(S7Object Client, tm *DateTime);
func Cli_SetPlcDateTime(cli S7Object, dataTime Tm) (err error) {
	var code C.int = C.Cli_SetPlcDateTime(cli, (*C.tm)(unsafe.Pointer(&dataTime)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_SetPlcSystemDateTime(S7Object Client);
func Cli_SetPlcSystemDateTime(cli S7Object) (err error) {
	var code C.int = C.Cli_SetPlcSystemDateTime(cli)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetOrderCode(S7Object Client, TS7OrderCode *pUsrData);
func Cli_GetOrderCode(cli S7Object) (pUsrData TS7OrderCode, err error) {
	var code C.int = C.Cli_GetOrderCode(cli, (*C.TS7OrderCode)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetCpuInfo(S7Object Client, TS7CpuInfo *pUsrData);
func Cli_GetCpuInfo(cli S7Object) (pUsrData TS7CpuInfo, err error) {
	var code C.int = C.Cli_GetCpuInfo(cli, (*C.TS7CpuInfo)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetCpInfo(S7Object Client, TS7CpInfo *pUsrData);
func Cli_GetCpInfo(cli S7Object) (pUsrData TS7CpInfo, err error) {
	var code C.int = C.Cli_GetCpInfo(cli, (*C.TS7CpInfo)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ReadSZL(S7Object Client, int ID, int Index, TS7SZL *pUsrData, int *Size);
func Cli_ReadSZL(cli S7Object, id int, index int) (pUsrData TS7SZL, size int, err error) {
	var code C.int = C.Cli_ReadSZL(cli, C.int(id), C.int(index), (*C.TS7SZL)(unsafe.Pointer(&pUsrData)), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ReadSZLList(S7Object Client, TS7SZLList *pUsrData, int *ItemsCount);
func Cli_ReadSZLList(cli S7Object) (pUsrData TS7SZLList, itemsCount int, err error) {
	var code C.int = C.Cli_ReadSZLList(cli, (*C.TS7SZLList)(unsafe.Pointer(&pUsrData)), (*C.int)(unsafe.Pointer(&itemsCount)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_PlcHotStart(S7Object Client);
func Cli_PlcHotStart(cli S7Object) (err error) {
	var code C.int = C.Cli_PlcHotStart(cli)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_PlcColdStart(S7Object Client);
func Cli_PlcColdStart(cli S7Object) (err error) {
	var code C.int = C.Cli_PlcColdStart(cli)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_PlcStop(S7Object Client);
func Cli_PlcStop(cli S7Object) (err error) {
	var code C.int = C.Cli_PlcStop(cli)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_CopyRamToRom(S7Object Client, int Timeout);
func Cli_CopyRamToRom(cli S7Object, timeout int) (err error) {
	var code C.int = C.Cli_CopyRamToRom(cli, C.int(timeout))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_Compress(S7Object Client, int Timeout);
func Cli_Compress(cli S7Object, timeout int) (err error) {
	var code C.int = C.Cli_Compress(cli, C.int(timeout))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetPlcStatus(S7Object Client, int *Status);
func Cli_GetPlcStatus(cli S7Object) (status S7CpuStatus, err error) {
	var code C.int = C.Cli_GetPlcStatus(cli, (*C.int)(unsafe.Pointer(&status)))
	err = Cli_ErrorText(code)
	return
}
