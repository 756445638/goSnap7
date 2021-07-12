package snap7go

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

// Area ID
const byte S7AreaPE = 0x81
const byte S7AreaPA = 0x82
const byte S7AreaMK = 0x83
const byte S7AreaDB = 0x84
const byte S7AreaCT = 0x1C
const byte S7AreaTM = 0x1D

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
