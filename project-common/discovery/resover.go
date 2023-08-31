package discovery

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"time"
)

const (
	schema = "etcd"
)

// Resolver for grpc client
type Resolver struct {
	schema      string
	EtcdAddrs   []string
	DialTimeout int

	closeCh      chan struct{}
	watchCh      clientv3.WatchChan
	cli          *clientv3.Client
	keyPrifix    string
	srvAddrsList []resolver.Address

	cc     resolver.ClientConn
	logger *zap.Logger
}

// NewResolver create a new resolver.Builder base on etcd
func NewResolver(etcdAddrs []string, logger *zap.Logger) *Resolver {
	return &Resolver{
		schema:      schema,
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Scheme returns the scheme supported by this resolver.
func (r *Resolver) Scheme() string { //返回解析器支持的协议方案
	//在目标地址字符串中，协议方案会跟随解析器名并用 "://" 分隔。例如，etcd://localhost:2379 中的 etcd 就是协议方案。
	return r.schema
}

// Build creates a new resolver.Resolver for the given target
// 实现build(命名解析器构建器)接口，返回resolver接口(命名解析器)

// 在 gRPC 调用 grpc.Dial 函数时，会自动将创建的 resolver.ClientConn 实例传递给解析器构建器的 Build 方法
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 而build函数中的target参数，就是grpc目标地址中“etcd”协议方案的后一部分
	//对于请求目标地址，"etcd://localhost:2379/myservice"，target.URL.Path 就是 /myservice。
	//target.URL.Host就是 localhost:2379。

	r.cc = cc
	r.keyPrifix = BuildPrefix(Server{Name: target.URL.Path, Version: target.URL.Host}) //填充keyPrifix字段：服务名-服务版本
	fmt.Println("target.URL.Path:", target.URL.Path+" "+"target.URL.HOST:", target.URL.Host)
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

// ResolveNow resolver.Resolver接口，是干什么的？
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close resolver.Resolver interface
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

// start
func (r *Resolver) start() (chan<- struct{}, error) { //连接etcd，查询etcd中所有r.keyPrefix前缀的值
	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	resolver.Register(r) //注册自定义的解析器构建器（builder），也就是注册【服务发现】的构造器。【传入了实现了builder接口的结构体地址】
	r.closeCh = make(chan struct{})

	if err = r.sync(); err != nil {
		return nil, err
	}

	go r.watch() //服务监控

	return r.closeCh, nil
}

// watch update events
func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Minute)
	r.watchCh = r.cli.Watch(context.Background(), r.keyPrifix, clientv3.WithPrefix())
	//↑watch函数为etcd监听变更的命令
	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh: //一旦监听到变更事件就将事件对象发送到通道，可以通过通道获取信息
			if ok {
				r.update(res.Events) //如果发现变更，就手动updateState更新到grpc里面
			}
		case <-ticker.C: //不断获取，刷新信息
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed", zap.Error(err))
			}
		}
	}
}

// update
func (r *Resolver) update(events []*clientv3.Event) { //传入watch函数返回的变更事件
	for _, ev := range events {
		var info Server
		var err error

		switch ev.Type {
		case mvccpb.PUT:
			info, err = ParseValue(ev.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
			if !Exist(r.srvAddrsList, addr) {
				r.srvAddrsList = append(r.srvAddrsList, addr)
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
			}
		case mvccpb.DELETE:
			info, err = SplitPath(string(ev.Kv.Key))
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr}
			if s, ok := Remove(r.srvAddrsList, addr); ok {
				r.srvAddrsList = s
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
			}
		}
	}
}

// sync 同步获取所有地址信息
func (r *Resolver) sync() error { //发现服务
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := r.cli.Get(ctx, r.keyPrifix, clientv3.WithPrefix()) //查询etcd中的键值对，查询有前缀的所有键值对，如果不存在，返回的是一个空列表而不是错误
	if err != nil {
		return err
	}
	//fmt.Println("前缀：", r.keyPrifix)

	if len(res.Kvs) == 0 {
		fmt.Println("没有找到服务")
		return errors.New("没有找到服务")
	} /*else {输出发现的服务
		for _, kv := range res.Kvs {
			fmt.Printf(" 键%s——值%s\n", kv.Key, kv.Value)
		}
	}*/

	r.srvAddrsList = []resolver.Address{}
	for _, v := range res.Kvs {
		info, err := ParseValue(v.Value)
		if err != nil {
			continue
		}
		addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
		r.srvAddrsList = append(r.srvAddrsList, addr) //添加所有服务
	}

	r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList}) //将服务真正更新进grpc
	return nil
}
