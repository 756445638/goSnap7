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
	client := NewS7Client()
	err := client.ConnectTo("127.0.0.1", 0, 2)
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
