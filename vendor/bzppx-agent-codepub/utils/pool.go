package utils

import (
	"fmt"
	"sync"
)

//ConnPool to use
type ConnPool interface {
	Get() (conn interface{}, err error)
	Put(conn interface{})
	ReleaseAll()
	Len() (length int)
}
type poolConfig struct {
	Factory    func() (interface{}, error)
	IsActive   func(interface{}) bool
	Release    func(interface{})
	InitialCap int
	MaxCap     int
}

func newNetPool(poolConfig poolConfig) (pool ConnPool, err error) {
	p := netPool{
		config: poolConfig,
		conns:  make(chan interface{}, poolConfig.MaxCap),
		lock:   &sync.Mutex{},
	}
	for i := 0; i < poolConfig.InitialCap; i++ {
		c, err := poolConfig.Factory()
		if err != nil {
			err = fmt.Errorf("factory is not able to fill the pool: %s", err)
			if pool != nil {
				pool.ReleaseAll()
			}
			break
		}
		p.conns <- c
	}
	return &p, nil
}

type netPool struct {
	conns  chan interface{}
	lock   *sync.Mutex
	config poolConfig
}

func (p *netPool) Get() (conn interface{}, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for {
		select {
		case conn = <-p.conns:
			if p.config.IsActive(conn) {

				return
			}
			p.config.Release(conn)
		default:
			conn, err = p.config.Factory()
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
	}
}

func (p *netPool) Put(conn interface{}) {
	if conn == nil {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if !p.config.IsActive(conn) {
		p.config.Release(conn)
	}
	select {
	case p.conns <- conn:
	default:
		p.config.Release(conn)
	}
}
func (p *netPool) ReleaseAll() {
	p.lock.Lock()
	defer p.lock.Unlock()
	close(p.conns)
	for c := range p.conns {
		p.config.Release(c)
	}
	p.conns = make(chan interface{}, p.config.InitialCap)

}
func (p *netPool) Len() (length int) {
	return len(p.conns)
}
