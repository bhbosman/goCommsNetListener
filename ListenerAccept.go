package goCommsNetListener

import (
	"context"
	"net"
)

type listenerAccept struct {
	Listener interface {
		Accept() (net.Conn, error)
	}
}

func (self *listenerAccept) AcceptWithContext() (net.Conn, context.CancelFunc, error) {
	conn, err := self.Listener.Accept()
	if err != nil {
		return nil, nil, err
	}
	return conn,
		func() {
			// do nothing
		},
		nil
}
