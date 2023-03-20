package goCommsNetListener

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConn"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"net/url"
	"time"
)

func NewNetListenApp(
	name string,
	connectionInstancePrefix string,
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
	settings ...common.INetManagerSettingsApply) common.NetAppFuncInParamsCallback {
	return func(params common.NetAppFuncInParams) messages.CreateAppCallback {
		return messages.CreateAppCallback{
			Name: name,
			Callback: func() (messages.IApp, goConn.ICancellationContext, error) {

				netListenSettings := &netListenManagerSettings{
					NetManagerSettings: common.NewNetManagerSettings(512),
				}
				namedLogger := params.ZapLogger.Named(name)
				ctx, cancelFunc := context.WithCancel(params.ParentContext)
				cancellationContext := goConn.NewCancellationContext(name, cancelFunc, ctx, namedLogger, nil)

				for _, setting := range settings {
					if setting == nil {
						continue
					}
					err := setting.ApplyNetManagerSettings(&netListenSettings.NetManagerSettings)
					if err != nil {
						return nil, cancellationContext, err
					}
				}

				callbackForConnectionInstance, err := netListenSettings.Build()
				if err != nil {
					return nil, nil, err
				}

				options := common.ConnectionApp(
					time.Hour,
					time.Hour,
					name,
					connectionInstancePrefix,
					params,
					cancellationContext,
					namedLogger,
					callbackForConnectionInstance,
					fx.Options(netListenSettings.MoreOptions...),
					fx.Supply(netListenSettings),
					goCommsDefinitions.ProvideUrl("ConnectionUrl", ConnectionUrl),
					goCommsDefinitions.ProvideUrl("ProxyUrl", ProxyUrl),
					goCommsDefinitions.ProvideBool("UseProxy", UseProxy),
					fx.Provide(fx.Annotated{Target: NewNetListenManager}),
					ProvideCreateListenAcceptResource(),
					ProvideCreateListenResource(),
					fx.Invoke(InvokeStartConnectionManagerListenForConnections),
					common.InvokeCancelContext(),
					common.InvokeListenerClose(),
				)
				fxApp := fx.New(options)
				return fxApp, cancellationContext, fxApp.Err()
			},
		}
	}
}
