package snap7go

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
		return fmt.Errorf(" C.Cli_ErrorText failed code :%d,origin code %d", errCode, code)
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
		return fmt.Errorf(" C.Srv_ErrorText failed code :%d,origin code %d", errCode, code)
	}
	return errors.New(string(buf[:]))
}
