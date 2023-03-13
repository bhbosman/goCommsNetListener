package goCommsNetListener

import (
	"github.com/bhbosman/gocomms/common"
	"net"
)

type netListenManagerSettings struct {
	common.NetManagerSettings
	//userContext           interface{}
	netListenerFactory interface{} //func() (net.Listener, error)
	//listenerAcceptFactory interface{} //func(ISshListenerAccept, err)
}

//func (self *netListenManagerSettings) setListenerAcceptFactory(listenerAcceptFactory func(listener net.Listener) (IListenerAccept, error)) {
//	self.listenerAcceptFactory = listenerAcceptFactory
//}

func (self *netListenManagerSettings) setListenerFactory(netListenerFactory func() (net.Listener, error)) {
	self.netListenerFactory = netListenerFactory
}
