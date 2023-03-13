package goCommsNetListener

import (
	"context"
	"net"
)

type IListenerAccept interface {
	AcceptWithContext() (net.Conn, context.CancelFunc, error)
}
