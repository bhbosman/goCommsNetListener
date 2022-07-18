package goCommsNetListener

import (
	"github.com/bhbosman/gocomms/common"
)

type userContextValue struct {
	userContext interface{}
}

func (self userContextValue) ApplyNetManagerSettings(settings *common.NetManagerSettings) error {
	return nil
}

func (self userContextValue) apply(settings *netListenManagerSettings) error {
	//settings.userContext = self.userContext
	return nil
}

func UserContextValue(userContext interface{}) *userContextValue {
	return &userContextValue{userContext: userContext}
}
