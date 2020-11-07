package main

import (
	"errors"
	"io"
	"sync"
	"time"
)

/**
实现通用模拟连接池的效果,通过channel的阻塞获取实现
如果使用slice或者其他方式实现需要设置一个获取连接池的队列
*/


var (
	ErrInvalidConfig = errors.New("invalid pool config")
	ErrPoolClosed = errors.New("pool closed")
)

type factory func() (io.Closer, error)


type Pool interface {
	Acquire() (io.Closer, error)
	Release(io.Closer) error
	Close(io.Closer) error
	ShutDown() error
}

type GenericPool struct{
	sync.Mutex
	pool           chan io.Closer
	maxOpen        int
	numOpen        int
	minOpen        int
	closed         bool
	maxLifetime    time.Duration
	factory        factory

}


func NewGenericPool(minOpen, maxOpen int, maxLifetime time.Duration, factory factory) (*GenericPool, error){
	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}

	p := &GenericPool{
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		maxLifetime: maxLifetime,
		factory:     factory,
		pool:        make(chan io.Closer, maxOpen),
	}


	for i := 0; i < minOpen; i ++{
		closed, err := factory()
		if err != nil {
			continue
		}
		p.numOpen ++
		p.pool <- closed
	}
	return p, nil
}



func (p *GenericPool) Acquire() (io.Closer, error){
	if p.closed {
		return nil, ErrPoolClosed
	}

	for {
		closer, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}
		// 需要factory里实现检测超时的方法，此处调用方法传入maxLifeTime验证是否超时
		return closer, nil
	}
}


func (p *GenericPool) getOrCreate() (io.Closer, error){
	select {
	case closer := <- p.pool:
		return closer, nil
	default:
	}

	p.Lock()
	if p.numOpen >= p.maxOpen{
		closer := <- p.pool
		p.Unlock()
		return closer, nil
	}

	//新建连接
	closer, err := p.factory()
	if err != nil {
		return nil, err
	}
	p.numOpen++
	p.Unlock()
	return closer, nil
}


//释放单个资源
func (p *GenericPool) Release(closer io.Closer) error{
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	p.pool <- closer
	p.Unlock()
	return nil
}

//关闭单个资源

func (p *GenericPool) Close(closer io.Closer) error {
	p.Lock()
	closer.Close()
	p.numOpen--
	p.Unlock()
	return nil
}

//关闭整个资源
func (p *GenericPool) ShutDown() error{
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	close(p.pool)
	for closer := range p.pool {
		closer.Close()
		p.numOpen--
	}
	p.closed = true
	p.Unlock()
	return nil

}