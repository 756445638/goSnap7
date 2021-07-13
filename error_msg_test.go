package snap7go

import (
	"testing"
)

func TestCliErrMsg(t *testing.T) {
	err := Cli_ErrorText(0x00300000)
	if err == nil {
		t.Fatalf("0x00300000 is a error")
	}

}
