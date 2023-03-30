all: IConnectionReactorFactory


IConnectionReactorFactory:
	mockgen -package goCommsNetListener -generateWhat mockgen -destination ConnMock.go net Conn,Addr,Listener
	mockgen -package goCommsNetListener -generateWhat mockgen -destination IListenerAcceptMock.go . IListenerAccept







