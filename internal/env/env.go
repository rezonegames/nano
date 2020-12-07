// env represents the environment of the current process, includes
// work path and config path etc.
package env

import (
	"net/http"
	"time"

	"nano/serialize"
	"nano/serialize/protobuf"
	"google.golang.org/grpc"
)

var (
	Wd          string                   // working path
	Die         chan bool                // wait for end application
	Heartbeat   time.Duration            // Heartbeat internal
	CheckOrigin func(*http.Request) bool // check origin when websocket enabled
	Debug       bool                     // enable Debug
	WSPath      string                   // WebSocket path(eg: ws://127.0.0.1/WSPath)

	// timerPrecision indicates the precision of timer, default is time.Second
	TimerPrecision = time.Second

	// globalTicker represents global ticker that all cron job will be executed
	// in globalTicker.
	GlobalTicker *time.Ticker

	Serializer serialize.Serializer

	GrpcOptions = []grpc.DialOption{grpc.WithInsecure()}
)

func init() {
	Die = make(chan bool)
	Heartbeat = 30 * time.Second
	Debug = false
	CheckOrigin = func(_ *http.Request) bool { return true }
	Serializer = protobuf.NewSerializer()
}
