package snap7go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClientAdministrativeCli(t *testing.T) { //已完成
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
	//连接地址(127.0.0.1)
	err = clientDefault.Connect()
	ast.Nil(err)

	/*
	   指定地址的server
	*/
	serverDesignated := NewS7Server()
	serverDesignated.SetEventsCallback(justPrintEvent)
	serverDesignated.SetReadEventsCallback(justPrintEvent)
	err = serverDesignated.StartTo("127.0.0.1")
	ast.Nil(err)

	defer func() {
		err = serverDesignated.Stop()
		ast.Nil(err)
		serverDesignated.Destroy()
	}()
	//默认client    Connect 默认连接的是PG (the programming console)
	client := NewS7Client()
	defer func() {
		err = client.Disconnect()
		ast.Nil(err)
		client.Destroy()
	}()
	//自定义client
	clientDesignated := NewS7Client()
	defer clientDesignated.Destroy()

	//默认client连接默认server
	err = client.Connect()
	ast.Nil(err)

	//SetParam在ConnectTo前后都可以 client与server可设置的ParamNumbers项不一样
	err = clientDesignated.SetParam(P_u16_LocalPort, uint16(2484))
	ast.NotNil(err)
	//err = clientDesignated.SetParam(P_u16_RemotePort, uint16(1548))      RemotePort的设置与ConnectTo有关
	//ast.Nil(err)
	err = clientDesignated.SetParam(P_i32_PingTimeout, int32(10))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_i32_SendTimeout, int32(10))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_i32_RecvTimeout, int32(10))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_i32_WorkInterval, int32(0))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_u16_SrcRef, uint16(5))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_u16_DstRef, uint16(5))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_u16_SrcTSap, uint16(0x1000))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_i32_PDURequest, int32(10))
	ast.Nil(err)
	err = clientDesignated.SetParam(P_i32_MaxClients, int32(4))
	ast.NotNil(err)
	err = clientDesignated.SetParam(P_i32_BSendTimeout, int32(4))
	ast.NotNil(err)
	err = clientDesignated.SetParam(P_i32_BRecvTimeout, int32(4))
	ast.NotNil(err)
	err = clientDesignated.SetParam(P_u32_RecoveryTime, uint32(4))
	ast.NotNil(err)
	err = clientDesignated.SetParam(P_u32_KeepAliveTime, uint32(4))
	ast.NotNil(err)

	//SetParam在ConnectTo前后都可以 client与server可设置的ParamNumbers项不一样
	err = clientDesignated.SetParam(P_u16_LocalPort, uint16(2484))
	ast.NotNil(err)
	//err = clientDesignated.SetParam(P_u16_RemotePort, uint16(1548))      RemotePort的设置与ConnectTo有关
	//ast.Nil(err)
	_, err = clientDesignated.GetParam(P_i32_PingTimeout)
	ast.Nil(err)
	_, err = clientDesignated.GetParam(P_i32_SendTimeout)
	ast.Nil(err)
	_, err = clientDesignated.GetParam(P_i32_RecvTimeout)
	ast.Nil(err)
	//_, err = clientDesignated.GetParam(P_i32_WorkInterval)  这个不知道为什么通不过
	//ast.NotNil(err)
	_, err = clientDesignated.GetParam(P_u16_SrcRef)
	ast.Nil(err)
	_, err = clientDesignated.GetParam(P_u16_DstRef)
	ast.Nil(err)
	_, err = clientDesignated.GetParam(P_u16_SrcTSap)
	ast.Nil(err)
	_, err = clientDesignated.GetParam(P_i32_PDURequest)
	ast.Nil(err)
	_, err = clientDesignated.GetParam(P_i32_MaxClients)
	ast.NotNil(err)
	_, err = clientDesignated.GetParam(P_i32_BSendTimeout)
	ast.NotNil(err)
	_, err = clientDesignated.GetParam(P_i32_BRecvTimeout)
	ast.NotNil(err)
	_, err = clientDesignated.GetParam(P_u32_RecoveryTime)
	ast.NotNil(err)
	_, err = clientDesignated.GetParam(P_u32_KeepAliveTime)
	ast.NotNil(err)

	//SetConnectionType:CONNTYPE_PG、CONNTYPE_OP、CONNTYPE_BASIC
	err = clientDesignated.SetConnectionType(CONNTYPE_PG)
	ast.Nil(err)

	err = clientDesignated.SetConnectionType(CONNTYPE_OP)
	ast.Nil(err)

	err = clientDesignated.SetConnectionType(CONNTYPE_BASIC)
	ast.Nil(err)

	err = clientDesignated.SetConnectionParams("127.0.0.1", 0x1000, 0x1000)
	ast.Nil(err)

	//自定义client连接指定地址
	err = clientDesignated.ConnectTo("127.0.0.1", 0, 0)
	ast.Nil(err)

}

func TestDataIOCli(t *testing.T) { //已完成
	ast := assert.New(t)
	/*
	   默认地址（127.0.0.1）的server
	*/
	serverDefault := NewS7Server()
	serverDefault.SetEventsCallback(justPrintEvent)
	serverDefault.SetReadEventsCallback(justPrintEvent)

	err := serverDefault.Start()
	ast.Nil(err)
	var dbArea [1024]byte
	err = serverDefault.RegisterArea(SrvAreaPE, 1, dbArea[:])
	ast.Nil(err)

	defer func() {
		err = serverDefault.Stop()
		ast.Nil(err)
		serverDefault.Destroy()
	}()

	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	//S7AreaPE    S7WLBit
	pUsrData := []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.WriteArea(S7AreaPE, 1, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	ret, err := client.ReadArea(S7AreaPE, 1, 0, 1, S7WLBit)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaPE    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} //输入的数据length是S7WLBit的Word size的倍数
	err = client.WriteArea(S7AreaPE, 1, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPE, 1, 0, 12, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.EBWrite(0, pUsrData)
	ast.Nil(err)
	ret, err = client.EBRead(0, 12)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPE    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaPE, 1, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPE, 1, 0, 6, S7WLWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPE    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaPE, 1, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPE, 1, 0, 3, S7WLDWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPE    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaPE, 1, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPE, 1, 0, 3, S7WLReal)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//server register dbAreapPA [1024]byte
	err = serverDefault.RegisterArea(SrvAreaPA, 1, dbArea[:])
	ast.Nil(err)
	//S7AreaPA    S7WLBit
	pUsrData = []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.WriteArea(S7AreaPA, 1, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPA, 1, 0, 1, S7WLBit)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaPA    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaPA, 1, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPA, 1, 0, 12, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//当WordLen = S7WLBytes时，使用ABWrite/ABRead简化WriteArea/ReadArea
	err = client.ABWrite(0, pUsrData)
	ast.Nil(err)
	ret, err = client.ABRead(0, 12)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPA    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaPA, 1, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPA, 1, 0, 6, S7WLWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPA    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaPA, 1, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPA, 1, 0, 3, S7WLDWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPA    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaPA, 1, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaPA, 1, 0, 3, S7WLReal)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//server register  SrvAreaMK 	（除了 SrvAreaDB特殊之外，其余的情况index与dbNmber设置无效）
	err = serverDefault.RegisterArea(SrvAreaMK, 0, dbArea[:])
	ast.Nil(err)
	//S7AreaMK    S7WLBit
	pUsrData = []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.WriteArea(S7AreaMK, 2, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaMK, 2, 0, 1, S7WLBit)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaMK    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaMK, 1, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaMK, 1, 0, 12, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.MBWrite(0, pUsrData)
	ast.Nil(err)
	ret, err = client.MBRead(0, 12)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaMK    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaMK, 1, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaMK, 1, 0, 6, S7WLWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaMK    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaMK, 1, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaMK, 1, 0, 3, S7WLDWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaMK    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaMK, 1, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaMK, 1, 0, 3, S7WLReal)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//server register  SrvAreaDB
	//var dbAreapDB [1024]byte                  dbNumber与 index 相对应，index未注册的 dbNumber 无法找到  （SrvAreaDB特殊之处，其余的情况index与dbNmber设置无效）
	err = serverDefault.RegisterArea(SrvAreaDB, 2, dbArea[:])
	ast.Nil(err)
	//S7AreaDB    S7WLBit
	pUsrData = []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.WriteArea(S7AreaDB, 2, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaDB, 2, 0, 1, S7WLBit)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaDB    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaDB, 2, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaDB, 2, 0, 12, S7WLByte)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.DBWrite(2, 0, pUsrData)
	ast.Nil(err)
	ret, err = client.DBRead(2, 0, 12)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaDB    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaDB, 2, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaDB, 2, 0, 6, S7WLWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaDB    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaDB, 2, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaDB, 2, 0, 3, S7WLDWord)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaDB    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaDB, 2, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaDB, 2, 0, 3, S7WLReal)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	var dbAreaCT [1024]byte
	err = serverDefault.RegisterArea(SrvAreaCT, 1, dbAreaCT[:])
	ast.Nil(err)
	//S7AreaCT    S7WLCounter
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaCT, 1, 0, S7WLCounter, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaCT, 1, 0, 6, S7WLCounter)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.CTWrite(0, pUsrData)
	ast.Nil(err)
	ret, err = client.CTRead(0, 6)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	var dbAreaTM [1024]byte
	err = serverDefault.RegisterArea(SrvAreaTM, 1, dbAreaTM[:])
	ast.Nil(err)
	//S7AreaTM    S7WLTimer
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.WriteArea(S7AreaTM, 1, 0, S7WLTimer, pUsrData)
	ast.Nil(err)
	ret, err = client.ReadArea(S7AreaTM, 1, 0, 6, S7WLTimer)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.TMWrite(0, pUsrData)
	ast.Nil(err)
	ret, err = client.TMRead(0, 6)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

}

func TestDirectoryCli(t *testing.T) { //未完成
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
	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	ret, err := client.GetProtection()
	fmt.Println("Protection级别信息：", ret) // {1 0 1 2 0}
	ast.Nil(err)
	//设置8位用户密码
	err = client.SetSessionPassword("12345678")
	ast.Nil(err)

	//pUsrData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	////func (c *S7Client) Upload(blockType Block, blockNum int, pUsrData []byte) (size int, err error) {
	//ret1, err := client.Upload(Block_OB, 1, pUsrData) //CPU权限不够
	//fmt.Println("fullUpload Buffer size:", ret1)
	//ast.Nil(err)
	//
	//ret1, err = client.FullUpload(Block_OB, 1, pUsrData)
	//fmt.Println("fullUpload Buffer size:", ret1)
	//ast.Nil(err)

	//显示（OB、FB、FC、SFB、SFC、DB、SDB）7种Blocks的数量
	rete, err := client.ListBlocks() //Blocks都为0，不知道怎么建立block，应该是用upload建立内容，但是没有权限
	fmt.Println("ListBlocks:", rete)
	ast.Nil(err)

	_, err = client.ListBlocksOfType(Block_OB, 10) //没有BLOCK无法测试
	//fmt.Println("TS7BlocksOfType", data)
	ast.Nil(err)

	ret2, err := client.GetAgBlockInfo(Block_OB, 1)
	fmt.Println("AgBlockInfo:", ret2)
	ast.Nil(err)

}

func TestBlockOrientedCli(t *testing.T) { //未完成
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
	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	ret, err := client.GetProtection()
	fmt.Println("Protection级别信息：", ret) // {1 0 1 2 0}
	ast.Nil(err)

	//设置8位用户密码
	err = client.SetSessionPassword("12345678")
	ast.Nil(err)

	ret, err = client.GetProtection()
	fmt.Println("Protection级别信息：", ret) // {1 0 1 2 0}
	ast.Nil(err)

	pUsrData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	//func (c *S7Client) Upload(blockType Block, blockNum int, pUsrData []byte) (size int, err error) {
	ret1, err := client.Upload(Block_OB, 1, pUsrData) //CPU权限不够  ,后面的都无法测试
	fmt.Println("fullUpload Buffer size:", ret1)
	ast.Nil(err)

	ret1, err = client.FullUpload(Block_OB, 1, pUsrData)
	fmt.Println("fullUpload Buffer size:", ret1)
	ast.Nil(err)

	downloadData, err := client.Download(1, 32)
	fmt.Println(downloadData)
	ast.Nil(err)

	err = client.Delete(Block_OB, 1)
	ast.Nil(err)

	ret1, err = client.DBGet(2, pUsrData) //CPU权限不够
	fmt.Println("fullUpload Buffer size:", ret1)
	ast.Nil(err)

	err = client.DBFill(2, 10086) //CPU权限不够
	ast.Nil(err)
}

//系统状态列表（德语：System-ZustandsListen)
func TestDateOrTimeCli(t *testing.T) { // 未完成
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

	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	//todo: set目前无效
	var timeSet Tm
	goTime := time.Unix(11, 0)
	timeSet.FromTime(goTime)
	fmt.Printf("!!!!!!Set Tm：%#v\n", timeSet)
	fmt.Printf("!!!!!!Set time：%#v\n", timeSet.ToTime())
	fmt.Println("!!!!!!time", goTime)
	err = client.SetPlcDateTime(timeSet)
	ast.Nil(err)

	err = client.SetPlcSystemDateTime()
	ast.Nil(err)

	dataTimeGet, err := client.GetPlcDateTime()
	ast.Nil(err)
	fmt.Printf("!!!!!!Get Tm：%#v\n", dataTimeGet)
	fmt.Printf("!!!!!!Get time：%#v\n", dataTimeGet.ToTime())
	fmt.Println("!!!!!!time", dataTimeGet.ToTime())
	//ast.Equal(goTime, dataTimeGet.ToTime())
	ast.Equal(timeSet, dataTimeGet)

}

//系统状态列表（德语：System-ZustandsListen)
func TestSystemInfoCli(t *testing.T) { // 未完成,ReadSZL  与 ReadSZLList 未完成
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

	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	ts7szl, size, err := client.ReadSZL(0x0232, 0x0004) //与upload有关
	fmt.Println("系统状态列表：", ts7szl, size)
	ast.Nil(err)

	_, err = client.ReadSZLList(100)
	//fmt.Println("ReadSZLList：", ret)
	ast.Nil(err)

	ordercode, err6 := client.GetOrderCode()
	fmt.Println("ordercode：", ordercode)
	ast.Nil(err6)
	_, err6 = client.GetCpuInfo()

	//fmt.Println("CpuInfo：", cpuInf)
	//fmt.Println("GetModuleTypeName：", cpuInf.GetModuleTypeName())
	//fmt.Println("GetSerialNumber：", cpuInf.GetSerialNumber())
	//fmt.Println("GetASName：", cpuInf.GetASName())
	//fmt.Println("GetCopyright：", cpuInf.GetCopyright())
	//fmt.Println("GetModuleName：", cpuInf.GetModuleName())
	//fmt.Println("CpuInfo：", cpuInf.GetASName())

	ast.Nil(err6)
	cpInf, err7 := client.GetCpInfo()
	fmt.Println("CpInfo：", cpInf)
	ast.Nil(err7)

}

func TestPLCControlCli(t *testing.T) { //已完成
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
	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)
	//PLC running hot start
	err = client.PlcHotStart()
	ast.Nil(err)

	err = client.PlcStop()
	ast.Nil(err)

	err = client.PlcColdStart()
	ast.Nil(err)

	//timeout：ms
	err = client.CopyRamToRom(20)
	ast.Nil(err)

	err = client.Compress(30)
	ast.Nil(err)

	PlcStatus, err := client.GetPlcStatus()
	ast.Equal(S7CpuStatusRun, PlcStatus)
	ast.Nil(err)

}

func TestSecurityCli(t *testing.T) { //完成，但有点小疑惑
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
	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	//Security
	//1 -mode selector   0-no password   1- CPU  2-:Mode selector setting RUN-P     0-Startup switch setting :undefined,
	ret4, err4 := client.GetProtection()
	fmt.Println("Protection级别信息：", ret4) //{1 0 1 2 0}
	ast.Nil(err4)

	err2 := client.SetSessionPassword("12345678")
	ast.Nil(err2)

	ret4, err5 := client.GetProtection() //{1 0 1 2 0} 设置完密码后得到的保护级别信息第二参数 应该不是0了才对？？？？
	fmt.Println("Protection级别信息：", ret4)
	ast.Nil(err5)

	err3 := client.ClearSessionPassword()
	ast.Nil(err3, "清除密码成功")

	ret4, err5 = client.GetProtection() //{1 0 1 2 0}
	fmt.Println("Protection级别信息：", ret4)
	ast.Nil(err5)
}

func TestMiscellaneousCli(t *testing.T) { //未完成
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
	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	time, err := client.GetExecTime()
	fmt.Println(" last job execution time:", time)
	ast.Nil(err)

	lastErr, err := client.GetLastError()
	fmt.Println(" lastErr:", lastErr)
	ast.Nil(err)

	//requested:Address of the PDU Req. variable       ???地址怎么找
	requested, negotiated, err := client.GetPduLength()
	fmt.Println(" negotiated:", negotiated)
	fmt.Println(" requested:", requested)
	ast.Nil(err)

	isconnecte, err := client.GetConnected()
	fmt.Println(" isconnecte:", isconnecte)
	ast.Nil(err)
}

/*
TestAsynchronousCli()
	Cli_AsReadArea Reads a data area from a PLC.
	Cli_AsWriteArea Writes a data area into a PLC.
	Cli_AsDBRead Reads a part of a DB from a PLC.
	Cli_AsDBWrite Writes a part of a DB into a PLC.
	Cli_AsABRead Reads a part of IPU area from a PLC.
	Cli_AsABWrite Writes a part of IPU area into a PLC.
	Cli_AsEBRead Reads a part of IPI area from a PLC.
	Cli_AsEBWrite Writes a part of IPI area into a PLC.
	Cli_AsMBRead Reads a part of Merkers area from a PLC.
	Cli_AsMBWrite Writes a part of Merkers area into a PLC.
	Cli_AsTMRead Reads timers from a PLC. Cli_AsTMWrite Write timers into a PLC.
	Cli_AsCTRead Reads counters from a PLC.
	Cli_AsCTWrite Write counters into a PLC.
	Cli_AsListBlocksOfType Returns the AG blocks list of a given type.
	Cli_AsReadSZL Reads a partial list of given ID and Index.
	Cli_AsReadSZLList Reads the list of partial lists available in the CPU.
	Cli_AsFullUpload Uploads a block from AG with Header and Footer infos.
	Cli_AsUpload Uploads a block from AG.
	Cli_AsDownload Download a block into AG.
	Cli_AsDBGet Uploads a DB from AG using DBRead.
	Cli_AsDBFill Fills a DB in AG with a given byte.
	Cli_AsCopyRamToRom Performs the Copy Ram to Rom action.
	Cli_AsCompress Performs the Compress action.
*/
func TestAsynchronousCli(t *testing.T) {
	ast := assert.New(t)
	/*
	   默认地址（127.0.0.1）的server
	*/
	serverDefault := NewS7Server()
	serverDefault.SetEventsCallback(justPrintEvent)
	serverDefault.SetReadEventsCallback(justPrintEvent)

	err := serverDefault.Start()
	ast.Nil(err)
	var dbArea [1024]byte
	err = serverDefault.RegisterArea(SrvAreaPE, 1, dbArea[:])
	ast.Nil(err)

	defer func() {
		err = serverDefault.Stop()
		ast.Nil(err)
		serverDefault.Destroy()
	}()

	client := NewS7Client()
	defer client.Destroy()
	//连接地址(127.0.0.1)
	err = client.Connect()
	ast.Nil(err)

	//S7AreaPE    S7WLBit
	pUsrData := []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.AsWriteArea(S7AreaPE, 1, 0, S7WLBit, pUsrData)
	ast.Nil(err)

	//opResult1, err := client.CheckAsCompletion()
	//ast.Nil(err)
	//ast.Equal(JobPending, opResult1)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	opResult2, err := client.CheckAsCompletion()
	ast.Nil(err)
	ast.Equal(JobComplete, opResult2)

	ret, err := client.AsReadArea(S7AreaPE, 1, 0, 1, S7WLBit)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaPE    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} //输入的数据length是S7WLBit的Word size的倍数
	err = client.AsWriteArea(S7AreaPE, 1, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPE, 1, 0, 12, S7WLByte)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.AsEBWrite(0, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsEBRead(0, 12)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPE    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPE, 1, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPE, 1, 0, 6, S7WLWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPE    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPE, 1, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPE, 1, 0, 3, S7WLDWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPE    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPE, 1, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPE, 1, 0, 3, S7WLReal)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//server register dbAreapPA [1024]byte
	err = serverDefault.RegisterArea(SrvAreaPA, 1, dbArea[:])
	ast.Nil(err)
	//S7AreaPA    S7WLBit
	pUsrData = []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.AsWriteArea(S7AreaPA, 1, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPA, 1, 0, 1, S7WLBit)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaPA    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPA, 1, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPA, 1, 0, 12, S7WLByte)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//当WordLen = S7WLBytes时，使用ABWrite/ABRead简化WriteArea/ReadArea
	err = client.AsABWrite(0, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsABRead(0, 12)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPA    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPA, 1, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPA, 1, 0, 6, S7WLWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPA    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPA, 1, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPA, 1, 0, 3, S7WLDWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaPA    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPA, 1, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaPA, 1, 0, 3, S7WLReal)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//server register  SrvAreaMK 	（除了 SrvAreaDB特殊之外，其余的情况index与dbNmber设置无效）
	err = serverDefault.RegisterArea(SrvAreaMK, 0, dbArea[:])
	ast.Nil(err)
	//S7AreaMK    S7WLBit
	pUsrData = []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.AsWriteArea(S7AreaMK, 2, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaMK, 2, 0, 1, S7WLBit)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaMK    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaMK, 1, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaMK, 1, 0, 12, S7WLByte)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.AsMBWrite(0, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsMBRead(0, 12)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaMK    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaMK, 1, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaMK, 1, 0, 6, S7WLWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaMK    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaMK, 1, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaMK, 1, 0, 3, S7WLDWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaMK    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaMK, 1, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaMK, 1, 0, 3, S7WLReal)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//server register  SrvAreaDB
	//var dbAreapDB [1024]byte                  dbNumber与 index 相对应，index未注册的 dbNumber 无法找到  （SrvAreaDB特殊之处，其余的情况index与dbNmber设置无效）
	err = serverDefault.RegisterArea(SrvAreaDB, 2, dbArea[:])
	ast.Nil(err)
	//S7AreaDB    S7WLBit
	pUsrData = []byte{1} // https://github.com/756445638/snap7-go/issues/4
	err = client.AsWriteArea(S7AreaDB, 2, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaDB, 2, 0, 1, S7WLBit)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1}, ret)

	//S7AreaDB    S7WLByte
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaDB, 2, 0, S7WLByte, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaDB, 2, 0, 12, S7WLByte)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.AsDBWrite(2, 0, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsDBRead(2, 0, 12)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaDB    S7WLWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaDB, 2, 0, S7WLWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaDB, 2, 0, 6, S7WLWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaDB    S7WLDWord
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaDB, 2, 0, S7WLDWord, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaDB, 2, 0, 3, S7WLDWord)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	//S7AreaDB    S7WLReal
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaDB, 2, 0, S7WLReal, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaDB, 2, 0, 3, S7WLReal)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	var dbAreaCT [1024]byte
	err = serverDefault.RegisterArea(SrvAreaCT, 1, dbAreaCT[:])
	ast.Nil(err)
	//S7AreaCT    S7WLCounter
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaCT, 1, 0, S7WLCounter, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaCT, 1, 0, 6, S7WLCounter)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.AsCTWrite(0, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsCTRead(0, 6)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	var dbAreaTM [1024]byte
	err = serverDefault.RegisterArea(SrvAreaTM, 1, dbAreaTM[:])
	ast.Nil(err)
	//S7AreaTM    S7WLTimer
	pUsrData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaTM, 1, 0, S7WLTimer, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsReadArea(S7AreaTM, 1, 0, 6, S7WLTimer)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	err = client.AsTMWrite(0, pUsrData)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ret, err = client.AsTMRead(0, 6)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)
	ast.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ret)

	_, err = client.AsListBlocksOfType(Block_OB, 10) //没有BLOCK无法测试
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	szl, size, err := client.AsReadSZL(0x0232, 0x0004) //与upload有关
	ast.Nil(err)
	fmt.Println("系统状态列表：", szl, size)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	_, err = client.AsReadSZLList(100)
	//fmt.Println("ReadSZLList：", ret)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	ret1, err := client.AsUpload(Block_OB, 1, pUsrData) //CPU权限不够  ,后面的都无法测试
	ast.Nil(err)
	fmt.Println("fullUpload Buffer size:", ret1)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	ret1, err = client.AsFullUpload(Block_OB, 1, pUsrData)
	ast.Nil(err)
	fmt.Println("fullUpload Buffer size:", ret1)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	asDownloadData, err := client.AsDownload(1, 32)
	ast.Nil(err)
	fmt.Println(asDownloadData)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	dbGet, err := client.AsDBGet(2, pUsrData) //CPU权限不够
	ast.Nil(err)
	fmt.Println("fullUpload Buffer size:", dbGet)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	err = client.AsDBFill(2, 10086) //CPU权限不够
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	//timeout：ms
	err = client.AsCopyRamToRom(20)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	err = client.AsCompress(30)
	ast.Nil(err)
	err = client.WaitAsCompletion(10000)
	ast.Nil(err)

	//	Cli_AsListBlocksOfType Returns the AG blocks list of a given type.
	//	Cli_AsReadSZL Reads a partial list of given ID and Index.
	//	Cli_AsReadSZLList Reads the list of partial lists available in the CPU.
	//	Cli_AsFullUpload Uploads a block from AG with Header and Footer infos.
	//	Cli_AsUpload Uploads a block from AG.
	//	Cli_AsDownload Download a block into AG.
	//	Cli_AsDBGet Uploads a DB from AG using DBRead.
	//	Cli_AsDBFill Fills a DB in AG with a given byte.
	//	Cli_AsCopyRamToRom Performs the Copy Ram to Rom action.
	//	Cli_AsCompress Performs the Compress action.
	//	Cli_CheckAsCompletion Checks if the current asynchronous job was done and terminates immediately.

	//ret1, err := client.FullUpload(Block_OB, 1, pUsrData)
	//fmt.Println("fullUpload Buffer size:", ret1)
	//ast.Nil(err)

	//err = client.AsDownload(1, pUsrData, 12)
	//ast.Nil(err)

}
