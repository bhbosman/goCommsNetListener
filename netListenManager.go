package goCommsNetListener

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/netBase"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"net"
)

type NetListenManager struct {
	netBase.ConnNetManager
	Listener           ISshListenerAccept
	MaxConnections     int
	OnCreateConnection goCommsDefinitions.IOnCreateConnection
}

func (self *NetListenManager) ListenForNewConnections() error {
	goFunc := func() {
		n := 0
		sem := semaphore.NewWeighted(int64(self.MaxConnections))
	loop:
		for self.CancelCtx.Err() == nil {
			n++
			self.ZapLogger.Info("Trying to accept connections", zap.Int("Connection Count", n))
			conn, connCancelFunc, err := self.acceptWithContext()
			if err != nil || err == nil && conn == nil {
				self.ZapLogger.Error("Error on accept", zap.Error(err))
				break loop
			}
			if sem.TryAcquire(1) {
				self.ZapLogger.Info("Accepted connection...")
				conn = common.NewNetConnWithSemaphoreWrapper(conn, sem)
				self.acceptNewClientConnection(
					self.UniqueSessionNumber.Next(self.ConnectionInstancePrefix),
					self.GoFunctionCounter,
					conn,
					connCancelFunc)
				continue
			}
			_, _ = conn.Write([]byte("ERR: To many connections\n"))
			_ = conn.Close()
		}
		self.ZapLogger.Info("Leaving accept loop")
	}
	var err error = nil
	// check if connection manager state is IFxService.Started
	actualState := self.ConnectionManager.State()
	if actualState != IFxService.Started {
		newError := IFxService.NewServiceStateError(
			self.ConnectionManager.ServiceName(),
			"Failed to start connection Listener",
			IFxService.Started,
			actualState)
		err = multierr.Append(err, newError)
	}

	if err == nil {
		// this function is part of the GoFunctionCounter count
		go func() {
			functionName := self.GoFunctionCounter.CreateFunctionName("NetListenManager.ListenForNewConnections")
			defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
				_ = GoFunctionCounter.Remove(name)
			}(self.GoFunctionCounter, functionName)
			_ = self.GoFunctionCounter.Add(functionName)

			//
			goFunc()
		}()
	}
	return err
}

// acceptNewClientConnection will create the new connection instance. uber/fx wraps the connection, and will take care
//of its initialization and de-initialization.
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
) {
	// this function is part of the GoFunctionCounter count
	f := func(
		conn net.Conn,
	) {
		functionName := self.GoFunctionCounter.CreateFunctionName("NetListenManager.acceptNewClientConnection.01")
		defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
			_ = GoFunctionCounter.Remove(name)
		}(self.GoFunctionCounter, functionName)
		_ = self.GoFunctionCounter.Add(functionName)

		//
		self.ZapLogger.Info(fmt.Sprintf("Accepted %s-%s", conn.RemoteAddr(), conn.LocalAddr()),
			zap.String("Remote Address", conn.RemoteAddr().String()),
			zap.String("LocalAddr Address", conn.LocalAddr().String()))
		connectionApp, ctx, cancelFunc := self.NewConnectionInstance(
			uniqueReference,
			goFunctionCounter,
			model.ServerConnection,
			conn,
		)
		err := connectionApp.Err()
		if ctx != nil {
			err = multierr.Append(err, ctx.Err())
		}
		if err != nil {
			self.ZapLogger.Error("Error in fxApp.Err() when creating NewConnectionInstance()",
				zap.Error(connectionApp.Err()))
			if cancelFunc != nil {
				cancelFunc()
			}
			err = conn.Close()
			if err != nil {
				self.ZapLogger.Error("Informational error on connection.close(). No serious issue here",
					zap.Error(connectionApp.Err()))
			}
		}
		if self.OnCreateConnection != nil {
			self.OnCreateConnection.OnCreateConnection(uniqueReference, connectionApp.Err(), ctx, cancelFunc)
		}

		if err != nil {
			if connCancelFunc != nil {
				connCancelFunc()
			}
			return
		}
		// TODO: Adhere to timeouts
		err = connectionApp.Start(context.Background())
		if err != nil {
			if connCancelFunc != nil {
				connCancelFunc()
			}
			return
		}

		// this function is part of the GoFunctionCounter count
		go func(app *fx.App, ctx context.Context, cancelFunc context.CancelFunc) {
			functionName := self.GoFunctionCounter.CreateFunctionName("NetListenManager.acceptNewClientConnection.02")
			defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
				_ = GoFunctionCounter.Remove(name)
			}(self.GoFunctionCounter, functionName)
			_ = self.GoFunctionCounter.Add(functionName)

			//
			<-ctx.Done()
			// TODO: Adhere to timeouts
			errInGoRoutine := app.Stop(context.Background())
			if errInGoRoutine != nil {
				self.ZapLogger.Error(
					"Stopping error. not really a problem. informational",
					zap.Error(errInGoRoutine))
			}
			if connCancelFunc != nil {
				cancelFunc()
			}
		}(connectionApp, ctx, connCancelFunc)
	}

	// this function is part of the GoFunctionCounter count
	go func() {
		functionName := self.GoFunctionCounter.CreateFunctionName("NetListenManager.acceptNewClientConnection.03")
		defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
			_ = GoFunctionCounter.Remove(name)
		}(self.GoFunctionCounter, functionName)
		_ = self.GoFunctionCounter.Add(functionName)

		//
		f(conn)
	}()
}

func (self *NetListenManager) acceptWithContext() (net.Conn, context.CancelFunc, error) {
	return self.Listener.AcceptWithContext()
}
