package sqlite3

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	//ErrClosed 连接池已经关闭Error
	ErrClosed = errors.New("pool is closed")
)

type PoolConfig struct {
	Min         int
	Max         int
	Create      func() (interface{}, error)
	Close       func(interface{}) error
	Ping        func(interface{}) error
	IdleTimeout time.Duration
}
type BasePool interface {
	Release()
	Get() (interface{}, error)
	Put(interface{}) error
	Ping(conn interface{}) error
	Close(conn interface{}) error
}
type Pool struct {
	mu          sync.Mutex
	conns       chan *idleConn
	create      func() (interface{}, error)
	close       func(interface{}) error
	ping        func(interface{}) error
	idleTimeout time.Duration
}

type idleConn struct {
	conn interface{}
	t    time.Time
}

func CreatePool(PoolConfig *PoolConfig) (BasePool, error) {
	if PoolConfig.Min < 0 || PoolConfig.Max <= 0 {
		return nil, errors.New("invalid capacity settings")
	}
	if PoolConfig.Create == nil {
		return nil, errors.New("create func settings")
	}
	if PoolConfig.Close == nil {
		return nil, errors.New("close func settings")
	}
	c := &Pool{
		conns:       make(chan *idleConn, PoolConfig.Max),
		create:      PoolConfig.Create,
		close:       PoolConfig.Close,
		idleTimeout: PoolConfig.IdleTimeout,
	}
	if PoolConfig.Ping != nil {
		c.ping = PoolConfig.Ping
	}
	for i := 0; i < PoolConfig.Min; i++ {
		conn, err := c.create()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}
	return c, nil
}
func (c *Pool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}
func (c *Pool) Get() (interface{}, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					c.Close(wrapConn.conn)
					continue
				}
			}
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					fmt.Println("conn is not able to be connected: ", err)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			c.mu.Lock()
			if c.create == nil {
				c.mu.Unlock()
				continue
			}
			conn, err := c.create()
			c.mu.Unlock()
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
	}
}

func (c *Pool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("ccontinue is nil. rejecting")
	}
	c.mu.Lock()
	if c.conns == nil {
		c.mu.Unlock()
		return c.Close(conn)
	}
	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		c.mu.Unlock()
		return nil
	default:
		c.mu.Unlock()
		return c.Close(conn)
	}
}

func (c *Pool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.close == nil {
		return nil
	}
	return c.close(conn)
}

func (c *Pool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.ping(conn)
}

func (c *Pool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.create = nil
	c.ping = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()
	if conns == nil {
		return
	}
	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}
func (c *Pool) Len() int {
	return len(c.getConns())
}
