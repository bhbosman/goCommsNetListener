package goCommsNetListener

import (
	"go.uber.org/fx"
	"net"
)

func ProvideCreateListenAcceptResource(
	params struct {
		fx.In
		Listener net.Listener
	}) (IListenerAccept, error) {
	return &listenerAccept{
		Listener: params.Listener,
	}, nil
}
