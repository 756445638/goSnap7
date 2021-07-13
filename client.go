package snap7go

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
import "C"
import (
	"errors"
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
func Cli_ReadMultiVars(cli S7Object, items []TS7DataItem) (datas [][]byte, err error) {
	itemsCount := len(items)
	for _, v := range items {
		t := make([]byte, dataLength(S7WL(v.WordLen), v.Amount, v.Start))
		datas = append(datas, t)
		v.Pdata = &t[0]
	}
	var code C.int = C.Cli_ReadMultiVars(cli, (C.PS7DataItem)(unsafe.Pointer(&items[0])), C.int(itemsCount))
	err = Cli_ErrorText(code)
	return
}
func Cli_WriteMultiVars(cli S7Object, items []TS7DataItem) (datas [][]byte, err error) {
	itemsCount := len(items)
	var code C.int = C.Cli_WriteMultiVars(cli, (C.PS7DataItem)(unsafe.Pointer(&items[0])), C.int(itemsCount))
	err = Cli_ErrorText(code)
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
func Cli_Upload(cli S7Object, blockType Block, blockNum int, pUsrData []byte, size int) (err error) {
	pUsrData = make([]byte, size)
	var code C.int = C.Cli_Upload(cli, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_FullUpload(S7Object Client, int BlockType, int BlockNum, void *pUsrData, int *Size);
func Cli_FullUpload(cli S7Object, blockType Block, blockNum int, pUsrData []byte, size int) (err error) {
	pUsrData = make([]byte, size)
	var code C.int = C.Cli_FullUpload(cli, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

func Cli_GetCpuInfo(cli S7Object) (info TS7CpuInfo, err error) {
	var code C.int = C.Cli_GetCpuInfo(cli, (*C.TS7CpuInfo)(unsafe.Pointer(&info)))
	err = Cli_ErrorText(code)
	return
}

var cliErrorsTable = map[int]error{
	0x001: errors.New("error during PDU negotiation."),
	0x002: errors.New("Invalid param(s) supplied to the current function."),
	0x003: errors.New("A Job is pending : there is an async function in progress."),
	0x004: errors.New("More than 20 items where passed to a MultiRead/Write area function."),
	0x005: errors.New("Invalid Wordlen param supplied to the current function"),
	0x006: errors.New("Partial data where written : The target area is smaller than the DataSize supplied."),
	0x007: errors.New("A MultiRead/MultiWrite function has datasize over the PDU size."),
	0x008: errors.New("Invalid answer from the PLC."),
	0x009: errors.New("An address out of range was specified."),
	0x00A: errors.New("Invalid Transportsize parameter was supplied to a Read/WriteArea function."),
	0x00B: errors.New("Invalid datasize parameter supplied to the current function."),
	0x00C: errors.New("Item requested was not found in the PLC."),
	0x00D: errors.New("Invalid value supplied to the current function."),
	0x00E: errors.New("PLC cannot be started."),
	0x00F: errors.New("PLC is already in RUN stare."),
	0x010: errors.New("PLC cannot be stopped."),
	0x011: errors.New("Cannot copy RAM to ROM : the PLC is running or doesn’t support this function."),
	0x012: errors.New("Cannot compress : the PLC is running or doesn’t support this function."),
	0x013: errors.New("PLC is already in STOP state."),
	0x014: errors.New("Function not available."),
	0x015: errors.New("Block upload sequence failed."),
	0x016: errors.New("Invalid data size received from the PLC."),
	0x017: errors.New("Invalid block type supplied to the current function."),
	0x018: errors.New("Invalid block supplied to the current function."),
	0x019: errors.New("Invalid block size supplied to the current function."),
	0x01A: errors.New("Block download sequence failed."),
	0x01B: errors.New("Insert command (implicit command sent after a block download) refused."),
	0x01C: errors.New("Delete command refused."),
	0x01D: errors.New("This operation is password protected."),
	0x01E: errors.New("Invalid password supplied."),
	0x01F: errors.New("There is no password to set or clear : the protection is OFF."),
	0x020: errors.New("Job timeout."),
	0x021: errors.New("Partial data where read : The source area is greater than the DataSize supplied."),
	0x022: errors.New("The buffer supplied is too small."),
	0x023: errors.New("Function refused by the PLC."),
	0x024: errors.New("Invalid param number suppilied to Get/SetParam."),
	0x025: errors.New("Cannot perform : the client is destroying."),
	0x026: errors.New("Cannot change parameter because connected."),
}
