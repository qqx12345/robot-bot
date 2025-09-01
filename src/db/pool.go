package db

import (
	"context"
	"errors"
	"sync"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

//milvus连接池
type Pool struct {
	connections chan *milvusclient.Client
	factory func()(*milvusclient.Client,error)
	close bool
	size int
	ctx context.Context
	mu sync.Mutex
}

func Newpool(factory func() (*milvusclient.Client,error), size int)(*Pool ,error) {
	if size<=0 {
		return nil, errors.New("invalid size")
	}
	ctx:=context.Background()
	p:=&Pool{
		connections: make(chan *milvusclient.Client,size),
		factory: factory,
		size: size,
		ctx: ctx,
	}
	for i:=0; i<size; i++ {
		conn,err:=p.factory()
		if err!=nil {
			p.Close()
			return nil,err
		}
		p.connections<-conn
	}
	return p,nil
}

func (p *Pool) Close() {
    p.mu.Lock()
	defer p.mu.Unlock()
    if p.close {
        return
    }
    p.close = true
    close(p.connections)

    for conn := range p.connections {
        if conn != nil {
            conn.Close(p.ctx)
        }
    }
}

func (p *Pool) Get() (*milvusclient.Client,error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.close {
		return nil,errors.New("pool is closed")
	}
	select {
	    case conn := <-p.connections:
			return conn,nil
		default:
			conn,err:=p.factory()
			if err!=nil {
				return nil,err
			}
			return conn,nil
	}
}

func (p *Pool) Put(conn *milvusclient.Client) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.close {
		return conn.Close(p.ctx)
	}
	select {
	case p.connections <- conn:
		return nil
	default:
		return conn.Close(p.ctx)
	}
}