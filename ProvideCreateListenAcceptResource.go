package goCommsNetListener

import (
	"go.uber.org/fx"
	"net"
)

func ProvideCreateListenAcceptResource() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(params struct {
				fx.In
				Listener net.Listener
			}) (IListenerAccept, error) {
				return &listenerAccept{
					Listener: params.Listener,
				}, nil
			},
		},
	)
}
