package snap7go

//#cgo CFLAGS: -I .
//#include "snap7.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type S7Object = C.S7Object

//------------------------------------------------------------------------------
//                                  PARAMS LIST
//------------------------------------------------------------------------------

type ParamNumber int32

const P_u16_LocalPort ParamNumber = 1
const P_u16_RemotePort ParamNumber = 2
const P_i32_PingTimeout ParamNumber = 3
const P_i32_SendTimeout ParamNumber = 4
const P_i32_RecvTimeout ParamNumber = 5
const P_i32_WorkInterval ParamNumber = 6
const P_u16_SrcRef ParamNumber = 7
const P_u16_DstRef ParamNumber = 8
const P_u16_SrcTSap ParamNumber = 9
const P_i32_PDURequest ParamNumber = 10
const P_i32_MaxClients ParamNumber = 11
const P_i32_BSendTimeout ParamNumber = 12
const P_i32_BRecvTimeout ParamNumber = 13
const P_u32_RecoveryTime ParamNumber = 14
const P_u32_KeepAliveTime ParamNumber = 15

// Client/Partner Job status
type JobStatus int32

const JobComplete JobStatus = 0
const JobPending JobStatus = 1
const errLibInvalidObject JobStatus = -2

//******************************************************************************
//                                   CLIENT
//******************************************************************************
// Error codes
type CliErrorCode int32

const errNegotiatingPDU CliErrorCode = 0x00100000
const errCliInvalidParams CliErrorCode = 0x00200000
const errCliJobPending CliErrorCode = 0x00300000
const errCliTooManyItems CliErrorCode = 0x00400000
const errCliInvalidWordLen CliErrorCode = 0x00500000
const errCliPartialDataWritten CliErrorCode = 0x00600000
const errCliSizeOverPDU CliErrorCode = 0x00700000
const errCliInvalidPlcAnswer CliErrorCode = 0x00800000
const errCliAddressOutOfRange CliErrorCode = 0x00900000
const errCliInvalidTransportSize CliErrorCode = 0x00A00000
const errCliWriteDataSizeMismatch CliErrorCode = 0x00B00000
const errCliItemNotAvailable CliErrorCode = 0x00C00000
const errCliInvalidValue CliErrorCode = 0x00D00000
const errCliCannotStartPLC CliErrorCode = 0x00E00000
const errCliAlreadyRun CliErrorCode = 0x00F00000
const errCliCannotStopPLC CliErrorCode = 0x01000000
const errCliCannotCopyRamToRom CliErrorCode = 0x01100000
const errCliCannotCompress CliErrorCode = 0x01200000
const errCliAlreadyStop CliErrorCode = 0x01300000
const errCliFunNotAvailable CliErrorCode = 0x01400000
const errCliUploadSequenceFailed CliErrorCode = 0x01500000
const errCliInvalidDataSizeRecvd CliErrorCode = 0x01600000
const errCliInvalidBlockType CliErrorCode = 0x01700000
const errCliInvalidBlockNumber CliErrorCode = 0x01800000
const errCliInvalidBlockSize CliErrorCode = 0x01900000
const errCliDownloadSequenceFailed CliErrorCode = 0x01A00000
const errCliInsertRefused CliErrorCode = 0x01B00000
const errCliDeleteRefused CliErrorCode = 0x01C00000
const errCliNeedPassword CliErrorCode = 0x01D00000
const errCliInvalidPassword CliErrorCode = 0x01E00000
const errCliNoPasswordToSetOrClear CliErrorCode = 0x01F00000
const errCliJobTimeout CliErrorCode = 0x02000000
const errCliPartialDataRead CliErrorCode = 0x02100000
const errCliBufferTooSmall CliErrorCode = 0x02200000
const errCliFunctionRefused CliErrorCode = 0x02300000
const errCliDestroying CliErrorCode = 0x02400000
const errCliInvalidParamNumber CliErrorCode = 0x02500000
const errCliCannotChangeParam CliErrorCode = 0x02600000

// Client Connection Type
type CONNTYPE uint16

const CONNTYPE_PG CONNTYPE = 0x0001    // Connect to the PLC as a PG
const CONNTYPE_OP CONNTYPE = 0x0002    // Connect to the PLC as an OP
const CONNTYPE_BASIC CONNTYPE = 0x0003 // Basic connection

// Area ID
type S7Area int32

const S7AreaPE S7Area = 0x81
const S7AreaPA S7Area = 0x82
const S7AreaMK S7Area = 0x83
const S7AreaDB S7Area = 0x84
const S7AreaCT S7Area = 0x1C
const S7AreaTM S7Area = 0x1D

// Word Length
type S7WL int32

const S7WLBit S7WL = 0x01
const S7WLByte S7WL = 0x02
const S7WLWord S7WL = 0x04
const S7WLDWord S7WL = 0x06
const S7WLReal S7WL = 0x08
const S7WLCounter S7WL = 0x1C
const S7WLTimer S7WL = 0x1D

// Block type
type Block int32

const Block_OB Block = 0x38
const Block_DB Block = 0x41
const Block_SDB Block = 0x42
const Block_FC Block = 0x43
const Block_SFC Block = 0x44
const Block_FB Block = 0x45
const Block_SFB Block = 0x46

// Sub Block Type
// const byte SubBlk_OB  = 0x08;
// const byte SubBlk_DB  = 0x0A;
// const byte SubBlk_SDB = 0x0B;
// const byte SubBlk_FC  = 0x0C;
// const byte SubBlk_SFC = 0x0D;
// const byte SubBlk_FB  = 0x0E;
// const byte SubBlk_SFB = 0x0F;

// Block languages
// const byte BlockLangAWL       = 0x01;
// const byte BlockLangKOP       = 0x02;
// const byte BlockLangFUP       = 0x03;
// const byte BlockLangSCL       = 0x04;
// const byte BlockLangDB        = 0x05;
// const byte BlockLangGRAPH     = 0x06;

/*
	字节数组的长度
*/
func DataLength(wordLen S7WL, amount int32) int32 {
	t := wordLen.Size() * amount
	return t
}

func (s S7WL) Size() int32 {
	switch s {
	case S7WLBit:
		return 1
	case S7WLByte:
		return 1
	case S7WLWord:
		return 2
	case S7WLDWord:
		return 4
	case S7WLReal:
		return 4
	case S7WLCounter:
		return 2
	case S7WLTimer:
		return 2
	}
	panic(fmt.Sprintf("S7WL not exist:%d", s))
}

func Value_Pvalue(paraNumber ParamNumber, value interface{}) (pValue unsafe.Pointer) {
	switch paraNumber {
	case P_u16_LocalPort:
		t := new(uint16)
		*t = value.(uint16)
		pValue = unsafe.Pointer(t)
	case P_u16_RemotePort:
		t := new(uint16)
		*t = value.(uint16)
		pValue = unsafe.Pointer(t)
	case P_i32_PingTimeout:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_i32_SendTimeout:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_i32_RecvTimeout:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_i32_WorkInterval:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_u16_SrcRef:
		t := new(uint16)
		*t = value.(uint16)
		pValue = unsafe.Pointer(t)
	case P_u16_DstRef:
		t := new(uint16)
		*t = value.(uint16)
		pValue = unsafe.Pointer(t)
	case P_u16_SrcTSap:
		t := new(uint16)
		*t = value.(uint16)
		pValue = unsafe.Pointer(t)
	case P_i32_PDURequest:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_i32_MaxClients:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_i32_BSendTimeout:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_i32_BRecvTimeout:
		t := new(int32)
		*t = value.(int32)
		pValue = unsafe.Pointer(t)
	case P_u32_RecoveryTime:
		t := new(uint32)
		*t = value.(uint32)
		pValue = unsafe.Pointer(t)
	case P_u32_KeepAliveTime:
		t := new(uint32)
		*t = value.(uint32)
		pValue = unsafe.Pointer(t)
	}
	return
}

//******************************************************************************
//                                   SERVER
//******************************************************************************
type Operation int32

const OperationRead Operation = 0
const OperationWrite Operation = 1

type MaskKind int32

//MaskKind  Srv_GetMask/Srv_SetMask
const (
	MaskKindEvent MaskKind = 0
	MaskKindLog   MaskKind = 1
)

//ServerStatus
type S7ServerStatus int32

const (
	SrvStopped S7ServerStatus = 0 //The Server is stopped.
	SrvRunning S7ServerStatus = 1 //The Server is Running.
	SrvError   S7ServerStatus = 2 //Server Error.
)

// CPU status
type S7CpuStatus int32

const S7CpuStatusUnknown S7CpuStatus = 0x00
const S7CpuStatusRun S7CpuStatus = 0x08
const S7CpuStatusStop S7CpuStatus = 0x04

type SrvAreaType int32

const SrvAreaPE SrvAreaType = 0
const SrvAreaPA SrvAreaType = 1
const SrvAreaMK SrvAreaType = 2
const SrvAreaCT SrvAreaType = 3
const SrvAreaTM SrvAreaType = 4
const SrvAreaDB SrvAreaType = 5

type SrvErrCode int32

const errSrvCannotStart SrvErrCode = 0x00100000        // Server cannot start
const errSrvDBNullPointer SrvErrCode = 0x00200000      // Passed null as PData
const errSrvAreaAlreadyExists SrvErrCode = 0x00300000  // Area Re-registration
const errSrvUnknownArea SrvErrCode = 0x00400000        // Unknown area
const errSrvInvalidParams SrvErrCode = 0x00500000      // Invalid param(s) supplied
const errSrvTooManyDB SrvErrCode = 0x00600000          // Cannot register DB
const errSrvInvalidParamNumber SrvErrCode = 0x00700000 // Invalid param (srv_get/set_param)
const errSrvCannotChangeParam SrvErrCode = 0x00800000  // Cannot change because running

func justPrintEvent(e *TSrvEvent) {
	s, err := Srv_EventText(e)
	if err != nil {
		fmt.Printf("Srv_EventText() failed,code:%d err:%v\n", e.EvtCode, err)
		return
	}
	fmt.Println(s)
}
