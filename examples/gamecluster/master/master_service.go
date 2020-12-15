package master

import (
	"fmt"
	"nano/component"
	"nano/examples/gamecluster/protocol"
	"nano/internal/runtime"
	"nano/session"
	log "github.com/sirupsen/logrus"
	"math"
	"net/http"
	"github.com/rs/cors"
)

type MasterService struct {
	component.Base
	gates	map[string]int
}

func NewMasterService() *MasterService{
	return &MasterService{
		gates: map[string]int{},
	}
}

func (m *MasterService) Init() {
	go m.startup()
}

func (m *MasterService) AfterInit()  {
	session.Lifetime.OnClosed(func(s *session.Session) {
	})
}

func (m *MasterService) startup()  {
	mux := http.NewServeMux()
	mux.HandleFunc("/addr", func(w http.ResponseWriter, r *http.Request) {
		gl, ok :=  runtime.CurrentNode.Handler().RemoteMember("GateService")
		if !ok {
			fmt.Fprintf(w, "")
		}
		num := math.MaxInt64
		addr:= ""
		for _, mi := range gl {
		//	[]*clusterpb.MemberInfo
			if n, _ := m.gates[mi.ClientAddr]; n < num {
				log.Printf("gate addr: %s n: %d", mi.ClientAddr, n)
				addr = mi.ClientAddr
				num = n
			}
		}
		log.Printf("advise gate addr: %s", addr)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("{\"addr\": \"%s\",\"path\": \"%s\"}", addr, "/nano")))
	})
	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(":8090", handler); err != nil {
		panic(err)
	}
}

func (m *MasterService) NewUser(s *session.Session, msg *protocol.NewUserRequest) error {
	log.Printf("MasterService.NewUser user: %+v", msg)
	s.Bind(msg.UId)
	s.Set("name", msg.Name)
	s.Set("addr", msg.Addr)
	m.gates[msg.Addr]++
	return nil
}

func (m *MasterService) userDisconnected(s *session.Session)  {
	log.Println("User session disconnected", s.UID())
	addr := s.String("addr")
	m.gates[addr]--
}