package goCommsNetListener

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/netBase"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"net"
	"sync"
)

type NetListenManager struct {
	netBase.ConnNetManager
	Listener            IListenerAccept
	MaxConnections      int
	CancellationContext common.ICancellationContext
}

func (self *NetListenManager) ListenForNewConnections() error {
	actualState := self.ConnectionManager.State()
	if actualState != IFxService.Started {
		newError := IFxService.NewServiceStateError(
			self.ConnectionManager.ServiceName(),
			"Failed to start connection Listener",
			IFxService.Started,
			actualState)
		return newError
	}

	return self.GoFunctionCounter.GoRun(
		"NetListenManager.ListenForNewConnections",
		func() {
			//
			n := 0
			sem := semaphore.NewWeighted(int64(self.MaxConnections))
		loop:
			for self.CancelCtx.Err() == nil {
				n++
				self.ZapLogger.Info(
					"Trying to accept connections",
					zap.Int("Connection Count", n),
				)
				conn, connCancelFunc, err := self.acceptWithContext()
				if err != nil || err == nil && conn == nil {
					self.ZapLogger.Error(
						"Error on accept",
						zap.Error(err),
					)
					break loop
				}
				if sem.TryAcquire(1) {
					self.ZapLogger.Info("Accepted connection...")
					conn, _ = common.NewNetConnWithSemaphoreWrapper(
						conn,
						func() {
							sem.Release(1)
						},
					)
					_ = self.acceptNewClientConnection(
						self.UniqueSessionNumber.Next(self.ConnectionInstancePrefix),
						self.GoFunctionCounter,
						conn,
						connCancelFunc)
					continue
				}
				_, _ = conn.Write([]byte("ERR: To many connections\n"))
				_ = conn.Close()
			}
			//
			self.ZapLogger.Info("Leaving accept loop")
		},
	)
}

// acceptNewClientConnection will create the new connection instance. uber/fx wraps the connection, and will take care
// of its initialization and de-initialization.
//
// net.Con parameter is the new connection that was acquired.
// context.CancelFunc is a context construct that was created when the connection was formed. This will be called when
// an error occurred on the construction of the fx.App, or on the start of the initialization. It will be called when
// the exit of the de-initialization. It can assist in test cases to give an indication that the connection is closed
// and de-initialized
func (self *NetListenManager) acceptNewClientConnection(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	conn net.Conn,
	connCancelFunc context.CancelFunc,
) error {
	vv := func() {
		self.ZapLogger.Info(fmt.Sprintf("Accepted %s-%s", conn.RemoteAddr(), conn.LocalAddr()),
			zap.String("Remote Address", conn.RemoteAddr().String()),
			zap.String("LocalAddr Address", conn.LocalAddr().String()))

		connectionInstance := netBase.NewConnectionInstance(
			self.ConnectionUrl,
			self.UniqueSessionNumber,
			self.ConnectionManager,
			self.CancelCtx,
			self.AdditionalFxOptionsForConnectionInstance,
			self.ZapLogger,
		)
		connectionApp, instanceAppCtx, cancellationContext, err := connectionInstance.NewConnectionInstance(
			uniqueReference,
			goFunctionCounter,
			model.ServerConnection,
			conn,
		)
		if instanceAppCtx != nil {
			err = multierr.Append(err, instanceAppCtx.Err())
		}
		onErr := func() {
			if connCancelFunc != nil {
				connCancelFunc()
			}
			if cancellationContext != nil {
				cancellationContext.Cancel()
			}
			err = multierr.Append(err, conn.Close())
		}
		if err != nil {
			onErr()
			return
		}
		err = connectionApp.Start(context.Background())
		if err != nil {
			onErr()
			return
		}
		connectionShutdown := registerConnectionShutdown(
			uniqueReference,
			connectionApp,
			self.ZapLogger,
			self.CancellationContext,
		)
		_, _ = self.CancellationContext.Add(uniqueReference, connectionShutdown)
	}
	return self.GoFunctionCounter.GoRun(
		"NetListenManager.acceptNewClientConnection.03",
		vv,
	)
	//vv()
	//return nil
}

func registerConnectionShutdown(
	connectionId string,
	connectionApp messages.IApp,
	logger *zap.Logger,
	CancellationContext ...common.ICancellationContext,
) func(cancelCtx common.ICancellationContext) {
	mutex := sync.Mutex{}
	cancelCalled := false
	return func(cancelCtx common.ICancellationContext) {
		mutex.Lock()
		b := cancelCalled
		cancelCalled = true
		mutex.Unlock()
		if !b {
			errInGoRoutine := connectionApp.Stop(context.Background())
			if errInGoRoutine != nil {
				logger.Error(
					"Stopping error. not really a problem. informational",
					zap.Error(errInGoRoutine))
			}
			for _, instance := range CancellationContext {
				_ = instance.Remove(connectionId)
			}
		}
	}
}

func (self *NetListenManager) acceptWithContext() (net.Conn, context.CancelFunc, error) {
	return self.Listener.AcceptWithContext()
}
