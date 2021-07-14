package snap7go

import (
	"fmt"
	"testing"
	"time"
)

func TestSomeCallBack(t *testing.T) {
	server := NewS7Server()
	server.SetEventsCallback(func(e *TSrvEvent) {
		s, err := Srv_EventText(e)
		if err != nil {
			t.Fatalf("Srv_EventText failed,err:%v\n", err)
			return
		}
		fmt.Println(s)
	})
	var data [1024]byte
	go func() {
		for ; ; time.Sleep(time.Second) {
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
}
