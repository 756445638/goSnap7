package snap7go

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func Cli_Create() (cli S7Object) {
	cli = C.Cli_Create()
	return
}

type S7Object = C.S7Object

func (c *S7Client) Destroy() {
	C.Cli_Destroy((*C.S7Object)(unsafe.Pointer(&c.client)))
	return
}
func (c *S7Client) ConnectTo(address string, rack int32, slot int32) (err error) {
	s := C.CString(address)
	defer func() {
		C.free(unsafe.Pointer(s))
	}()
	var code C.int = C.Cli_ConnectTo(c.client, s, C.int(rack), C.int(slot))
	err = Cli_ErrorText(code)
	return
}

func (c *S7Client) SetConnectionParams(address string, localTSAP uint16, remoteTSAP uint16) (err error) {
	s := C.CString(address)
	defer func() {
		C.free(unsafe.Pointer(s))
	}()
	var code C.int = C.Cli_SetConnectionParams(c.client, s, C.word(localTSAP), C.word(remoteTSAP))
	err = Cli_ErrorText(code)
	return
}
func (c *S7Client) SetConnectionType(connectionType CONNTYPE) (err error) {
	var code C.int = C.Cli_SetConnectionType(c.client, C.word(connectionType))
	err = Cli_ErrorText(code)
	return
}
func (c *S7Client) Connect() (err error) {
	var code C.int = C.Cli_Connect(c.client)
	err = Cli_ErrorText(code)
	return
}
func (c *S7Client) Disconnect() (err error) {
	var code C.int = C.Cli_Disconnect(c.client)
	err = Cli_ErrorText(code)
	return
}

/*
   ParamNumber 为P_u16_LocalPort的时候 value的数据是uint16 其他情况类似的
*/
func (c *S7Client) GetParam(paraNumber ParamNumber) (value interface{}, err error) {
	var pValue unsafe.Pointer
	switch paraNumber {
	case P_u16_RemotePort:
		pValue = unsafe.Pointer(new(uint16))
	case P_i32_PingTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_SendTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_RecvTimeout:
		pValue = unsafe.Pointer(new(int32))
	case P_u16_SrcRef:
		pValue = unsafe.Pointer(new(uint16))
	case P_u16_DstRef:
		pValue = unsafe.Pointer(new(uint16))
	case P_u16_SrcTSap:
		pValue = unsafe.Pointer(new(uint16))
	case P_i32_PDURequest:
		pValue = unsafe.Pointer(new(int32))
	}
	var code C.int = C.Cli_GetParam(c.client, C.int(paraNumber), pValue)
	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	switch paraNumber {
	case P_u16_RemotePort:
		value = *(*uint16)(pValue)
	case P_i32_PingTimeout:
		value = *(*int32)(pValue)
	case P_i32_SendTimeout:
		value = *(*int32)(pValue)
	case P_i32_RecvTimeout:
		value = *(*int32)(pValue)
	case P_u16_SrcRef:
		value = *(*uint16)(pValue)
	case P_u16_DstRef:
		value = *(*uint16)(pValue)
	case P_u16_SrcTSap:
		value = *(*uint16)(pValue)
	case P_i32_PDURequest:
		value = *(*int32)(pValue)
	}
	return
}

/*
   P_u16_LocalPort 设定端口为uint16
*/
func (c *S7Client) SetParam(paraNumber ParamNumber, value interface{}) (err error) {
	pvalue := Value_Pvalue(paraNumber, value)
	var code C.int = C.Cli_SetParam(c.client, C.int(paraNumber), pvalue)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ReadArea(S7Object Client, int Area, int DBNumber, int Start, int Amount, int WordLen, void *pUsrData);
func (c *S7Client) ReadArea(area S7Area, dBNumber int32, start int32, amount int32, wordLen S7WL) (pUsrData []byte, err error) {
	pUsrData = make([]byte, dataLength(wordLen, amount))
	var code C.int = C.Cli_ReadArea(c.client, C.int(area), C.int(dBNumber), C.int(start), C.int(amount), C.int(wordLen), unsafe.Pointer(&pUsrData[0]))
	err = Cli_ErrorText(code)
	return
}
func (c *S7Client) checkWriteAmount(pUsrData []byte, wordLen S7WL) (amount int32, err error) {
	if len(pUsrData)%int(wordLen.size()) != 0 {
		err = fmt.Errorf("length of pUserData != wordLen size * amount")
		return
	}
	amount = int32(len(pUsrData)) / wordLen.size()
	return
}
func (c *S7Client) WriteArea(area S7Area, dBNumber int32, start int32, wordLen S7WL, pUsrData []byte) (err error) {
	amount, err := c.checkWriteAmount(pUsrData, wordLen)
	if err != nil {
		return
	}
	var code C.int = C.Cli_WriteArea(c.client, C.int(area), C.int(dBNumber), C.int(start), C.int(amount), C.int(wordLen), unsafe.Pointer(&pUsrData[0]))
	err = Cli_ErrorText(code)
	return
}
func (c *S7Client) freeItemsC(items []TS7DataItem) {
	for k := range items {
		C.free(unsafe.Pointer(items[k].Pdata))
	}
}
func (c *S7Client) ReadMultiVars(items []TS7DataItemGo) (err error) {
	itemsCount := len(items)
	itemsC := make([]TS7DataItem, itemsCount)
	for k := range items {
		itemsC[k] = items[k].ToC()
	}
	for k := range itemsC {
		itemsC[k].Pdata = (*byte)(C.malloc(
			C.size_t(dataLength(S7WL(itemsC[k].WordLen), itemsC[k].Amount)),
		))
	}
	defer c.freeItemsC(itemsC)
	var code C.int = C.Cli_ReadMultiVars(
		c.client,
		(C.PS7DataItem)(unsafe.Pointer(&itemsC[0])),
		C.int(itemsCount),
	)

	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	for k, v := range itemsC {
		items[k].Result = v.Result
		items[k].Pdata = v.GetBytes()
	}
	return
}
func (c *S7Client) WriteMultiVars(items []TS7DataItemGo) (err error) {
	itemsCount := len(items)
	itemsC := make([]TS7DataItem, itemsCount)
	for k, _ := range items {
		itemsC[k] = items[k].ToC()
		items[k].CopyPdata(&itemsC[k])
	}
	defer c.freeItemsC(itemsC)
	var code C.int = C.Cli_WriteMultiVars(
		c.client,
		(C.PS7DataItem)(unsafe.Pointer(&itemsC[0])),
		C.int(itemsCount))

	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	for k, v := range itemsC {
		items[k].Result = v.Result
	}
	return
}

func (c *S7Client) DBRead(dBNumber int32, start int32, size int32) (pUsrData []byte, err error) {
	return c.ReadArea(S7AreaDB, dBNumber, start, size, S7WLByte)
}
func (c *S7Client) DBWrite(dBNumber int32, start int32, pUsrData []byte) (err error) {
	return c.WriteArea(S7AreaDB, dBNumber, start, S7WLByte, pUsrData)
}
func (c *S7Client) MBRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.ReadArea(S7AreaMK, 0, start, size, S7WLByte)
}
func (c *S7Client) MBWrite(start int32, pUsrData []byte) (err error) {
	return c.WriteArea(S7AreaMK, 0, start, S7WLByte, pUsrData)
}
func (c *S7Client) EBRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.ReadArea(S7AreaPE, 0, start, size, S7WLByte)
}
func (c *S7Client) EBWrite(start int32, pUsrData []byte) (err error) {
	return c.WriteArea(S7AreaPE, 0, start, S7WLByte, pUsrData)
}
func (c *S7Client) ABRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.ReadArea(S7AreaPA, 0, start, size, S7WLByte)
}
func (c *S7Client) ABWrite(start int32, pUsrData []byte) (err error) {
	return c.WriteArea(S7AreaPA, 0, start, S7WLByte, pUsrData)
}
func (c *S7Client) TMRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.ReadArea(S7AreaTM, 0, start, size, S7WLTimer)
}
func (c *S7Client) TMWrite(start int32, pUsrData []byte) (err error) {
	return c.WriteArea(S7AreaTM, 0, start, S7WLTimer, pUsrData)
}
func (c *S7Client) CTRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.ReadArea(S7AreaCT, 0, start, size, S7WLCounter)
}
func (c *S7Client) CTWrite(start int32, pUsrData []byte) (err error) {
	return c.WriteArea(S7AreaCT, 0, start, S7WLCounter, pUsrData)
}
func (c *S7Client) ListBlocks() (pUsrData TS7BlocksList, err error) {
	var code C.int = C.Cli_ListBlocks(c.client, (*C.TS7BlocksList)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}
func (c *S7Client) GetAgBlockInfo(blockType Block, blockNum int32) (pUsrData TS7BlockInfo, err error) {
	var code C.int = C.Cli_GetAgBlockInfo(c.client, C.int(blockType), C.int(blockNum), (*C.TS7BlockInfo)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetPgBlockInfo(S7Object Client, void *pBlock, TS7BlockInfo *pUsrData, int Size);
func (c *S7Client) GetPgBlockInfo(pBlock []byte) (pUsrData TS7BlockInfo, err error) {
	size := len(pBlock)
	var code C.int = C.Cli_GetPgBlockInfo(c.client, unsafe.Pointer(&pBlock[0]), (*C.TS7BlockInfo)(unsafe.Pointer(&pUsrData)), C.int(size))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ListBlocksOfType(S7Object Client, int BlockType, TS7BlocksOfType *pUsrData, int *ItemsCount);
func (c *S7Client) ListBlocksOfType(blockType Block, cap int32) (pUsrData []TS7BlocksOfType, err error) {
	pUsrData = make([]TS7BlocksOfType, cap)
	var code C.int = C.Cli_ListBlocksOfType(
		c.client,
		C.int(blockType),
		(*C.TS7BlocksOfType)(unsafe.Pointer(&pUsrData[0])),
		(*C.int)(unsafe.Pointer(&cap)))
	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	pUsrData = pUsrData[:cap]
	return
}

//int S7API Cli_Upload(S7Object Client, int BlockType, int BlockNum, void *pUsrData, int *Size);
func (c *S7Client) Upload(blockType Block, blockNum int32, pUsrData []byte) (size int32, err error) {
	size = int32(len(pUsrData))
	var code C.int = C.Cli_Upload(c.client, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_FullUpload(S7Object Client, int BlockType, int BlockNum, void *pUsrData, int *Size);
func (c *S7Client) FullUpload(blockType Block, blockNum int32, pUsrData []byte) (size int32, err error) {
	size = int32(len(pUsrData))
	var code C.int = C.Cli_FullUpload(c.client, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_Download(S7Object Client, int BlockNum, void *pUsrData, int Size);
func (c *S7Client) Download(blockNum int32, size int32) (pUsrData []byte, err error) {
	pUsrData = make([]byte, size)
	var code C.int = C.Cli_Download(c.client, C.int(blockNum), unsafe.Pointer(&pUsrData[0]), C.int(size))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_Delete(S7Object Client, int BlockType, int BlockNum);
func (c *S7Client) Delete(blockType Block, blockNum int32) (err error) {
	var code C.int = C.Cli_Delete(c.client, C.int(blockType), C.int(blockNum))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_DBGet(S7Object Client, int DBNumber, void *pUsrData, int *Size);
func (c *S7Client) DBGet(dBNumber int32, pUsrData []byte) (size int32, err error) {
	return c.Upload(Block_DB, dBNumber, pUsrData)
}

//int S7API Cli_DBFill(S7Object Client, int DBNumber, int FillChar);
func (c *S7Client) DBFill(dBNumber int32, fillChar int32) (err error) {
	var code C.int = C.Cli_DBFill(c.client, C.int(dBNumber), C.int(fillChar))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetPlcDateTime(S7Object Client, tm *DateTime);
func (c *S7Client) GetPlcDateTime() (dataTime Tm, err error) {
	var code C.int = C.Cli_GetPlcDateTime(c.client, (*C.tm)(unsafe.Pointer(&dataTime)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_SetPlcDateTime(S7Object Client, tm *DateTime);
func (c *S7Client) SetPlcDateTime(dataTime Tm) (err error) {
	var code C.int = C.Cli_SetPlcDateTime(c.client, (*C.tm)(unsafe.Pointer(&dataTime)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_SetPlcSystemDateTime(S7Object Client);
func (c *S7Client) SetPlcSystemDateTime() (err error) {
	var code C.int = C.Cli_SetPlcSystemDateTime(c.client)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetOrderCode(S7Object Client, TS7OrderCode *pUsrData);
func (c *S7Client) GetOrderCode() (pUsrData TS7OrderCode, err error) {
	var code C.int = C.Cli_GetOrderCode(c.client, (*C.TS7OrderCode)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetCpuInfo(S7Object Client, TS7CpuInfo *pUsrData);
func (c *S7Client) GetCpuInfo() (pUsrData TS7CpuInfo, err error) {
	var code C.int = C.Cli_GetCpuInfo(c.client, (*C.TS7CpuInfo)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetCpInfo(S7Object Client, TS7CpInfo *pUsrData);
func (c *S7Client) GetCpInfo() (pUsrData TS7CpInfo, err error) {
	var code C.int = C.Cli_GetCpInfo(c.client, (*C.TS7CpInfo)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ReadSZL(S7Object Client, int ID, int Index, TS7SZL *pUsrData, int *Size);
func (c *S7Client) ReadSZL(id int32, index int32) (pUsrData TS7SZL, size int32, err error) {
	var code C.int = C.Cli_ReadSZL(c.client, C.int(id), C.int(index), (*C.TS7SZL)(unsafe.Pointer(&pUsrData)), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ReadSZLList(S7Object Client, TS7SZLList *pUsrData, int *ItemsCount);
func (c *S7Client) ReadSZLList(capacity int32) (ret []TS7SZLList, err error) {
	var itemsCount = capacity
	ret = make([]TS7SZLList, capacity)
	var code C.int = C.Cli_ReadSZLList(c.client, (*C.TS7SZLList)(unsafe.Pointer(&ret[0])), (*C.int)(unsafe.Pointer(&capacity)))
	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	ret = ret[:itemsCount]
	return
}

//int S7API Cli_PlcHotStart(S7Object Client);
func (c *S7Client) PlcHotStart() (err error) {
	var code C.int = C.Cli_PlcHotStart(c.client)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_PlcColdStart(S7Object Client);
func (c *S7Client) PlcColdStart() (err error) {
	var code C.int = C.Cli_PlcColdStart(c.client)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_PlcStop(S7Object Client);
func (c *S7Client) PlcStop() (err error) {
	var code C.int = C.Cli_PlcStop(c.client)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_CopyRamToRom(S7Object Client, int Timeout);
func (c *S7Client) CopyRamToRom(timeout int32) (err error) {
	var code C.int = C.Cli_CopyRamToRom(c.client, C.int(timeout))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_Compress(S7Object Client, int Timeout);
func (c *S7Client) Compress(timeout int32) (err error) {
	var code C.int = C.Cli_Compress(c.client, C.int(timeout))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetPlcStatus(S7Object Client, int *Status);
func (c *S7Client) GetPlcStatus() (status S7CpuStatus, err error) {
	var code C.int = C.Cli_GetPlcStatus(c.client, (*C.int)(unsafe.Pointer(&status)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetProtection(S7Object Client, TS7Protection *pUsrData);
func (c *S7Client) GetProtection() (pUsrData TS7Protection, err error) {
	var code C.int = C.Cli_GetProtection(c.client, (*C.TS7Protection)(unsafe.Pointer(&pUsrData)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_SetSessionPassword(S7Object Client, char *Password);
func (c *S7Client) SetSessionPassword(Password string) (err error) {
	password := C.CString(Password)
	defer func() {
		C.free(unsafe.Pointer(password))
	}()
	var code C.int = C.Cli_SetSessionPassword(c.client, password)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ClearSessionPassword(S7Object Client);
func (c *S7Client) ClearSessionPassword() (err error) {
	var code C.int = C.Cli_ClearSessionPassword(c.client)
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_IsoExchangeBuffer(S7Object Client, void *pUsrData, int *Size);
func (c *S7Client) IsoExchangeBuffer(pUsrData []byte) (size int32, err error) {
	size = int32(len(pUsrData))
	var code C.int = C.Cli_IsoExchangeBuffer(c.client, unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetExecTime(S7Object Client, int *Time);
func (c *S7Client) GetExecTime() (time int32, err error) {
	var code C.int = C.Cli_GetExecTime(c.client, (*C.int)(unsafe.Pointer(&time)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetLastError(S7Object Client, int *LastError);
func (c *S7Client) GetLastError() (lastError CliErrorCode, err error) {
	var code C.int = C.Cli_GetLastError(c.client, (*C.int)(unsafe.Pointer(&lastError)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_GetPduLength(S7Object Client, int *Requested, int *Negotiated);
func (c *S7Client) GetPduLength() (requested int32, negotiated int32, err error) {
	var code C.int = C.Cli_GetPduLength(c.client, (*C.int)(unsafe.Pointer(&requested)), (*C.int)(unsafe.Pointer(&negotiated)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_ErrorText(int Error, char *Text, int TextLen);

//int S7API Cli_GetConnected(S7Object Client, int *Connected);
func (c *S7Client) GetConnected() (connected int32, err error) {
	var code C.int = C.Cli_GetConnected(c.client, (*C.int)(unsafe.Pointer(&connected)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_AsReadArea(S7Object Client, int Area, int DBNumber, int Start, int Amount, int WordLen, void *pUsrData);
func (c *S7Client) AsReadArea(area S7Area, dBNumber int32, start int32, amount int32, wordLen S7WL) (pUsrData []byte, err error) {
	pUsrData = make([]byte, dataLength(wordLen, amount))
	var code C.int = C.Cli_AsReadArea(c.client, C.int(area), C.int(dBNumber), C.int(start), C.int(amount), C.int(wordLen), unsafe.Pointer(&pUsrData[0]))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsWriteArea(S7Object Client, int Area, int DBNumber, int Start, int Amount, int WordLen, void *pUsrData);
func (c *S7Client) AsWriteArea(area S7Area, dBNumber int32, start int32, wordLen S7WL, pUsrData []byte) (err error) {
	if len(pUsrData)%int(wordLen.size()) != 0 {
		err = fmt.Errorf("length of pUserData != wordLen size * amount")
		return
	}
	amount := int32(len(pUsrData)) / wordLen.size()
	var code C.int = C.Cli_AsWriteArea(c.client, C.int(area), C.int(dBNumber), C.int(start), C.int(amount), C.int(wordLen), unsafe.Pointer(&pUsrData[0]))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsDBRead(S7Object Client, int DBNumber, int Start, int Size, void *pUsrData);
func (c *S7Client) AsDBRead(dBNumber int32, start int32, size int32) (pUsrData []byte, err error) {
	return c.AsReadArea(S7AreaDB, dBNumber, start, size, S7WLByte)
}

// int S7API Cli_AsDBWrite(S7Object Client, int DBNumber, int Start, int Size, void *pUsrData);
func (c *S7Client) AsDBWrite(dBNumber int32, start int32, pUsrData []byte) (err error) {
	return c.AsWriteArea(S7AreaDB, dBNumber, start, S7WLByte, pUsrData)
}

// int S7API Cli_AsMBRead(S7Object Client, int Start, int Size, void *pUsrData);
func (c *S7Client) AsMBRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.AsReadArea(S7AreaMK, 0, start, size, S7WLByte)
}

// int S7API Cli_AsMBWrite(S7Object Client, int Start, int Size, void *pUsrData);
func (c *S7Client) AsMBWrite(start int32, pUsrData []byte) (err error) {
	return c.AsWriteArea(S7AreaMK, 0, start, S7WLByte, pUsrData)
}

// int S7API Cli_AsEBRead(S7Object Client, int Start, int Size, void *pUsrData);
func (c *S7Client) AsEBRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.AsReadArea(S7AreaPE, 0, start, size, S7WLByte)
}

// int S7API Cli_AsEBWrite(S7Object Client, int Start, int Size, void *pUsrData);
func (c *S7Client) AsEBWrite(start int32, pUsrData []byte) (err error) {
	return c.AsWriteArea(S7AreaPE, 0, start, S7WLByte, pUsrData)
}

// int S7API Cli_AsABRead(S7Object Client, int Start, int Size, void *pUsrData);
func (c *S7Client) AsABRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.AsReadArea(S7AreaPA, 0, start, size, S7WLByte)
}

// int S7API Cli_AsABWrite(S7Object Client, int Start, int Size, void *pUsrData);
func (c *S7Client) AsABWrite(start int32, pUsrData []byte) (err error) {
	return c.AsWriteArea(S7AreaPA, 0, start, S7WLByte, pUsrData)
}

// int S7API Cli_AsTMRead(S7Object Client, int Start, int Amount, void *pUsrData);
func (c *S7Client) AsTMRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.AsReadArea(S7AreaTM, 0, start, size, S7WLTimer)
}

// int S7API Cli_AsTMWrite(S7Object Client, int Start, int Amount, void *pUsrData);
func (c *S7Client) AsTMWrite(start int32, pUsrData []byte) (err error) {
	return c.AsWriteArea(S7AreaTM, 0, start, S7WLTimer, pUsrData)
}

// int S7API Cli_AsCTRead(S7Object Client, int Start, int Amount, void *pUsrData);
func (c *S7Client) AsCTRead(start int32, size int32) (pUsrData []byte, err error) {
	return c.AsReadArea(S7AreaCT, 0, start, size, S7WLCounter)
}

// int S7API Cli_AsCTWrite(S7Object Client, int Start, int Amount, void *pUsrData);
func (c *S7Client) AsCTWrite(start int32, pUsrData []byte) (err error) {
	return c.AsWriteArea(S7AreaCT, 0, start, S7WLCounter, pUsrData)
}

// int S7API Cli_AsListBlocksOfType(S7Object Client, int BlockType, TS7BlocksOfType *pUsrData, int *ItemsCount);
func (c *S7Client) AsListBlocksOfType(blockType Block, cap int32) (pUsrData []TS7BlocksOfType, err error) {
	pUsrData = make([]TS7BlocksOfType, cap)
	var code C.int = C.Cli_AsListBlocksOfType(
		c.client,
		C.int(blockType),
		(*C.TS7BlocksOfType)(unsafe.Pointer(&pUsrData[0])),
		(*C.int)(unsafe.Pointer(&cap)))
	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	pUsrData = pUsrData[:cap]
	return
}

// int S7API Cli_AsReadSZL(S7Object Client, int ID, int Index, TS7SZL *pUsrData, int *Size);
func (c *S7Client) AsReadSZL(id int32, index int32) (pUsrData TS7SZL, size int32, err error) {
	var code C.int = C.Cli_AsReadSZL(c.client, C.int(id), C.int(index), (*C.TS7SZL)(unsafe.Pointer(&pUsrData)), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsReadSZLList(S7Object Client, TS7SZLList *pUsrData, int *ItemsCount);
func (c *S7Client) AsReadSZLList(capacity int32) (ret []TS7SZLList, err error) {
	var itemsCount = capacity
	ret = make([]TS7SZLList, capacity)
	var code C.int = C.Cli_AsReadSZLList(c.client, (*C.TS7SZLList)(unsafe.Pointer(&ret[0])), (*C.int)(unsafe.Pointer(&capacity)))
	err = Cli_ErrorText(code)
	if err != nil {
		return
	}
	ret = ret[:itemsCount]
	return
}

// int S7API Cli_AsUpload(S7Object Client, int BlockType, int BlockNum, void *pUsrData, int *Size);
func (c *S7Client) AsUpload(blockType Block, blockNum int32, pUsrData []byte) (size int32, err error) {
	size = int32(len(pUsrData))
	var code C.int = C.Cli_AsUpload(c.client, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsFullUpload(S7Object Client, int BlockType, int BlockNum, void *pUsrData, int *Size);
func (c *S7Client) AsFullUpload(blockType Block, blockNum int32, pUsrData []byte) (size int32, err error) {
	size = int32(len(pUsrData))
	var code C.int = C.Cli_AsFullUpload(c.client, C.int(blockType), C.int(blockNum), unsafe.Pointer(&pUsrData[0]), (*C.int)(unsafe.Pointer(&size)))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsDownload(S7Object Client, int BlockNum, void *pUsrData, int Size);
func (c *S7Client) AsDownload(blockNum int32, size int32) (pUsrData []byte, err error) {
	pUsrData = make([]byte, size)
	var code C.int = C.Cli_AsDownload(c.client, C.int(blockNum), unsafe.Pointer(&pUsrData[0]), C.int(size))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsCopyRamToRom(S7Object Client, int Timeout);
func (c *S7Client) AsCopyRamToRom(timeout int32) (err error) {
	var code C.int = C.Cli_AsCopyRamToRom(c.client, C.int(timeout))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsCompress(S7Object Client, int Timeout);
func (c *S7Client) AsCompress(timeout int32) (err error) {
	var code C.int = C.Cli_AsCompress(c.client, C.int(timeout))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_AsDBGet(S7Object Client, int DBNumber, void *pUsrData, int *Size);
func (c *S7Client) AsDBGet(dBNumber int32, pUsrData []byte) (size int32, err error) {
	return c.AsUpload(Block_DB, dBNumber, pUsrData)
}

// int S7API Cli_AsDBFill(S7Object Client, int DBNumber, int FillChar);
func (c *S7Client) AsDBFill(dBNumber int32, fillChar int32) (err error) {
	var code C.int = C.Cli_AsDBFill(c.client, C.int(dBNumber), C.int(fillChar))
	err = Cli_ErrorText(code)
	return
}

// int S7API Cli_CheckAsCompletion(S7Object Client, int *opResult);
func (c *S7Client) CheckAsCompletion() (opResult JobStatus, err error) {
	var code C.int = C.Cli_CheckAsCompletion(c.client, (*C.int)(unsafe.Pointer(&opResult)))
	err = Cli_ErrorText(code)
	return
}

//int S7API Cli_WaitAsCompletion(S7Object Client, int Timeout);
func (c *S7Client) WaitAsCompletion(timeout int32) (err error) {
	var code C.int = C.Cli_WaitAsCompletion(c.client, C.int(timeout))
	err = Cli_ErrorText(code)
	return
}
