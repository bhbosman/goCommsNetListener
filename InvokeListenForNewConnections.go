package goCommsNetListener

import (
	"context"
	"go.uber.org/fx"
)

func InvokeStartConnectionManagerListenForConnections(
	params struct {
		fx.In
		NetManager *NetListenManager
		Lifecycle  fx.Lifecycle
	},
) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.NetManager.ListenForNewConnections()
		},
	})
}
