package tcp_server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed     = errors.New("tcp: Server closed")
	ErrAbortHandler     = errors.New("tcp: abort TCPHandler")
	ServerContextKey    = &contextKey{"tcp-server"} // 不能使用string，而需要将string包装入一个struct
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type TCPHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

type TcpServer struct {
	Addr    string
	Handler TCPHandler
	err     error
	BaseCtx context.Context

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	KeepAliveTimeout time.Duration

	mu         sync.Mutex
	inShutdown int32
	doneChan   chan struct{}
	l          *onceCloseListener
}

type onceCloseListener struct {
	// A Listener is a generic network listener for stream-oriented protocols.
	// Accept() Close() Addr()
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() {
	oc.closeErr = oc.Listener.Close()
}

func (server *TcpServer) shuttingDown() bool {
	return atomic.LoadInt32(&server.inShutdown) != 0
}

func (server *TcpServer) Close() error {
	atomic.StoreInt32(&server.inShutdown, 1)
	close(server.doneChan) // 关闭channel
	server.l.Close()       // 执行listener关闭
	return nil
}

func (server *TcpServer) ListenAndServe() error {
	if server.shuttingDown() {
		return ErrServerClosed
	}

	if server.doneChan == nil {
		server.doneChan = make(chan struct{})
	}

	addr := server.Addr
	if addr == "" {
		return errors.New("need addr")
	}

	// Listen announces on the local network address.
	//
	// The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	//
	// For TCP networks, if the host in the address parameter is empty or
	// a literal unspecified IP address, Listen listens on all available
	// unicast and anycast IP addresses of the local system.
	// To only use IPv4, use network "tcp4".
	// The address can use a host name, but this is not recommended,
	// because it will create a listener for at most one of the host's IP
	// addresses.
	// If the port in the address parameter is empty or "0", as in
	// "127.0.0.1:" or "[::1]:0", a port number is automatically chosen.
	// The Addr method of Listener can be used to discover the chosen
	// port.
	//
	// See func Dial for a description of the network and address
	// parameters.
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return server.Serve(tcpKeepAliveListener{
		ln.(*net.TCPListener)})
}

func (server *TcpServer) Serve(l net.Listener) error {
	server.l = &onceCloseListener{Listener: l}
	defer server.l.Close() // 执行listener关闭

	if server.BaseCtx == nil {
		// Background returns a non-nil, empty Context. It is never canceled, has no
		// values, and has no deadline. It is typically used by the main function,
		// initialization, and tests, and as the top-level Context for incoming
		// requests.
		server.BaseCtx = context.Background()
	}
	baseCtx := server.BaseCtx

	// WithValue returns a copy of parent in which the value associated with key is
	// val.
	//
	// Use context Values only for request-scoped data that transits processes and
	// APIs, not for passing optional parameters to functions.
	//
	// The provided key must be comparable and should not be of type
	// string or any other built-in type to avoid collisions between
	// packages using context. Users of WithValue should define their own
	// types for keys. To avoid allocating when assigning to an
	// interface{}, context keys often have concrete type
	// struct{}. Alternatively, exported context key variables' static
	// type should be a pointer or interface.
	ctx := context.WithValue(baseCtx, ServerContextKey, server)

	for {
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-server.getDoneChan():
				return ErrServerClosed
			default:
			}
			fmt.Printf("accept fail, err: %v\n", e)
			continue
		}
		conn := server.newConn(rw)
		go conn.serve(ctx)
	}

	return nil
}

func (server *TcpServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: server,
		rwc:    rwc,
	}
	// 设置参数
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := c.server.KeepAliveTimeout; d != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
	return c
}

func (server *TcpServer) getDoneChan() <-chan struct{} {
	server.mu.Lock()
	defer server.mu.Unlock()
	if server.doneChan == nil {
		server.doneChan = make(chan struct{})
	}
	return server.doneChan
}
