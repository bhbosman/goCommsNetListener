package goCommsNetListener

import (
	"github.com/bhbosman/gocomms/common"
	"net"
)

type overrideListenerAcceptFactory struct {
	listenerAcceptFactory func(listener net.Listener) (IListenerAccept, error)
}

func (self *overrideListenerAcceptFactory) ApplyNetManagerSettings(settings *common.NetManagerSettings) error {
	return nil
}

func (self *overrideListenerAcceptFactory) apply(settings *netListenManagerSettings) error {
	settings.setListenerAcceptFactory(self.listenerAcceptFactory)
	return nil
}

func NewOverrideListenerAcceptFactory(listenerAcceptFactory func(listener net.Listener) (IListenerAccept, error)) *overrideListenerAcceptFactory {
	return &overrideListenerAcceptFactory{
		listenerAcceptFactory: listenerAcceptFactory,
	}
}
