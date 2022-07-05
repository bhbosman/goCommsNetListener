package goCommsNetListener

import (
	"context"
	"go.uber.org/fx"
	"net"
	"net/url"
)

func ProvideCreateListenResource(
	params struct {
		fx.In
		Lifecycle fx.Lifecycle

		UseProxy      bool     `name:"UseProxy"`
		ConnectionUrl *url.URL `name:"ConnectionUrl"`
		ProxyUrl      *url.URL `name:"ProxyUrl"`
	}) (net.Listener, error) {
	con, err := net.Listen(
		params.ConnectionUrl.Scheme,
		params.ConnectionUrl.Host,
	)
	if err != nil {
		return nil, err
	}
	params.Lifecycle.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return con.Close()
		},
	})
	return con, nil
}
