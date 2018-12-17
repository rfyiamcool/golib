package discovery

import (
	"log"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
)

func TestBasic(t *testing.T) {

	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 5 * time.Second,
	}
	c, err := client.New(cfg)

	if err != nil {
		log.Panicln(err)
	}

	kapi := client.NewKeysAPI(c)

	client := RegisterCtl{
		Options{
			ServiceName:  "test",
			InstanceName: "test1",
			BaseURL:      "127.0.0.1:8080",
		},
		kapi,
	}

	client.Register()
	response, _ := client.ServicesByName("test")
	if len(response) == 0 {
		t.Error("No service registered")
	}
	client.UnregisterAndDelete()
	response, _ = client.ServicesByName("test")
	if len(response) != 0 {
		t.Error("Service not  unregistered")
	}
}

func TestKeepAlive(t *testing.T) {
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 5 * time.Second,
	}
	c, err := client.New(cfg)

	if err != nil {
		log.Panicln(err)
	}

	client1 := RegisterCtl{
		Options{
			ServiceName:  "test",
			InstanceName: "test1",
			BaseURL:      "127.0.0.1:8080",
		},
		client.NewKeysAPI(c),
	}
	client1.Register()

	client2 := RegisterCtl{
		Options{
			ServiceName:  "test",
			InstanceName: "test2",
			BaseURL:      "192.168.1.199:3738",
		},
		client.NewKeysAPI(c),
	}
	client2.Register()

	time.Sleep(3 * time.Second)
	response, _ := client1.ServicesByName("test")
	log.Println(response)
}
