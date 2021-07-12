package snap7go

/*

	用c语言来调用go语言，只能是全局的函数
	1.注册的时候所有的回调函数都注册同一个函数GlobalAsCallback
	2.把注册时候的usrptr作为key存入registeredCallBacks
	3.GlobalAsCallback的时候根据usrptr的时候动态指派

	注意，再为每个client注册回调函数时候，务必使usrptr不同!!

	usrptr在这里充当了上下文的作用，caller去确定需要什么样的上下文。

*/

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
/*
	// real go implementation
	extern void S7API GlobalAsCallback(void *usrPtr, int opCode, int opResult);

*/
import "C"
import "unsafe"

var registeredCallBacks = make(map[uintptr]func(usrptr uintptr, opCode int, opResult int))

/*

 */
func Cli_SetAsCallback(
	client S7Object,
	callback func(usrptr uintptr, opCode int, opResult int),
	usrptr uintptr) error {
	var code C.int = C.Cli_SetAsCallback(
		client,
		/*
			todo
			为啥golang认为函数指针是这个类型
		*/
		(*[0]byte)(unsafe.Pointer(C.GlobalAsCallback)),
		unsafe.Pointer(usrptr))
	err := cliErrorsTable[int(code)]
	if err != nil {
		return err
	}
	registeredCallBacks[usrptr] = callback
	return err
}

//export GlobalAsCallback
func GlobalAsCallback(usrptr *C.void, opCode C.int, opResult C.int) {
	callback := registeredCallBacks[uintptr(unsafe.Pointer(usrptr))]
	if callback == nil {
		return
	}
	callback(uintptr(unsafe.Pointer(usrptr)), int(opCode), int(opResult))
}
