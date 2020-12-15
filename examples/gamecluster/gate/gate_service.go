package gate

import (
	"nano/benchmark/testdata"
	"nano/component"
	"nano/examples/gamecluster/db"
	"nano/examples/gamecluster/protocol"
	"nano/internal/runtime"
	"nano/session"
	"github.com/pingcap/errors"
	log "github.com/sirupsen/logrus"
)


type GateService struct {
	component.Base
	users 		map[int64]*db.User
	clientAddr	string
}

func NewGateService() *GateService {
	return &GateService{
		users:	map[int64]*db.User{},
	}
}

func (g *GateService) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		g.userDisconnected(s)
	})
}

func (h *GateService) Ping(s *session.Session, data *testdata.Ping) error {
	return s.Push("Pong", &protocol.Pong{})
}

func (g *GateService) Login(s *session.Session, msg *protocol.LoginRequest) error {
	deviceId := msg.DeviceId
	u, err := db.DB.FindAndCreateUser(deviceId)
	if err != nil {
		return errors.Trace(err)
	}

	log.Printf("GateService.Login user: %+v", u)

	s.Bind(u.UId)
	s.Set("name", u.Name)
	g.users[u.UId] = u

	if g.clientAddr == "" {
		g.clientAddr = runtime.CurrentNode.ClientAddr
	}
	n := &protocol.NewUserRequest{
		UId: u.UId,
		Name: u.Name,
		Addr: g.clientAddr,
	}
	s.RPC("MasterService.NewUser", n)
	s.RPC("RoomService.NewUser", n)

	return s.Response(&protocol.LoginResponse{UId: u.UId, Diamond: u.Diamond, Name: u.Name, Pic: u.Pic})
}

func (g *GateService) userDisconnected(s *session.Session)  {
	uId := s.UID()
	delete(g.users, uId)
	log.Println("User session disconnected", s.UID())
}