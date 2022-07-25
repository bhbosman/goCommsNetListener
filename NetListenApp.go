package goCommsNetListener

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
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
			Callback: func() (messages.IApp, context.CancelFunc, error) {
				cancelFunc := func() {}
				netListenSettings := &netListenManagerSettings{
					NetManagerSettings:    common.NewNetManagerSettings(512),
					netListenerFactory:    ProvideCreateListenResource,
					listenerAcceptFactory: ProvideCreateListenAcceptResource,
				}

				for _, setting := range settings {
					if setting == nil {
						continue
					}
					if listenAppSettingsApply, ok := setting.(iListenAppSettingsApply); ok {
						err := listenAppSettingsApply.apply(netListenSettings)
						if err != nil {
							return nil, cancelFunc, err
						}
					} else {
						err := setting.ApplyNetManagerSettings(&netListenSettings.NetManagerSettings)
						if err != nil {
							return nil, cancelFunc, err
						}
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
					callbackForConnectionInstance,
					fx.Options(netListenSettings.MoreOptions...),
					fx.Supply(netListenSettings),
					goCommsDefinitions.ProvideUrl("ConnectionUrl", ConnectionUrl),
					goCommsDefinitions.ProvideUrl("ProxyUrl", ProxyUrl),
					goCommsDefinitions.ProvideBool("UseProxy", UseProxy),

					fx.Provide(fx.Annotated{Target: NewNetListenManager}),
					fx.Provide(fx.Annotated{Target: netListenSettings.listenerAcceptFactory}),
					fx.Provide(fx.Annotated{Target: netListenSettings.netListenerFactory}),
					fx.Invoke(invokeListenForNewConnections),
				)
				fxApp := fx.New(options)
				return fxApp, cancelFunc, fxApp.Err()
			},
		}
	}
}
