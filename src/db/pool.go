package db

import (
	"context"
	"errors"
	"sync"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/milvus-io/milvus/client/v2/entity"
    "github.com/milvus-io/milvus/client/v2/index"
	"os"
	"log"
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

var GlobalPool *Pool 

func init(){
	var err error
	GlobalPool,err = Newpool(10)
    if err != nil {
        log.Printf("连接池初始化失败: %v", err)
    } else {
        log.Printf("连接池初始化成功")
    }    
	GlobalPool.initcollections("demo_collection")
}

func Newpool(size int)(*Pool ,error) {
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
	for len(p.connections)>0 {
		cli:=<-p.connections
		err:=p.isHealthy(cli)
		if err==nil{
			return cli,nil
		}
	}
	conn,err:=p.factory()
	if err!=nil {
		return nil,err
	}
	return conn,nil
}

func (p *Pool) isHealthy(cli *milvusclient.Client) error {
    _, err := cli.ListCollections(p.ctx,milvusclient.NewListCollectionOption())
    if err != nil {
        return err
    }
    return nil
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


func factory() (*milvusclient.Client, error) {
    ctx := context.Background()
    client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
        Address:  os.Getenv("MILVUS_ADDRESS"),
    })
    if err != nil {
        return nil, err
    }
    return client, nil
}

func (p *Pool) initcollections(CollectionName string) {
	client,err:=GlobalPool.Get()
	if err!=nil {
		log.Printf("连接获取失败：%v",err)
	}
	defer GlobalPool.Put(client)
	
	has, err := client.HasCollection(p.ctx, milvusclient.NewHasCollectionOption(CollectionName))
	if err!=nil {
		log.Printf("查找collections失败：%v",err)
	}
	if has {
		err = client.DropCollection(p.ctx, milvusclient.NewDropCollectionOption(CollectionName))
		if err != nil {
			log.Println(err.Error())
		}
	}

	schema := entity.NewSchema().
    WithField(entity.NewField().WithName("id").
        WithDataType(entity.FieldTypeInt64).
        WithIsPrimaryKey(true).
        WithIsAutoID(true),
    ).
    WithField(entity.NewField().WithName("user_openid").
        WithDataType(entity.FieldTypeVarChar).
        WithMaxLength(200),
    ).
    WithField(entity.NewField().WithName("role").
        WithDataType(entity.FieldTypeVarChar).
		WithMaxLength(20),
    ).
    WithField(entity.NewField().WithName("text").
        WithDataType(entity.FieldTypeVarChar).
        WithMaxLength(2000).
        WithEnableAnalyzer(true),
    ).
    WithField(entity.NewField().WithName("text_dense_vector").
        WithDataType(entity.FieldTypeFloatVector).
        WithDim(512),
    )

	err = client.CreateCollection(p.ctx, milvusclient.NewCreateCollectionOption(CollectionName, schema))
	if err !=nil {
		log.Printf("创建collections失败：%v",err)
	}
	log.Printf("创建collections成功")

    idx := index.NewIvfFlatIndex(entity.L2, 512)
	option := milvusclient.NewCreateIndexOption(CollectionName, "text_dense_vector", idx)
    _, err = client.CreateIndex(p.ctx, option)
	if err != nil {
		log.Printf("创建索引失败")
    }
	loadTask, err := client.LoadCollection(p.ctx, milvusclient.NewLoadCollectionOption(CollectionName))
	if err != nil {
		log.Println(err.Error())
	}
	err = loadTask.Await(p.ctx)
	if err != nil {
		log.Println(err.Error())
	}
}