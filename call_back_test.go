package snap7go

import (
	"fmt"
	"testing"
	"time"
)

func TestSomeCallBack(t *testing.T) {
	server := NewS7Server()
	server.SetEventsCallback(justPrintEvent)
	server.SetReadEventsCallback(justPrintEvent)
	const duration = time.Millisecond * 50
	var data [1024]byte

	go func() {
		for ; ; time.Sleep(duration) {
			for k, _ := range data {
				data[k]++
			}
		}
	}()
	err := server.RegisterArea(SrvAreaPA, 0, data[:])
	if err != nil {
		t.Fatal(err)
	}
	err = server.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop()
		if err != nil {
			t.Fatal(err)
			return
		}
		server.Destroy()
	}()
	client := NewS7Client()
	err = client.ConnectTo("127.0.0.1", 0, 1)
	if err != nil {
		t.Fatal(err)
		return
	}
	for i := 0; i < 10; func() {
		i++
		time.Sleep(duration)
	}() {
		data, err := client.ReadArea(S7AreaPA, 0, 3, 10, S7WLWord)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("read data:", data)
	}

	{
		// 读取区域的最后十个字节
		data := []TS7DataItemGo{
			{
				Area:     int32(S7AreaPA),
				WordLen:  int32(S7WLByte),
				DBNumber: 0,
				Start:    1014,
				Amount:   10,
			},
			{
				Area:     int32(S7AreaPA),
				WordLen:  int32(S7WLByte),
				DBNumber: 0,
				Start:    100,
				Amount:   10,
			},
		}
		err := client.ReadMultiVars(data)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("!!!!!!!!read data:", data[0].Pdata)
		for k, _ := range data[0].Pdata {
			data[0].Pdata[k] = 100
		}
		// 把值全部改成100
		err = client.WriteMultiVars(data)
		if err != nil {
			t.Fatal(err)
			return
		}
		err = client.ReadMultiVars(data)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("read data:", data[0].Pdata)
		if data[0].Pdata[0] != 100 {
			t.Fatalf("value shoudle be 100\n")
		}
	}

}

/*

 */
func TestSomeWordLenStart(t *testing.T) {
	server := NewS7Server()
	server.SetEventsCallback(justPrintEvent)
	server.SetReadEventsCallback(justPrintEvent)
	const duration = time.Millisecond * 50
	var data [10]byte

	for k, _ := range data {
		data[k] = 100
	}

	err := server.RegisterArea(SrvAreaDB, 0, data[:])
	if err != nil {
		t.Fatal(err)
	}
	err = server.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop()
		if err != nil {
			t.Fatal(err)
			return
		}
		server.Destroy()
	}()
	client := NewS7Client()
	err = client.ConnectTo("127.0.0.1", 0, 2)
	if err != nil {
		t.Fatal(err)
		return
	}
	{
		data, err := client.ReadArea(S7AreaDB, 0, 1, 1, S7WLBit)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("bit value:", data)
		err = client.WriteArea(S7AreaDB, 0, 1, S7WLBit, []byte{1})
		if err != nil {
			t.Fatal(err)
			return
		}
		data, err = client.ReadArea(S7AreaDB, 0, 1, 1, S7WLBit)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("bit value:", data)
	}

	/*
		这里测试的是start是否包含长度信息
	*/
	{
		_, err := client.ReadArea(S7AreaDB, 0, 6, 1, S7WLReal)
		if err != nil {
			fmt.Println("start不包含长度")
		} else {
			fmt.Println("start包含长度")
		}
	}
}

func TestSomeWordLenStart222(t *testing.T) {
	if testing.Short() {
		return
	}
	server := NewS7Server()
	var data [10]byte

	for k, _ := range data {
		data[k] = 100
	}

	err := server.RegisterArea(SrvAreaDB, 1, data[:])
	if err != nil {
		t.Fatal(err)
	}
	err = server.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop()
		if err != nil {
			t.Fatal(err)
			return
		}
		server.Destroy()
	}()
	client := NewS7Client()
	err = client.ConnectTo("127.0.0.1", 0, 2)
	if err != nil {
		t.Fatal(err)
		return
	}
	{
		data, err := client.ReadArea(S7AreaDB, 1, 1, 1, S7WLBit)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("!!!!!!!!!", data)
		err = client.WriteArea(S7AreaDB, 1, 1, S7WLBit, []byte{1})
		if err != nil {
			t.Fatal(err)
			return
		}
		data, err = client.ReadArea(S7AreaDB, 1, 1, 1, S7WLBit)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("!!!!!!!!!", data)
	}
}
func TestSomeSetRWAreaCallback(t *testing.T) {
	server := NewS7Server()
	err := server.SetEventsCallback(justPrintEvent)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = server.SetReadEventsCallback(justPrintEvent)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = server.SetRWAreaCallback(func(sender int32, operation Operation, tag *PS7Tag, userData uintptr) SrvErrCode {
		if operation == 0 {
			// read
			CopyToC([]byte{1, 2, 3}, userData)
			return 0
		} else {
			data := GetBytesFromC(userData, int(dataLength(S7WL(tag.WordLen), tag.Size)))
			if data[0] != 4 || data[1] != 5 || data[2] != 6 {
				panic("data not right")
			}
			// write
			// false error
			return 0x20000
		}
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	err = server.Start()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func() {
		err = server.Stop()
		if err != nil {
			t.Fatal(err)
			return
		}
		server.Destroy()
	}()
	client := NewS7Client()
	err = client.ConnectTo("127.0.0.1", 0, 2)
	if err != nil {
		t.Fatal(err)
		return
	}
	data, err := client.ReadArea(S7AreaDB, 0, 0, 3, S7WLByte)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println("data:", data)
	err = client.WriteArea(S7AreaDB, 0, 0, S7WLByte, []byte{4, 5, 6})
	if err == nil {
		t.Fatalf("shoudle be a false error:%v\n", err)
		return
	}
}

func TestSetRWAreaCallbackInterface(t *testing.T) {
	server := NewS7Server()
	err := server.SetEventsCallback(justPrintEvent)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = server.SetReadEventsCallback(justPrintEvent)
	if err != nil {
		t.Fatal(err)
		return
	}
	var handle1 handle
	err = server.SetRWAreaCallbackInterface(handle1)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = server.Start()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func() {
		err = server.Stop()
		if err != nil {
			t.Fatal(err)
			return
		}
		server.Destroy()
	}()
	client := NewS7Client()
	err = client.ConnectTo("127.0.0.1", 0, 2)
	if err != nil {
		t.Fatal(err)
		return
	}
	data, err := client.ReadArea(S7AreaDB, 0, 0, 5, S7WLByte)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println("data:", data)
	err = client.WriteArea(S7AreaDB, 0, 0, S7WLByte, []byte{6, 7, 8, 9})
	if err != nil {
		t.Fatal(err)
		return
	}
}

type handle struct{}

func (h handle) Read(sender int32, tag *PS7Tag, data []byte) (errCode SrvErrCode) {
	return 0
}
func (h handle) Write(sender int32, tag *PS7Tag, data []byte) (errCode SrvErrCode) {
	return 0
}

func TestSetAsCallback(t *testing.T) {
	server := NewS7Server()
	server.SetEventsCallback(justPrintEvent)
	server.SetReadEventsCallback(justPrintEvent)
	server.Start()
	defer func() {
		err := server.Stop()
		if err != nil {
			t.Fatal(err)
			return
		}
		server.Destroy()
	}()
	client := NewS7Client()
	err := client.ConnectTo("127.0.0.1", 0, 2)
	if err != nil {
		t.Fatal(err)
		return
	}
	var JobDone = false
	Pfn_CliCompletion := func(opCode int32, opResult int32) {
		JobDone = true
		fmt.Println("JobDone:", JobDone)
	}

	var db [1024]byte
	server.RegisterArea(SrvAreaPE, 1, db[:])

	err = client.SetAsCallback(Pfn_CliCompletion)
	if err != nil {
		t.Fatal(err)
		return

	}

	pUsrData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	err = client.AsWriteArea(S7AreaPE, 1, 0, S7WLByte, pUsrData)
	if err != nil {
		t.Fatal(err)
		return

	}
	for JobDone == false {
		time.Sleep(time.Millisecond)
	}
	JobDone = false
	_, err = client.AsReadArea(S7AreaPE, 1, 0, 12, S7WLByte)
	if err != nil {
		t.Fatal(err)
		return
	}
	for JobDone == false {
		time.Sleep(time.Millisecond)
	}
}
