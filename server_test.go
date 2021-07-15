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

	defer func() {
		err = serverDefault.Stop()
		ast.Nil(err)
		serverDefault.Destroy()
	}()

	clientDefault := NewS7Client()
	defer clientDefault.Destroy()
	//连接默认地址(127.0.0.1)
	err = clientDefault.Connect()
	ast.Nil(err)

	/*
		指定地址的server
	*/
	serverDesignated := NewS7Server()
	serverDesignated.SetEventsCallback(justPrintEvent)
	serverDesignated.SetReadEventsCallback(justPrintEvent)
	//SetParam须在start之前，否则errSrvCannotChangeParam
	err=serverDesignated.SetParam(P_i32_MaxClients,int32(12))
	ast.Nil(err)
	err = serverDesignated.StartTo("127.0.0.1")
	ast.Nil(err)

	defer func() {
		err = serverDesignated.Stop()
		ast.Nil(err)
		serverDesignated.Destroy()
	}()

	clientDesignated := NewS7Client()
	defer clientDesignated.Destroy()
	//连接指定地址
	err = clientDesignated.ConnectTo("127.0.0.1",0,2)
	ast.Nil(err)

	getValue,err:=serverDesignated.GetParam(P_i32_MaxClients)
	ast.Nil(err)
	ast.Equal(int32(12),getValue)
}
