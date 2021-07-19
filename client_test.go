package snap7go

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	clientDesignated := NewS7Client()

	defer clientDesignated.Destroy()
	//CONNTYPE_PG、CONNTYPE_OP、CONNTYPE_BASIC
	err = clientDesignated.SetConnectionType(CONNTYPE_BASIC)
	ast.Nil(err)
	err = clientDesignated.SetConnectionParams("127.0.0.1", 0x1000, 0x1000)
	ast.Nil(err)

	//在ConnectTo前后都可以
	err = clientDesignated.SetParam(P_i32_SendTimeout, int32(4))
	ast.Nil(err)
	//连接指定地址(192.168.187.1)
	err = clientDesignated.ConnectTo("127.0.0.1", 0, 1)
	ast.Nil(err)

	paradata, err := clientDesignated.GetParam(P_i32_SendTimeout)
	ast.Nil(err)
	ast.Equal(int32(4), paradata)
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

	//Security
	//查询保护级别 1 -mode selector   0-no password   1- CPU  2-:Mode selector setting RUN-P     0-Startup switch setting :undefined,
	ret4, err4 := client.GetProtection()
	//fmt.Println("Protection级别信息：",ret4)     {1 0 1 2 0}
	ast.Nil(err4)
	//设置8位用户密码
	err2 := client.SetSessionPassword("12345678")
	ast.Nil(err2)

	ret4, err5 := client.GetProtection()
	fmt.Println("Protection级别信息：", ret4)
	ast.Nil(err5)

	pUsrData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	//func (c *S7Client) Upload(blockType Block, blockNum int, pUsrData []byte) (size int, err error) {

	ret, err := client.FullUpload(Block_OB, 1, pUsrData)
	fmt.Println("fullUpload Buffer size:", ret)
	ast.Nil(err)

	rete, err := client.ListBlocks()
	fmt.Println("ListBlocks:", rete)
	ast.Nil(err)

	data, itemCounter, err := client.ListBlocksOfType(Block_OB)
	fmt.Println("TS7BlocksOfType", data)
	fmt.Println(itemCounter)
	ast.Nil(err)

	ret1, err := client.GetAgBlockInfo(Block_OB, 1)
	fmt.Println("AgBlockInfo:", ret1)
	ast.Nil(err)

	//fmt.Println(ret)
	//ast.Nil(err)

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

	ts7szl, size, err2 := client.ReadSZL(0x0232, 0x0004) //与upload有关
	fmt.Println("系统状态列表：", ts7szl, size)
	ast.Nil(err2)

	//szlheader := SZL_HEADER {
	//	LENTHDR :18,
	//	DR   :  1,
	//}
	//
	//szlList := []TS7SZLList {   这个看exmple不太像是自己输入的，现在也不清楚header里面的LENTHDR以及DR什么形式
	//	{Header : szlheader,
	//		//List  : [8190]uint16
	//		},
	//	{Header : szlheader,
	//		//List  : [8190]uint16
	//	},
	//}

	//iteCount,err2 := client.ReadSZLList(szlList)
	//fmt.Println("ReadSZLList：",iteCount)
	//ast.Nil(err2)

	ordercode, err6 := client.GetOrderCode()
	fmt.Println("ordercode：", ordercode)
	ast.Nil(err6)
	cpuInf, err6 := client.GetCpuInfo()
	fmt.Println("CpuInfo：", cpuInf.GetASName())
	ast.Nil(err6)
	cpInf, err7 := client.GetCpInfo()
	fmt.Println("CpInfo：", cpInf)
	ast.Nil(err7)

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

func TestLowLevelCli(t *testing.T) { //未完成
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

	pUsrData1 := []byte{1, 2, 3, 4, 5, 6}
	size, err5 := client.IsoExchangeBuffer(pUsrData1)
	fmt.Println(size)
	ast.Nil(err5)
}
