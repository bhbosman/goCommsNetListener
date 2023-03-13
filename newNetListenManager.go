package goCommsNetListener

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/netBase"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/url"
)

func NewNetListenManager(
	params struct {
		fx.In
		UseProxy                                 bool     `name:"UseProxy"`
		ConnectionUrl                            *url.URL `name:"ConnectionUrl"`
		ProxyUrl                                 *url.URL `name:"ProxyUrl"`
		ListenerAccept                           IListenerAccept
		ConnectionManager                        goConnectionManager.IService
		CancelCtx                                context.Context
		CancellationContext                      common.ICancellationContext
		Settings                                 *netListenManagerSettings
		ZapLogger                                *zap.Logger
		ConnectionName                           string `name:"ConnectionName"`
		ConnectionInstancePrefix                 string `name:"ConnectionInstancePrefix"`
		UniqueSessionNumber                      interfaces.IUniqueReferenceService
		AdditionalFxOptionsForConnectionInstance func() fx.Option
		GoFunctionCounter                        GoFunctionCounter.IService
	},
) (*NetListenManager, error) {

	if params.ConnectionManager.State() != IFxService.Started {
		return nil, IFxService.NewServiceStateError(
			params.ConnectionManager.ServiceName(),
			"Service in incorrect state", IFxService.Started,
			params.ConnectionManager.State())
	}

	netManager, err := netBase.NewNetManager(
		params.ConnectionName,
		params.ConnectionInstancePrefix,
		params.UseProxy,
		params.ProxyUrl,
		params.ConnectionUrl,
		params.CancelCtx,
		params.CancellationContext,
		params.ConnectionManager,
		params.ZapLogger,
		params.UniqueSessionNumber,
		params.AdditionalFxOptionsForConnectionInstance,
		params.GoFunctionCounter,
	)
	if err != nil {
		return nil, err
	}

	return &NetListenManager{
		ConnNetManager: netBase.ConnNetManager{
			NetManager: netManager,
		},
		Listener:            params.ListenerAccept,
		MaxConnections:      params.Settings.MaxConnections,
		CancellationContext: params.CancellationContext,
	}, nil
}
