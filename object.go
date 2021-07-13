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
func (s *S7Server) SetRWAreaCallback(handle func(sender int, operation int, tag *PS7Tag, userData *C.void)) error {
	return Srv_SetRWAreaCallback(s.server, func(_ uintptr, sender int, operation int, tag *PS7Tag, userData *C.void) {
		handle(sender, operation, tag, userData)
	}, uintptr(s.server))
}

func (s *S7Server) Destroy() {
	Srv_Destroy(s.server)

}
func (s *S7Server) GetParam(paraNumber ParamNumber) (value interface{}, err error) {
	return Srv_GetParam(s.server, paraNumber)
}

func (s *S7Server) SetParam(paraNumber ParamNumber, value interface{}) (err error) {
	return Srv_SetParam(s.server, paraNumber, value)
}

func (s *S7Server) StartTo(Address string) (err error) {
	return Srv_StartTo(s.server, Address)
}

func (s *S7Server) Start() (err error) {
	return Srv_Start(s.server)
}

func (s *S7Server) Stop() (err error) {
	return Srv_Stop(s.server)
}

func (s *S7Server) RegisterArea(AreaCode int, Index uint16, pUsrData []byte, Size int) (err error) {
	return Srv_RegisterArea(s.server, AreaCode, Index, pUsrData, Size)

}

func (s *S7Server) UnregisterArea(AreaCode int, Index uint16) (err error) {
	return Srv_UnregisterArea(s.server, AreaCode, Index)
}

func (s *S7Server) LockArea(AreaCode int, Index uint16) (err error) {
	return Srv_LockArea(s.server, AreaCode, Index)
}

func (s *S7Server) UnlockArea(AreaCode int, Index uint16) (err error) {
	return Srv_UnlockArea(s.server, AreaCode, Index)
}

func (s *S7Server) GetStatus(CpuStatus int, ClientsCount int) (ServerStatus int, err error) {
	return Srv_GetStatus(s.server, CpuStatus, ClientsCount)
}

func (s *S7Server) SetCpuStatus(CpuStatus int) (err error) {
	return Srv_SetCpuStatus(s.server, CpuStatus)
}

func (s *S7Server) ClearEvents() (err error) {
	return Srv_ClearEvents(s.server)
}

func (s *S7Server) PickEvent(pEvent TSrvEvent, EvtReady int) (err error) {
	return Srv_PickEvent(s.server, pEvent, EvtReady)
}

func (s *S7Server) GetMask(MaskKind int) (Mask uint32, err error) {
	return Srv_GetMask(s.server, MaskKind)
}

func (s *S7Server) SetMask(MaskKind int, Mask uint32) (err error) {
	return Srv_SetMask(s.server, MaskKind, Mask)
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
