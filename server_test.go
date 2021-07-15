package snap7go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestServerSharedMemory(t *testing.T) {
	ast := assert.New(t)

	server := NewS7Server()
	server.SetEventsCallback(justPrintEvent)
	server.SetReadEventsCallback(justPrintEvent)

	/*
		RegisterArea(),服务器共享内存区域
	*/
	data := [10]byte{1, 2, 3}
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
	err = client.Connect()
	ast.Nil(err)
	//从S7AreaPA的第二个byte开始读，读取长度为3
	dataRead, err := client.ReadArea(S7AreaPA, 0, 1, 3, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{2, 3, 0}, dataRead)
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
	dataBeforeLockArea := [10]byte{3,2,1}
	err = server.RegisterArea(SrvAreaPE, 0, dataBeforeLockArea[:])
	ast.Nil(err)
	err = server.LockArea(SrvAreaPE, 0)
	ast.Nil(err)

}
