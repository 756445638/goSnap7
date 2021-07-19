package snap7go

import "C"

func NewS7Server() *S7Server {
	server := &S7Server{}
	server.server = Srv_Create()
	return server
}

type S7Server struct {
	server S7Object
}

func (s *S7Server) SetEventsCallback(handle func(*TSrvEvent)) error {
	return Srv_SetEventsCallback(s.server, func(usrPtr uintptr, event *TSrvEvent) {
		handle(event)
	}, uintptr(s.server))
}

func (s *S7Server) SetReadEventsCallback(handle func(*TSrvEvent)) error {
	return Srv_SetReadEventsCallback(s.server, func(usrPtr uintptr, event *TSrvEvent) {
		handle(event)
	}, uintptr(s.server))
}

func (s *S7Server) SetRWAreaCallback(handle func(sender int, operation Operation, tag *PS7Tag, userData uintptr) int) error {
	return Srv_SetRWAreaCallback(s.server,
		func(_ uintptr, sender int, operation int, tag *PS7Tag, userData uintptr) int {
			return handle(sender, Operation(operation), tag, userData)
		}, uintptr(s.server))
}

//Client
func NewS7Client() *S7Client {
	server := &S7Client{}
	server.client = Cli_Create()
	return server
}

type S7Client struct {
	client S7Object
}

func (c *S7Client) SetAsCallback(handle func(opCode int, opResult int)) error {
	return Cli_SetAsCallback(c.client, func(usrptr uintptr, opCode int, opResult int) {
		handle(opCode, opResult)
	}, uintptr(c.client))
}
