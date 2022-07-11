package goCommsNetListener

import (
	"context"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"net/url"
	"time"
)

func NewNetListenApp(
	name string,
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	connectionInstancePrefix string,
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
	//stackName string,
	settings ...common.INetManagerSettingsApply) common.NetAppFuncInParamsCallback {
	return func(params common.NetAppFuncInParams) messages.CreateAppCallback {
		return messages.CreateAppCallback{
			ServiceId:         serviceIdentifier,
			ServiceDependency: serviceDependentOn,
			Name:              name,
			Callback: func() (*fx.App, context.CancelFunc, error) {
				cancelFunc := func() {}
				netListenSettings := &netListenManagerSettings{
					NetManagerSettings:    common.NewNetManagerSettings(512),
					userContext:           nil,
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
					UseProxy,
					ProxyUrl,
					ConnectionUrl,
					//stackName,
					params,
					callbackForConnectionInstance,
					fx.Options(netListenSettings.MoreOptions...),
					fx.Supply(netListenSettings),
					fx.Provide(fx.Annotated{Target: NewNetListenManager}),
					fx.Provide(fx.Annotated{Target: netListenSettings.OnCreateConnectionFactory}),
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
