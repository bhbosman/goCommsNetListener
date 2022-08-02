package goCommsNetListener

import (
	"context"
	"net"
)

type ISshListenerAccept interface {
	AcceptWithContext() (net.Conn, context.CancelFunc, error)
}
