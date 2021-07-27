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

func (s *S7Server) SetRWAreaCallback(handle func(sender int32, operation Operation, tag *PS7Tag, userData uintptr) SrvErrCode) error {
	return Srv_SetRWAreaCallback(s.server,
		func(_ uintptr, sender int32, operation Operation, tag *PS7Tag, userData uintptr) SrvErrCode {
			return handle(sender, Operation(operation), tag, userData)
		}, uintptr(s.server))
}

func (s *S7Server) SetRWAreaCallbackInterface(handle RWAreaCallbackInterface) error {
	return Srv_SetRWAreaCallback(s.server,
		func(_ uintptr, sender int32, operation Operation, tag *PS7Tag, userData uintptr) SrvErrCode {
			if operation == OperationRead {
				data := make([]byte, dataLength(S7WL(tag.WordLen), tag.Size))
				errCode := handle.Read(sender, tag, data)
				if errCode != 0 {
					return errCode
				}
				CopyToC(data, userData)
				return errCode
			} else {
				// write
				data := GetBytesFromC(userData, int(dataLength(S7WL(tag.WordLen), tag.Size)))
				return handle.Write(sender, tag, data)
			}
		}, uintptr(s.server))
}

type RWAreaCallbackInterface interface {
	Read(sender int32, tag *PS7Tag, ret []byte) (errCode SrvErrCode)
	Write(sender int32, tag *PS7Tag, data []byte) (errCode SrvErrCode)
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

func (c *S7Client) SetAsCallback(handle func(opCode int32, opResult JobStatus)) error {
	return Cli_SetAsCallback(c.client, func(usrptr uintptr, opCode int32, opResult JobStatus) {
		handle(opCode, opResult)
	}, uintptr(c.client))
}
