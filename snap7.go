package snap7go

import (
	"fmt"
	"unsafe"
)

//------------------------------------------------------------------------------
//                                  PARAMS LIST
//------------------------------------------------------------------------------

type ParamNumber = int

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
type JobStatus int

const JobComplete JobStatus = 0
const JobPending JobStatus = 1
const errLibInvalidObject JobStatus = -2

//******************************************************************************
//                                   CLIENT
//******************************************************************************
// Error codes
type ErrorCode int

const errNegotiatingPDU ErrorCode = 0x00100000
const errCliInvalidParams ErrorCode = 0x00200000
const errCliJobPending ErrorCode = 0x00300000
const errCliTooManyItems ErrorCode = 0x00400000
const errCliInvalidWordLen ErrorCode = 0x00500000
const errCliPartialDataWritten ErrorCode = 0x00600000
const errCliSizeOverPDU ErrorCode = 0x00700000
const errCliInvalidPlcAnswer ErrorCode = 0x00800000
const errCliAddressOutOfRange ErrorCode = 0x00900000
const errCliInvalidTransportSize ErrorCode = 0x00A00000
const errCliWriteDataSizeMismatch ErrorCode = 0x00B00000
const errCliItemNotAvailable ErrorCode = 0x00C00000
const errCliInvalidValue ErrorCode = 0x00D00000
const errCliCannotStartPLC ErrorCode = 0x00E00000
const errCliAlreadyRun ErrorCode = 0x00F00000
const errCliCannotStopPLC ErrorCode = 0x01000000
const errCliCannotCopyRamToRom ErrorCode = 0x01100000
const errCliCannotCompress ErrorCode = 0x01200000
const errCliAlreadyStop ErrorCode = 0x01300000
const errCliFunNotAvailable ErrorCode = 0x01400000
const errCliUploadSequenceFailed ErrorCode = 0x01500000
const errCliInvalidDataSizeRecvd ErrorCode = 0x01600000
const errCliInvalidBlockType ErrorCode = 0x01700000
const errCliInvalidBlockNumber ErrorCode = 0x01800000
const errCliInvalidBlockSize ErrorCode = 0x01900000
const errCliDownloadSequenceFailed ErrorCode = 0x01A00000
const errCliInsertRefused ErrorCode = 0x01B00000
const errCliDeleteRefused ErrorCode = 0x01C00000
const errCliNeedPassword ErrorCode = 0x01D00000
const errCliInvalidPassword ErrorCode = 0x01E00000
const errCliNoPasswordToSetOrClear ErrorCode = 0x01F00000
const errCliJobTimeout ErrorCode = 0x02000000
const errCliPartialDataRead ErrorCode = 0x02100000
const errCliBufferTooSmall ErrorCode = 0x02200000
const errCliFunctionRefused ErrorCode = 0x02300000
const errCliDestroying ErrorCode = 0x02400000
const errCliInvalidParamNumber ErrorCode = 0x02500000
const errCliCannotChangeParam ErrorCode = 0x02600000

// Client Connection Type
type CONNTYPE = uint16

const CONNTYPE_PG CONNTYPE = 0x0001    // Connect to the PLC as a PG
const CONNTYPE_OP CONNTYPE = 0x0002    // Connect to the PLC as an OP
const CONNTYPE_BASIC CONNTYPE = 0x0003 // Basic connection

// Area ID
type S7Area = int

const S7AreaPE S7Area = 0x81
const S7AreaPA S7Area = 0x82
const S7AreaMK S7Area = 0x83
const S7AreaDB S7Area = 0x84
const S7AreaCT S7Area = 0x1C
const S7AreaTM S7Area = 0x1D

// Word Length
type S7WL int

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

func dataLength(wordLen S7WL, amount int32, start int32) int32 {
	return wordLen.size()*amount + start
}
func (s S7WL) size() int32 {
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
		return 1
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

type MaskKind = int

//MaskKind  Srv_GetMask/Srv_SetMask
const (
	MaskKindEvent = 0
	MaskKindLog   = 1
)

//ServerStatus
type S7ServerStatus int

const (
	SrvStopped S7ServerStatus = 0 //The Server is stopped.
	SrvRunning S7ServerStatus = 1 //The Server is Running.
	SrvError   S7ServerStatus = 2 //Server Error.
)

// CPU status
type S7CpuStatus int

const S7CpuStatusUnknown S7CpuStatus = 0x00
const S7CpuStatusRun S7CpuStatus = 0x08
const S7CpuStatusStop S7CpuStatus = 0x04
