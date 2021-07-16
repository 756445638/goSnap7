package snap7go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

)

func TestClientAdministrative(t *testing.T) {
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
	   指定本地IP地址（192.168.187.1）的server
	*/
	serverDesignated := NewS7Server()
	serverDesignated.SetEventsCallback(justPrintEvent)
	serverDesignated.SetReadEventsCallback(justPrintEvent)
	err = serverDesignated.StartTo("192.168.187.1")
	ast.Nil(err)

	defer func() {
		err = serverDesignated.Stop()
		ast.Nil(err)
		serverDesignated.Destroy()
	}()

	clientDesignated := NewS7Client()

	defer clientDesignated.Destroy()
    //CONNTYPE_PG、CONNTYPE_OP、CONNTYPE_BASIC
	err = clientDesignated.SetConnectionType(CONNTYPE_BASIC)
	ast.Nil(err)
	err = clientDesignated.SetConnectionParams("192.168.187.1",0x1000, 0x1000)
	ast.Nil(err)

	//在ConnectTo前后都可以
	err=clientDesignated.SetParam(P_i32_SendTimeout,int32(4))
	ast.Nil(err)
	//连接指定地址(192.168.187.1)
	err = clientDesignated.ConnectTo("192.168.187.1",0,1)
	ast.Nil(err)

	paradata,err := clientDesignated.GetParam(P_i32_SendTimeout)
	fmt.Println(paradata)
	ast.Nil(err)
	//ast.Equal(11,paradata)

}


func TestDataIO(t *testing.T) {
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

	err = serverDefault.RegisterArea(SrvAreaPE , 1 ,  dbArea[:])
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
	//WriteArea
	//Invalid Transport size            Address out of range
	//ReadArea      S7AreaPE        S7WLBit             /S7WLCounter/S7WLTimer/
	//ReadArea      S7AreaPA       S7WLBit           /S7WLCounter/S7WLTimer/
	pUsrData := []byte{1}
	err = client.WriteArea(S7AreaPE, 1, 0, S7WLBit, pUsrData)
	ast.Nil(err)
	ret,err := client.ReadArea(S7AreaPE, 1, 0,1 , S7WLBit)
	fmt.Println("S7WLBit",ret)
	ast.Nil(err)



}





func TestDirectory(t *testing.T) {
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
	ret4,err4 := client.GetProtection()
	fmt.Println("Protection级别信息：",ret4)
	ast.Nil(err4)

	err2 := client.SetSessionPassword("12345678")
	ast.Nil(err2,"设置密码成功")

	ret4,err5 := client.GetProtection()
	fmt.Println("Protection级别信息：",ret4)
	ast.Nil(err5)

	err3 := client.ClearSessionPassword()
	ast.Nil(err3,"清除密码成功")
	pUsrData1:= []byte{1,2,3,4,5,6}
	size,err5 := client.IsoExchangeBuffer(pUsrData1)
	fmt.Println(size)
	ast.Nil(err5)



	pUsrData:= []byte{1,2,3,4,5,6,7,8}
	//func (c *S7Client) Upload(blockType Block, blockNum int, pUsrData []byte) (size int, err error) {

	ret,err := client.FullUpload(Block_OB, 1, pUsrData)
	fmt.Println("fullUpload Buffer size:",ret)
	ast.Nil(err)


	rete,err := client.ListBlocks()
	fmt.Println("ListBlocks:",rete)
	ast.Nil(err)

	data,itemCounter, err:= client.ListBlocksOfType(Block_OB)
	fmt.Println("TS7BlocksOfType",data)
	fmt.Println(itemCounter)
	ast.Nil(err)

	ret1,err := client.GetAgBlockInfo(Block_OB, 1)
	fmt.Println("AgBlockInfo:",ret1)
	ast.Nil(err)


	//fmt.Println(ret)
	//ast.Nil(err)


}



//系统状态列表（德语：System-ZustandsListen)
func TestSystemInfo(t *testing.T) {
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

	ts7szl,size,err2 := client.ReadSZL(0x0011, 0x0000)    //与upload有关
	fmt.Println("系统状态列表：",ts7szl,size)
	ast.Nil(err2)


	szlheader := SZL_HEADER {
		LENTHDR :18,
		DR   :  1,
	}

	szlList := []TS7SZLList {
		{Header : szlheader,
			//List  : [8190]uint16
			},
		{Header : szlheader,
			//List  : [8190]uint16
		},
	}
	iteCount,err2 := client.ReadSZLList(szlList)
	fmt.Println("ReadSZLList：",iteCount)
	ast.Nil(err2)



	ordercode,err6 := client.GetOrderCode()
	fmt.Println("ordercode：",ordercode)
	ast.Nil(err6)
	cpuInf,err6 := client.GetCpuInfo()
	fmt.Println("CpuInfo：",cpuInf)
	ast.Nil(err6)
	cpInf,err7 := client.GetCpInfo()
	fmt.Println("CpInfo：",cpInf)
	ast.Nil(err7)

}

func TestSecurity(t *testing.T) {
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
	ret4,err4 := client.GetProtection()
	fmt.Println("Protection级别信息：",ret4)
	ast.Nil(err4)

	err2 := client.SetSessionPassword("12345678")
	ast.Nil(err2,"设置密码成功")

	ret4,err5 := client.GetProtection()
	fmt.Println("Protection级别信息：",ret4)
	ast.Nil(err5)

	err3 := client.ClearSessionPassword()
	ast.Nil(err3,"清除密码成功")


}

func TestLowLevel(t *testing.T) {
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


	pUsrData1:= []byte{1,2,3,4,5,6}
	size,err5 := client.IsoExchangeBuffer(pUsrData1)
	fmt.Println(size)
	ast.Nil(err5)
}










