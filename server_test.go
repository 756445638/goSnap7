package snap7go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerAdministrative(t *testing.T) {
	ast := assert.New(t)
	server := NewS7Server()

	err := server.Start()
	ast.Nil(err)

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
}
