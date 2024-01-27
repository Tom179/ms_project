package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Register for grpc server
type Register struct { //自定义注册器，注册服务
	EtcdAddrs   []string //etcd地址
	DialTimeout int      //

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo Server
	srvTTL  int64
	cli     *clientv3.Client
	logger  *zap.Logger
}

// NewRegister create a register base on etcd
func NewRegister(etcdAddrs []string, logger *zap.Logger) *Register { //传入etcd地址，建立注册器
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Register a service
func (r *Register) RegistService(srvInfo Server, ttl int64) (chan<- struct{}, error) { //传入要注册的服务
	var err error

	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}

	if r.cli, err = clientv3.New(clientv3.Config{ //连接etcd
		Endpoints:   r.EtcdAddrs, //获取etcd地址（刚刚new的）
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	fmt.Println("r.servInfo为:", r.srvInfo)
	r.srvTTL = ttl

	if err = r.register(); err != nil { //注册
		return nil, err
	}

	r.closeCh = make(chan struct{}) //创建通道

	go r.KeepMyServiceAlive() //监测通道，如果发现异常就手动注册。有关闭信号就注销服务
	return r.closeCh, nil
}

// Stop stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// register 注册节点
func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second) //创建上下文
	defer cancel()

	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL) //分配租约，指定时长为ttl
	if err != nil {
		return err
	}
	r.leasesID = leaseResp.ID
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil { //返回通道，在通道上定期发送续约响应
		return err
	}

	data, err := json.Marshal(r.srvInfo) //存入服务【Server结构体】的字节切片
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	//存入etcd的key为服务名/版本+地址，value为服务对象的json信息
	//Put 方法将服务实例的信息作为键值对存储在 etcd中
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegPath(r.srvInfo))
	return err
}

// keepAlive
func (r *Register) KeepMyServiceAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	//定时器创建后，您可以通过其内置的 C 通道（ticker.C）来接收定时触发的事件。
	//每当定时器的间隔时间到达时，C 通道会收到一个信号，您可以从中读取。
	for {
		select {
		case <-r.closeCh: //register的关闭管道，有值就注销服务
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		case res := <-r.keepAliveCh: //续约信道，续约失败会返回nil
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
			//fmt.Println("续约,存活时间为:", res.TTL)
		case <-ticker.C: //定时器管道
			if r.keepAliveCh == nil { //当续约失败的时候，手动重新注册节点
				if err := r.register(); err != nil { //租约失败
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		}
	}
}

// UpdateHandler return http handler
func (r *Register) UpdateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wi := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}
			_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("update server weight success"))
	})
}

func (r *Register) GetServerInfo() (Server, error) { //获取一个服务
	resp, err := r.cli.Get(context.Background(), BuildRegPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}
	info := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil { //从一堆服务实例中选取一个，解析出来
			return info, err
		}
	}
	return info, nil
}
