package snap7go

import "C"

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"
)

type S7Object = C.S7Object

func Cli_Create() (client S7Object) {
	client = C.Cli_Create()
	return
}
func Cli_Destroy(client S7Object) {
	C.Cli_Destroy((*C.S7Object)(unsafe.Pointer(&client)))
	return
}
func Cli_ConnectTo(Client S7Object, Address string, Rack int, Slot int) (err error) {
	s := C.CString(Address)
	defer func() {
		C.free(unsafe.Pointer(s))
	}()
	var code C.int = C.Cli_ConnectTo(Client, s, C.int(Rack), C.int(Slot))
	err = cliErrorsTable[int(code)]
	return
}

func Cli_GetParam(Client S7Object, paraNumber int) (interface{}, error) {

}
func Cli_GetCpuInfo(cli S7Object) (info TS7CpuInfo, err error) {
	var code C.int = C.Cli_GetCpuInfo(cli, (*C.TS7CpuInfo)(unsafe.Pointer(&info)))
	err = cliErrorsTable[int(code)]
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
