package goCommsNetListener

import (
	"go.uber.org/fx"
	"net"
)

type overrideListenerFactory struct {
	listenerFactory func() (net.Listener, error)
}

func (self *overrideListenerFactory) apply(settings *netListenManagerSettings) (fx.Option, error) {
	settings.setListenerFactory(self.listenerFactory)
	return nil, nil
}

func NewOverrideListener(listenerFactory func() (net.Listener, error)) *overrideListenerFactory {
	return &overrideListenerFactory{listenerFactory: listenerFactory}
}
