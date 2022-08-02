module github.com/bhbosman/goCommsNetListener

go 1.18

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20220718082038-833ca2ad99e2
	github.com/bhbosman/goConnectionManager v0.0.0-20220721070628-0f4b3c036d93
	github.com/bhbosman/gocommon v0.0.0-20220718213201-2711fee77ae4
	github.com/bhbosman/gocomms v0.0.0-20220614200341-e167364b814f
	github.com/golang/mock v1.6.0
	go.uber.org/fx v1.17.1
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.21.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require github.com/stretchr/testify v1.7.0 // indirect

require (
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf // indirect
	github.com/bhbosman/gomessageblock v0.0.0-20220617132215-32f430d7de62 // indirect
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7 // indirect
	github.com/cenkalti/backoff/v4 v4.0.0 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/reactivex/rxgo/v2 v2.5.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/bhbosman/gocomms => ../gocomms

replace github.com/bhbosman/goMessages => ../goMessages

replace github.com/bhbosman/gocommon => ../gocommon

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20220617134815-f277ff266f47

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

replace github.com/bhbosman/goFxApp => ../goFxApp

replace github.com/bhbosman/goerrors => ../goerrors

replace github.com/bhbosman/goConnectionManager => ../goConnectionManager

replace github.com/bhbosman/goprotoextra => ../goprotoextra
