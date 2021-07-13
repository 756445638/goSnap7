package snap7go

import "C"

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

func Cli_ErrorText(code int) error {
	if code == 0 {
		return nil
	}
	const length = 512
	var buf [length]byte
	var errCode = C.Cli_ErrorText(C.int(code), (*C.char)(unsafe.Pointer(&buf[0])), length)
	if errCode != 0 {
		if e, ok := cliErrorsTable[int(errCode)]; ok {
			return e
		} else {
			return fmt.Errorf("unknown error code %d", errCode)
		}
	}
	return errors.New(string(buf[:]))
}

func Srv_ErrorText(code int) error {
	if code == 0 {
		return nil
	}
	const length = 512
	var buf [length]byte
	var errCode = C.Srv_ErrorText(C.int(code), (*C.char)(unsafe.Pointer(&buf[0])), length)
	if errCode != 0 {
		if e, ok := serverErrorsTable[int(errCode)]; ok {
			return e
		} else {
			return fmt.Errorf("unknown error code %d", errCode)
		}
	}
	return errors.New(string(buf[:]))
}
