//数据源为etcd时，将storage初始化为该类型
package datasource

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/configcenter/config"
    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/util"
    "github.com/coreos/etcd/clientv3"
    "github.com/spf13/viper"
)

//type EtcdInfo struct {
//    EndPoint         []string `json:"endpoints"`
//    UserName         string   `json:"username"`
//    PassWord         string   `json:"password"`
//    OperationTimeout int      `json:"operationtimeout"`
//}

type EtcdType struct {
    client *clientv3.Client
    kv     clientv3.KV
    // cfgMap *EtcdInfo
}

func NewEtcdType() (*EtcdType, error) {
    instance := new(EtcdType)
    err := instance.ConnectToEtcd()
    if err != nil {
        return nil, err
    }
    return instance, nil
}

// ConnectToEtcd 读取etcd配置文件并初始化clientv3
func (e *EtcdType) ConnectToEtcd() error {
    //e.cfgMap = new(EtcdInfo)
    //file, err := os.Open(etcdConfigLocation)
    //if err != nil {
    //    return err
    //}
    //binaryFlie, err := ioutil.ReadAll(file)
    //if err != nil {
    //    return err
    //}
    //err = json.Unmarshal(binaryFlie, &e.cfgMap)
    //if err != nil {
    //    return err
    //}
    var err error
    e.client, err = clientv3.New(clientv3.Config{
        Endpoints:   viper.GetStringSlice(config.EtcdEndpoints),
        DialTimeout: time.Duration(2*viper.GetInt(config.EtcdOperationTimeout)) * time.Second,
        Username:    viper.GetString(config.EtcdUserName),
        Password:    viper.GetString(config.EtcdPassWord),
    })
    if err != nil {
        return err
    }
    e.kv = clientv3.NewKV(e.client)
    err = e.Put("testconn", "justfortest")
    if err != nil {
        msg := fmt.Sprintf("init etcd err: %v", err)
        log.Sugar().Info(msg)
        return errors.New(msg)
    }
    e.Delete("testconn")
    return nil
}

func (e *EtcdType) Put(key, value string) error {
    ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(viper.GetInt(config.EtcdOperationTimeout))*time.Second)
    _, err := e.kv.Put(ctx, key, value)
    cancel()
    if err != nil {
        return err
    }
    return nil
}

func (e *EtcdType) Get(key string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(viper.GetInt(config.EtcdOperationTimeout))*time.Second)
    resp, err := e.kv.Get(ctx, key)
    cancel()
    if err != nil {
        return nil, err
    }
    if resp.Kvs != nil {
        return resp.Kvs[0].Value, nil
    }
    return nil, nil
}

func (e *EtcdType) Delete(key string) error {
    ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(viper.GetInt(config.EtcdOperationTimeout))*time.Second)
    _, err := e.kv.Delete(ctx, key)
    cancel()
    if err != nil {
        return err
    }
    return nil
}

// GetbyPrefix 范围获取
func (e *EtcdType) GetbyPrefix(prefix string) (map[string][]byte, error) {
    prefix = util.GetPrefix(prefix)
    ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(viper.GetInt(config.EtcdOperationTimeout))*time.Second)
    resp, err := e.kv.Get(ctx, prefix, clientv3.WithPrefix())
    cancel()
    if err != nil {
        return nil, err
    }
    kvs := make(map[string][]byte)
    if resp.Kvs != nil {
        for _, data := range resp.Kvs {
            kvs[string(data.Key)] = data.Value
        }
        return kvs, nil
    }
    return nil, nil
}

// DeletebyPrefix 范围删除
func (e *EtcdType) DeletebyPrefix(prefix string) error {
    prefix = util.GetPrefix(prefix)
    ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(viper.GetInt(config.EtcdOperationTimeout))*time.Second)
    _, err := e.kv.Delete(ctx, prefix, clientv3.WithPrefix())
    cancel()
    if err != nil {
        return err
    }
    return nil
}

func (e *EtcdType) GetSourceDataorOperator() interface{} {
    return e.client
}

func (e *EtcdType) AtomicCommit(putMap map[string]string, deleteSlice []string) error {
    var operationSlice []clientv3.Op
    if putMap != nil && len(putMap) > 0 {
        for k, v := range putMap {
            operationSlice = append(operationSlice, clientv3.OpPut(k, v))
        }
    }
    if deleteSlice != nil && len(deleteSlice) > 0 {
        for _, v := range deleteSlice {
            //消除将要删除对象中可能的重复键值，因为etcd事务不支持键值重复
            if _, ok := putMap[v]; !ok {
                operationSlice = append(operationSlice, clientv3.OpDelete(v))
            }
        }
    }
    if len(operationSlice) == 0 {
        return nil
    }
    ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(3*viper.GetInt(config.EtcdOperationTimeout))*time.Second) //提交时限设为一般时限*提交次数
    txn := e.client.Txn(ctx)
    txn = txn.Then(operationSlice...)
    //重试3次
    var errIns error
    for i := 0; i < 3; i++ {
        _, err := txn.Commit()
        if err != nil {
            errIns = err
            continue
        }
        cancel()
        return nil
    }
    cancel()
    return errIns
}

func (e *EtcdType) GracefullyClose(ctx context.Context) {
    select {
    case <-ctx.Done():
        if e.client != nil {
            _ = e.client.Close()
            defer func() {
                if err := recover(); err != nil {
                    log.Sugar().Errorf("recover from:%v", err)
                    fmt.Println(fmt.Sprintf("recover from:%v", err))
                }
            }()
        }
        fmt.Println("etcd disconnected")
    }
}
