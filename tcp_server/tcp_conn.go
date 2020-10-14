package tcp_server

import (
	"context"
	"fmt"
	"net"
	"runtime"
)

// TCPListener is a TCP network listener. Clients should typically
// use variables of type Listener instead of assuming TCP.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept implements the Accept method in the Listener interface; it
// waits for the next call and returns a generic Conn.
// 继承方法覆写方法时，只能使用非指针接口
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	// AcceptTCP accepts the next incoming call and returns the new
	// connection.
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return tc, nil
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "tcp_proxy context value " + k.name
}

type conn struct {
	server     *TcpServer
	cancelCtx  context.CancelFunc
	rwc        net.Conn
	remoteAddr string
}

func (c *conn) close() {
	_ = c.rwc.Close()
}

func (c *conn) serve(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("tcp: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		c.close()
	}()
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	if c.server.Handler == nil {
		panic("handler empty")
	}
	c.server.Handler.ServeTCP(ctx, c.rwc)
}
