package tetris

import (
	"container/list"
	"fmt"
	"github.com/pingcap/errors"
	log "github.com/sirupsen/logrus"
	"nano"
	"nano/component"
	protocol2 "nano/examples/gamecluster/protocol"
	"nano/examples/gamecluster/tetris/protocol"
	"nano/session"
	"sync"
	"time"
)

type (
	Room struct {
		group  *nano.Group
		tables map[int32]*Table
		lock   sync.RWMutex
		cap    int
		//quickstart
		ql *list.List

		tablecap int
		//filter 		nano.SessionFilter
	}
)

func (r *Room) run() {
	// quick start
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			r.lock.Lock()

			k := 1

		EXIT:
			//log.Printf("Room.run ql.len: %d k:%d", r.ql.Len(), k)

			if r.ql.Len() != 0 {
				for _, v := range r.tables {
					if v.round.isRunning {
						continue
					}
					if r.ql.Len() <= 0 {
						break
					}

					lack := v.cap - v.group.Count()

					if lack == k {
						log.Printf("table state user count:%d round: %+v", v.group.Count(), v.round)
						for i := 0; i < lack; i++ {
							e := r.ql.Front()
							if e == nil {
								break
							}
							s := e.Value.(*session.Session)
							r.quickstart(s, v)
							r.ql.Remove(e)

						}
					}
				}

				if k <= r.tablecap {
					k++
					goto EXIT
				}
			}

			r.lock.Unlock()

		}
	}()
}

func (r *Room) getTable(tableId int32) (*Table, error) {
	if t, ok := r.tables[tableId]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("no table %d", tableId)
}

func (r *Room) setTable(tableId int32, table *Table) {
	r.tables[tableId] = table
	log.Printf("Room.setTable: tableid %d", table.tableId)
}

func (r *Room) quickstart(s *session.Session, t *Table) error {
	b := true
	if err := r.join(s); err == nil {
		if err := t.join(s); err != nil {
			// rollback join
			r.group.Leave(s)
		} else {
			b = false
		}
	}
	if b {
		s.Push("OnCancelQuickStart", &protocol.OnCancelQuickStart{})
	}
	return nil
}

func (r *Room) pushquickstart(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.ql.PushBack(s)
	return nil
}

func (r *Room) cancelquickstart(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for e := r.ql.Front(); e != nil; e = e.Next() {
		sss := e.Value.(*session.Session)
		if sss.UID() == s.UID() {
			r.ql.Remove(e)
			log.Printf("cancelquickstart userid: %d", s.UID())
			break
		}
	}

	return nil
}

func (r *Room) join(s *session.Session) error {
	// !important lock
	if r.group.Count() >= r.cap {
		return errors.Errorf("room is full")
	}
	r.group.Add(s)
	r.group.Broadcast("onJoinRoom", &protocol.UserInfo{
		UId:     s.UID(),
		Name:    s.String("name"),
		Content: fmt.Sprintf("welcome user %s", s.String("name")),
	})
	return nil
}

func (r *Room) leave(s *session.Session) error {
	// leave table
	if t, _ := SessionTable(s); t != nil {
		if err := t.leave(s); err != nil {
			return errors.Trace(err)
		}
	}
	//
	// leave room
	if err := r.group.Leave(s); err != nil {
		return errors.Trace(err)
	} else {
		r.group.Broadcast("onLeaveRoom", &protocol.UserInfo{
			UId:     s.UID(),
			Name:    s.String("name"),
			Content: fmt.Sprintf(" user: %s leave room", s.String("name")),
		})
	}
	return nil
}

func (r *Room) info() *protocol.JoinRoomResponse {
	tl := []protocol.TableInfo{}
	for _, v := range r.tables {
		tl = append(tl, *v.info())
	}
	return &protocol.JoinRoomResponse{
		Tables: tl,
	}
}

//room service

type RoomService struct {
	component.Base
	room *Room
}

func NewRoomService(opts ...Option) *RoomService {
	opt := options{}
	for _, option := range opts {
		option(&opt)
	}
	r := &Room{
		tables:   map[int32]*Table{},
		group:    nano.NewGroup("room"),
		ql:       list.New(),
		cap:      opt.cap,
		tablecap: opt.tablecap,
	}
	go r.run()
	return &RoomService{
		room: r,
	}
}

func (rs *RoomService) Room() *Room {
	return rs.room
}

func (rs *RoomService) AfterInit() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			log.Printf("room table count: %d user count: %d", len(rs.room.tables), rs.room.group.Count())
		}
	}()
	session.Lifetime.OnClosed(rs.userDisconnected)
}

func (rs *RoomService) QuickStart(s *session.Session, msg *protocol.QuickStartRequest) error {
	rs.room.pushquickstart(s)
	log.Printf("QuickStart userid: %d", s.UID())
	return s.Response(&protocol.QuickStartResponse{})
}

func (rs *RoomService) CancelQuickStart(s *session.Session, msg *protocol.CancelQuickStartRequest) error {
	rs.room.cancelquickstart(s)
	return s.Response(&protocol.CancelQuickStartResponse{})
}

func (rs *RoomService) Join(s *session.Session, msg *protocol.JoinRoomRequest) error {
	if err := rs.room.join(s); err != nil {
		return errors.Trace(err)
	}
	tl := rs.room.info()
	log.Printf("RoomService.Join %+v", tl)
	return s.Response(tl)
}

func (rs *RoomService) Leave(s *session.Session, msg *protocol.LeaveRoomRequest) error {
	if err := rs.room.leave(s); err != nil {
		return errors.Trace(err)
	}
	return s.Response(&protocol.LeaveRoomResponse{
		Code:    0,
		UId:     s.UID(),
		Content: fmt.Sprintf("player %d leave room:", s.UID()),
	})
}

func (rs *RoomService) userDisconnected(s *session.Session) {
	rs.room.leave(s)
}

func (rs *RoomService) NewUser(s *session.Session, msg *protocol2.NewUserRequest) error {
	s.Bind(msg.UId)
	s.Set("name", msg.Name)
	log.Printf("RoomService.NewUser user: %+v", msg)
	return nil
}
