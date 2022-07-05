package goCommsNetListener

import "github.com/bhbosman/gocomms/common"

type iListenAppSettingsApply interface {
	common.INetManagerSettingsApply
	apply(settings *netListenManagerSettings) error
}
