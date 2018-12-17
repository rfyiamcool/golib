package discovery

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"time"

	etcdClient "github.com/coreos/etcd/client"
)

const (
	TTL = 10 * time.Second

	KeepAlivePeriod = 3 * time.Second
)

type RegistryClient interface {
	Register() error

	Unregister() error

	ServicesByName(name string) ([]string, error)
}

type Options struct {
	EtcdEndpoints []string
	ServiceName   string
	InstanceName  string
	BaseURL       string

	etcdKey         string
	keepAliveTicker *time.Ticker
	cancel          context.CancelFunc
}

type RegisterCtl struct {
	Options
	etcdKApi etcdClient.KeysAPI
}

func New(config Options) (*RegisterCtl, error) {
	cfg := etcdClient.Config{
		Endpoints:               config.EtcdEndpoints,
		Transport:               etcdClient.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := etcdClient.New(cfg)

	if err != nil {
		return nil, err
	}

	etcdClient := &RegisterCtl{
		config,
		etcdClient.NewKeysAPI(c),
	}
	return etcdClient, nil
}

func (e *RegisterCtl) Register() error {
	e.etcdKey = buildKey(e.ServiceName, e.InstanceName)
	value := dto{
		e.BaseURL,
	}

	val, _ := json.Marshal(value)
	e.keepAliveTicker = time.NewTicker(KeepAlivePeriod)
	ctx, c := context.WithCancel(context.TODO())
	e.cancel = c

	insertFunc := func() error {
		_, err := e.etcdKApi.Set(context.Background(), e.etcdKey, string(val), &etcdClient.SetOptions{
			TTL: TTL,
		})
		return err
	}
	err := insertFunc()
	if err != nil {
		return err
	}

	// Exec the keep alive goroutine
	go func() {
		for {
			select {
			case <-e.keepAliveTicker.C:
				insertFunc()
				log.Printf("Keep alive routine for %s", e.ServiceName)
			case <-ctx.Done():
				log.Printf("Shutdown keep alive routine for %s", e.ServiceName)
				return
			}
		}
	}()
	return nil
}

func (e *RegisterCtl) Unregister() error {
	e.cancel()
	e.keepAliveTicker.Stop()
	_, err := e.etcdKApi.Delete(context.Background(), e.etcdKey, nil)
	return err
}

func (e *RegisterCtl) ServicesByName(name string) ([]string, error) {
	response, err := e.etcdKApi.Get(context.Background(), fmt.Sprintf("/%s", name), nil)
	ipList := make([]string, 0)
	if err != nil {
		return ipList, err
	}

	for _, node := range response.Node.Nodes {
		val := &dto{}
		json.Unmarshal([]byte(node.Value), val)
		ipList = append(ipList, val.BaseURL)
	}
	return ipList, nil
}

type dto struct {
	BaseURL string
}

func buildKey(servicetype, instanceName string) string {
	return fmt.Sprintf("%s/%s", servicetype, instanceName)
}
