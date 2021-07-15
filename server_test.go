package snap7go

import "testing"

func TestServerAdministrative(t *testing.T) {
	server := NewS7Server()

	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}

}
