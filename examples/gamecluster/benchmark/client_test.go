package benchmark_test

import (
	"encoding/json"
	"fmt"
	"nano/examples/gamecluster/benchmark"
	"nano/examples/gamecluster/protocol"
	protocol2 "nano/examples/gamecluster/tetris/protocol"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

var (
	addr = "127.0.0.1:34560"
	sg = make(chan os.Signal)
	chClient = make(chan struct{})
	chWait = make(chan struct{})
)

func start(connect *benchmark.Connector)  {

	connect.Request("TableService.Update", &protocol2.UpdateMessage{
		ID: 0,
		X:  0,
		Y:  0,
	}, func(data interface{}) {
		log.Printf("TableService.Update")
	})

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(60)
	log.Printf("before sleep %d", n)
	time.Sleep(time.Duration(n) *time.Second)
	log.Printf("after sleep")

	connect.Request("TableService.Over",  protocol2.OverRequest{}, func(data interface{}) {
		log.Printf("TableService.Over")
	})

	log.Printf("before reset")
	time.Sleep(60 *time.Second)
	log.Printf("after reset")

	connect.Request("TableService.Ready", &protocol2.ReadyRequest{}, func(data interface{}) {
		r := &protocol2.ReadyResponse{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("TableService.Ready")
	})

}

func clientss(deviceId string, b int)  {
	connect := benchmark.NewConnector()

	connect.OnConnected(func() {
		chWait<- struct{}{}
	})

	if err := connect.Start(addr); err != nil {
		log.Printf("%+v", err)
	}
	<-chWait


	// broadcatst

	connect.On("OnCancelQuickStart", func(data interface{}) {
		r := &protocol2.OnCancelQuickStart{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("OnCancelQuickStart %+v", r)
	})

	connect.On("onUpdate", func(data interface{}) {
		r := &protocol2.UpdateMessage{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onUpdate %+v", r)
	})

	connect.On("onOver", func(data interface{}) {
		r := &protocol2.OnOver{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onOver %+v", r)
	})

	connect.On("onStopAndSettle", func(data interface{}) {
		r := &protocol2.StopAndSettleBroadcast{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onStopAndSettle %+v", r)
	})

	connect.On("onJoinRoom", func(data interface{}) {
		r := &protocol2.UserInfo{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onJoinRoom %+v", r)
	})

	connect.On("onLeaveRoom", func(data interface{}) {
		r := &protocol2.UserInfo{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onLeaveRoom %+v", r)
	})


	connect.On("onJoinTable", func(data interface{}) {
		r := &protocol2.OnJoinTable{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onJoinTable %v", r)
	})

	connect.On("onLeaveTable", func(data interface{}) {
		r := &protocol2.UserInfo{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onLeaveTable %+v", r)
	})

	connect.On("onReady", func(data interface{}) {
		r := &protocol2.OnReady{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onReady %+v", r)
	})

	connect.On("onCountdown", func(data interface{}) {
		r := &protocol2.OnCountdown{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onCountdown %+v", r)
	})

	connect.On("onStart", func(data interface{}) {
		r := &protocol2.OnStart{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("onStart %+v", r)

		//go start(connect)

	})

	connect.On("onReady", func(data interface{}) {
		r := &protocol2.OnReady{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("OnReady %+v", r)
	})

	connect.On("Pong", func(data interface{}) {})


	// request
	connect.Request("GateService.Login", &protocol.LoginRequest{DeviceId: deviceId}, func(data interface{}) {
		r := &protocol.LoginResponse{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("GateService.Login %+v", r)
	})

	time.Sleep(1 * time.Second)


	if b== 1 || b== 2 {
		connect.Request("RoomService.Join", &protocol2.JoinRoomRequest{Offset: 10}, func(data interface{}) {
			r := &protocol2.JoinRoomResponse{}
			json.Unmarshal(data.([]byte), r)
			log.Printf("RoomService.Join %v", r)
		})
	}

	if b == 1{

		connect.Request("TableService.Create", &protocol2.CreateTableRequest{
			Name: "我的第一个桌子",
			Desc: "",
		}, func(data interface{}) {
			r := &protocol2.CreateTableResponse{}
			json.Unmarshal(data.([]byte), r)
			log.Printf("TableService.Create %+v", r)
		})
	}
	if b == 2 {
		connect.Request("TableService.Join", &protocol2.JoinTableRequest{
			TableId: 1,
		}, func(data interface{}) {
			r := &protocol2.JoinTableResponse{}
			json.Unmarshal(data.([]byte), r)
			log.Printf("TableService.Join %+v", r)
		})
	}

	if b == 3 {
		connect.Request("RoomService.QuickStart", &protocol2.QuickStartRequest{}, func(data interface{}) {
			r := &protocol2.QuickStartResponse{}
			json.Unmarshal(data.([]byte), r)
			log.Printf("RoomService.QuickStart %+v", r)
		})
	}


	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	log.Printf("before sleep %d", n)
	time.Sleep(time.Duration(n) *time.Second)
	log.Printf("after sleep")
	connect.Request("TableService.Ready", &protocol2.ReadyRequest{}, func(data interface{}) {
		r := &protocol2.ReadyResponse{}
		json.Unmarshal(data.([]byte), r)
		log.Printf("TableService.Ready")
	})


	for {
		connect.Notify("GateService.Ping", &protocol.Ping{})
		time.Sleep(2 * time.Second)
	}
}

func TestCreate(t *testing.T)  {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	go clientss("xxxxxxx", 1)

	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	<-sg
}


func TestJoin(t *testing.T)  {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	corn := 4
	for i :=0; i<corn; i++ {
		deviceId := fmt.Sprintf("xxx-%d", i)
		go clientss(deviceId, 2)
	}

	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	<-sg

}

func TestQuickstart(t *testing.T)  {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	corn := 1
	for i :=0; i<corn; i++ {
		deviceId := fmt.Sprintf("qqq-%d", i)
		go clientss(deviceId, 3)
	}

	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	<-sg
}