package snap7go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

/*
Srv_Create Creates a Server Object.
Srv_Destroy Destroys a Server Object.
Srv_StartTo Starts a Server Object onto a given IP Address.
Srv_Start Starts a Server Object onto the default adapter.
Srv_Stop Stops the Server.
Srv_GetParam Reads an internal Server parameter.
Srv_SetParam Writes an internal Server Parameter
*/
func TestServerAdministrative(t *testing.T) {
	ast := assert.New(t)
	/*
		默认地址（127.0.0.1）的server
	*/
	serverDefault := NewS7Server()
	serverDefault.SetEventsCallback(justPrintEvent)
	serverDefault.SetReadEventsCallback(justPrintEvent)

	err := serverDefault.Start()
	ast.Nil(err)

	//后面会测试Stop()和Destroy()，所以不用defer
	//defer func() {
	//	err = serverDefault.Stop()
	//	ast.Nil(err)
	//	serverDefault.Destroy()
	//}()

	clientDefault := NewS7Client()
	defer clientDefault.Destroy()
	//连接默认地址(127.0.0.1)
	err = clientDefault.Connect()
	ast.Nil(err)

	err = serverDefault.Stop()
	ast.Nil(err)
	serverDefault.Destroy()

	/*
		指定地址的server
	*/
	serverDesignated := NewS7Server()
	serverDesignated.SetEventsCallback(justPrintEvent)
	serverDesignated.SetReadEventsCallback(justPrintEvent)
	//SetParam须在start之前，否则errSrvCannotChangeParam
	err = serverDesignated.SetParam(P_i32_MaxClients, int32(12))
	ast.Nil(err)
	err = serverDesignated.StartTo("127.0.0.1")
	ast.Nil(err)

	clientDesignated := NewS7Client()
	defer clientDesignated.Destroy()
	//连接指定地址
	err = clientDesignated.ConnectTo("127.0.0.1", 0, 2)
	ast.Nil(err)

	getValue, err := serverDesignated.GetParam(P_i32_MaxClients)
	ast.Nil(err)
	ast.Equal(int32(12), getValue)

	//Stop后client无法连接
	err = serverDesignated.Stop()
	ast.Nil(err)
	err = clientDesignated.ConnectTo("127.0.0.1", 0, 2)
	ast.NotNil(err)
	//重新startTo后，能够连接
	err = serverDesignated.StartTo("127.0.0.1")
	ast.Nil(err)
	err = clientDesignated.ConnectTo("127.0.0.1", 0, 2)
	ast.Nil(err)

	//Destroy后无法startTo
	err = serverDesignated.Stop()
	ast.Nil(err)
	serverDesignated.Destroy()
	err = serverDesignated.StartTo("127.0.0.1")
	ast.NotNil(err)
}

/*
Srv_RegisterArea Shares a given memory area with the server.
Srv_UnRegisterArea “Unshares” a memory area previously shared.
Srv_LockArea Locks a shared memory area.
Srv_UnlockArea Unlocks a previously locked shared memory area.
*/
func TestServerSharedMemory(t *testing.T) {
	ast := assert.New(t)

	server := NewS7Server()
	server.SetEventsCallback(justPrintEvent)
	server.SetReadEventsCallback(justPrintEvent)

	/*
		RegisterArea(),服务器共享内存区域
	*/

	data := [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	err := server.RegisterArea(SrvAreaPA, 0, data[:])
	ast.Nil(err)

	err = server.Start()
	ast.Nil(err)

	defer func() {
		err = server.Stop()
		ast.Nil(err)
		server.Destroy()
	}()

	client := NewS7Client()
	defer client.Destroy()
	err = client.Connect()
	ast.Nil(err)

	//WordLen=S7WLBit,一次只读一个byte，即amount必须等于1
	dataReadS7WLBit, err := client.ReadArea(S7AreaPA, 0, 0, 1, S7WLBit)
	ast.Nil(err)
	//data的byte为1时，dataReadS7WLBit=[]byte{1},否则dataReadS7WLBit=[]byte{0}
	ast.Equal([]byte{1}, dataReadS7WLBit)
	dataReadS7WLBit2, err := client.ReadArea(S7AreaPA, 0, 2, 1, S7WLBit)
	ast.Nil(err)
	ast.Equal([]byte{0}, dataReadS7WLBit2)

	//从S7AreaPA的第二个byte开始读，读取长度为3
	dataRead, err := client.ReadArea(S7AreaPA, 0, 1, 3, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{2, 3, 4}, dataRead)
	//不能读取不同Area的共享信息
	_, err = client.ReadArea(S7AreaPE, 0, 1, 3, S7WLByte)
	ast.NotNil(err)

	/*
		UnregisterArea(),取消服务器共享内存区域
	*/
	err = server.UnregisterArea(SrvAreaPA, 0)
	ast.Nil(err)
	//取消共享后，无法读取S7AreaPA区域
	_, err = client.ReadArea(S7AreaPA, 0, 1, 3, S7WLByte)
	ast.NotNil(err)

	/*
		LockArea(),锁定区域
	*/
	dataBeforeLockArea := [10]byte{3, 2, 1}
	err = server.RegisterArea(SrvAreaPE, 0, dataBeforeLockArea[:])
	ast.Nil(err)
	//锁定前能读
	dataReadBeforeLockArea, err := client.ReadArea(S7AreaPE, 0, 1, 3, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{2, 1, 0}, dataReadBeforeLockArea)
	//锁定区域
	err = server.LockArea(SrvAreaPE, 0)
	ast.Nil(err)
	//锁定区域后，client无法读取该区域的信息
	_, err = client.ReadArea(S7AreaPE, 0, 1, 3, S7WLByte)
	ast.NotNil(err)

	/*
		UnlockArea(),解锁区域
	*/
	err = server.UnlockArea(SrvAreaPE, 0)
	ast.Nil(err)
	//解锁区域后，client能够再次读取该区域的信息
	dataReadAfterUnlockArea, err := client.ReadArea(S7AreaPE, 0, 1, 3, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{2, 1, 0}, dataReadAfterUnlockArea)
}

/*
Srv_SetEventsCallback Sets the user callback that the Server object has to call when an event is created.
Srv_SetRWAreaCallback Sets the user callback that the Server object has to call on a read or write request.
Srv_GetMask Reads the specified filter mask.
Srv_SetMask Writes the specified filter mask.
Srv_PickEvent Extracts an event (if available) from the Events queue.
Srv_ClearEvents Empties the Event queue.
*/
func TestServerControlFlow(t *testing.T) {
	ast := assert.New(t)

	server := NewS7Server()
	err := server.SetEventsCallback(justPrintEvent)
	ast.Nil(err)
	err = server.SetReadEventsCallback(justPrintEvent)
	ast.Nil(err)

	/*
		SetMask(),设置掩码
		GetMask(),获取掩码
	*/
	err = server.SetMask(MaskKindEvent, uint32(123))
	ast.Nil(err)
	getMask, err := server.GetMask(MaskKindEvent)
	ast.Nil(err)
	ast.Equal(uint32(123), getMask)
	//不同 MaskKind
	getMaskKindLog, err := server.GetMask(MaskKindLog)
	ast.Nil(err)
	ast.NotEqual(uint32(123), getMaskKindLog)

	/*
		PickEvent(),提取事件
	*/
	//无事件
	pEventNil, err := server.PickEvent()
	ast.Nil(pEventNil)
	ast.Nil(err)

	//开启服务器，存在事件
	err = server.Start()
	ast.Nil(err)
	defer func() {
		err = server.Stop()
		ast.Nil(err)
		server.Destroy()
	}()
	pEvent, err := server.PickEvent()
	ast.Nil(err)
	ast.NotNil(pEvent)

	/*
		ClearEvents(),清空事件
	*/
	err = server.ClearEvents()
	ast.Nil(err)
	//清空事件后，提取事件为nil
	pEventAfterClearEvents, err := server.PickEvent()
	ast.Nil(err)
	ast.Nil(pEventAfterClearEvents)
}

/*
Srv_GetStatus Returns the last job execution time in milliseconds.
Srv_SetCpuStatus Returns the last job result.
Srv_EventText Returns a textual explanation of a given event.
Srv_ErrorText Returns a textual explanation of a given error number
*/
func TestServerMiscellaneous(t *testing.T) {
	ast := assert.New(t)

	server := NewS7Server()
	server.SetEventsCallback(justPrintEvent)
	server.SetReadEventsCallback(justPrintEvent)

	/*
		GetStatus(),Reads the server status, the Virtual CPU status and the number of the clients connected.
	*/
	serverStatus, cpuStatus, clientsCount, err := server.GetStatus()
	ast.Nil(err)
	ast.Equal(SrvStopped, serverStatus)
	ast.Equal(S7CpuStatusRun, cpuStatus)
	ast.Equal(0, clientsCount)

	server.Start()
	defer func() {
		err = server.Stop()
		ast.Nil(err)
		server.Destroy()
	}()
	serverStatusAfterStart, _, _, err := server.GetStatus()
	ast.Equal(SrvRunning, serverStatusAfterStart)

	//一个客户端连接server后，clientsCount+=1
	client1 := NewS7Client()
	defer client1.Destroy()
	client1.Connect()
	_, _, clientsCount, err = server.GetStatus()
	ast.Equal(1, clientsCount)
	client2 := NewS7Client()
	defer client2.Destroy()
	client2.Connect()
	_, _, clientsCount, err = server.GetStatus()
	ast.Equal(2, clientsCount)

	/*
		SetCpuStatus(),设置CPU状态
	*/
	err = server.SetCpuStatus(S7CpuStatusStop)
	ast.Nil(err)
	_, cpuStatusAfterSet, _, err := server.GetStatus()
	ast.Nil(err)
	ast.Equal(S7CpuStatusStop, cpuStatusAfterSet)

	/*
		Srv_EventText(),
	*/
	//存在Server start事件
	pEvent, err := server.PickEvent()
	ast.Nil(err)
	evenText, err := Srv_EventText(pEvent)
	ast.Nil(err)
	ast.True(strings.Contains(evenText, "Server started"))

	/*
		ErrorText(),
	*/
	//code在1到8之间
	err = Srv_ErrorText(0x00100000)
	ast.NotNil(err)
	ast.True(strings.Contains(fmt.Sprintf("%s", err), "Server cannot start"))
	//超出1到8的范围，err 为 Unknown error
	err = Srv_ErrorText(0x00900000)
	ast.NotNil(err)
	ast.True(strings.Contains(fmt.Sprintf("%s", err), "Unknown error"))
}
