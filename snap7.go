package snap7go

import "fmt"

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

//******************************************************************************
//                                   CLIENT
//******************************************************************************
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

func (s S7WL) Size() int {
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
	panic(fmt.Sprintf("S7WL not exist:", s))
}
