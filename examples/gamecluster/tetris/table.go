package tetris

import (
	"fmt"
	"github.com/pingcap/errors"
	log "github.com/sirupsen/logrus"
	"nano"
	"nano/component"
	"nano/examples/gamecluster/tetris/protocol"
	"nano/scheduler"
	"nano/session"
	"sync"
	"sync/atomic"
	"time"
)

func SessionTable(s *session.Session) (*Table, error) {
	if !s.HasKey(TableKey) {
		return nil, fmt.Errorf("not join room yet")
	}
	return s.Value(TableKey).(*Table), nil
}

func SessionRound(s *session.Session) (*Round, error) {
	if !s.HasKey(RoundKey) {
		return nil, fmt.Errorf("not join room yet")
	}
	return s.Value(RoundKey).(*Round), nil
}

const (
	TableKey = "table"
	RoundKey = "round"
)

var (
	nextTableId int32
	nextRoundId int32
)

type (
	User struct {
		uId       int64
		ready     bool
		over      bool
		pos       int
		team      int
		isLeave   bool
		entertime int64
	}

	Round struct {
		roundId int32
		rank    []*User
		//0,1,2 red team
		//3,4,5 blue team
		red  map[int64]*User
		blue map[int64]*User

		pending map[int64]*session.Session

		cap int

		chReady chan *session.Session
		chOver  chan *session.Session

		status int

		table *Table

		isRunning bool

		owner int64
	}

	Table struct {
		group   *nano.Group
		cap     int
		tableId int32
		name    string
		desc    string
		chDie   chan struct{}
		round   *Round
		room    *Room
		lock    sync.RWMutex
	}
)

func newTable(opts ...Option) *Table {
	opt := options{}
	for _, option := range opts {
		option(&opt)
	}
	atomic.AddInt32(&nextTableId, 1)
	atomic.AddInt32(&nextRoundId, 1)
	r := &Round{
		roundId: nextRoundId,
		rank:    nil,
		red:     make(map[int64]*User),
		blue:    make(map[int64]*User),
		chReady: make(chan *session.Session),
		chOver:  make(chan *session.Session),
		cap:     opt.cap,
	}
	t := &Table{
		group:   nano.NewGroup(fmt.Sprintf("table-%d", nextTableId)),
		cap:     opt.cap,
		tableId: nextTableId,
		name:    opt.name,
		desc:    opt.desc,
		chDie:   nil,
		round:   r,
		room:    opt.room,
	}
	r.table = t
	go r.run()
	return t
}

func (r *Round) reset() {
	atomic.AddInt32(&nextRoundId, 1)
	r.roundId = nextRoundId
	r.rank = nil

	for k, v := range r.red {
		v.ready = false
		v.over = false
		if v.isLeave {
			delete(r.red, k)
		}
	}
	for k, v := range r.blue {
		v.ready = false
		v.over = false
		if v.isLeave {
			delete(r.blue, k)
		}
	}

	for _, v := range r.pending {
		r.add2recommend(v, true)
	}

	r.isRunning = false

	log.Printf("Round.reset r: %+v", r)
}

func (r *Round) change2pos(s *session.Session, pos int) error {
	if r.isRunning {
		r.pending[s.UID()] = s
		return nil
	}
	if pos >= r.cap/2 {
		if len(r.red) >= r.cap/2 {
			return fmt.Errorf("red full %d", r.roundId)
		}
		r.red[s.UID()] = &User{
			uId:   s.UID(),
			ready: false,
			over:  false,
			pos:   pos,
			team:  0,
		}
	} else {
		if len(r.blue) >= r.cap/2 {
			return fmt.Errorf("blue full %d", r.roundId)
		}
		r.blue[s.UID()] = &User{
			uId:   s.UID(),
			ready: false,
			over:  false,
			pos:   pos,
			team:  1,
		}
	}
	return nil
}

func (r *Round) add2recommend(s *session.Session, bForce bool) error {
	if r.isRunning && !bForce {
		r.pending[s.UID()] = s
		return nil
	}
	rc := len(r.red)
	bc := len(r.blue)
	if rc > bc {
		r.blue[s.UID()] = &User{
			uId:       s.UID(),
			ready:     false,
			over:      false,
			pos:       bc + 1,
			team:      0,
			entertime: time.Now().UnixNano(),
		}
	} else {
		r.red[s.UID()] = &User{
			uId:       s.UID(),
			ready:     false,
			over:      false,
			pos:       rc + 1,
			team:      1,
			entertime: time.Now().UnixNano(),
		}
	}
	s.Set(RoundKey, r)

	if r.owner == 0 {
		r.selectowner()
	}

	return nil
}

func (r *Round) selectowner() {
	oldover := r.owner
	n := time.Now().UnixNano()
	for k, v := range r.red {
		if v.entertime < n {
			n = v.entertime
			r.owner = k
		}
	}
	for k, v := range r.blue {
		if v.entertime < n {
			n = v.entertime
			r.owner = k
		}
	}
	if oldover == r.owner {
		r.owner = 0
	}
}

func (r *Round) user(s *session.Session) *User {
	if u, ok := r.red[s.UID()]; ok {
		return u
	}
	return r.blue[s.UID()]
}

func (r *Round) leave(s *session.Session) error {
	// pending user leave
	if _, ok := r.pending[s.UID()]; ok {
		delete(r.pending, s.UID())
		return nil
	}

	u := r.user(s)
	if r.isRunning {
		u.isLeave = true
		r.chOver <- s
	} else {
		delete(r.blue, s.UID())
		delete(r.red, s.UID())
	}

	if r.owner == s.UID() {
		r.selectowner()
	}

	s.Remove(RoundKey)
	return nil
}

func (r *Round) ready(s *session.Session) error {
	if r.isRunning {
		return fmt.Errorf("round is running")
	}
	//log.Printf("Round.ready before userid: %d", s.UID())
	r.chReady <- s
	//log.Printf("Round.ready after userid: %d", s.UID())

	return nil
}

func (r *Round) cancelready(s *session.Session) error {
	if r.isRunning {
		return fmt.Errorf("round is running")
	}
	u := r.user(s)
	u.ready = false
	//broadcast
	r.table.cancelready(s)
	return nil
}

func (r *Round) over(s *session.Session) error {
	if !r.isRunning {
		return fmt.Errorf("round alredy over")
	}
	r.chOver <- s
	return nil
}

func (r *Round) kick(s *session.Session) error {
	if r.isRunning {
		return fmt.Errorf("round is running")
	}
	return r.leave(s)
}

func (r *Round) run() {
	for {
		select {
		case s := <-r.chReady:
			b := true
			//log.Printf("Round.run ready userid: %d", s.UID())

			blue := []int64{}
			red := []int64{}

			r.table.ready(s)
			for k, v := range r.red {
				if k == s.UID() {
					v.ready = true
				}
				if !v.ready {
					b = false
				}
				red = append(red, k)
			}
			for k, v := range r.blue {
				if k == s.UID() {
					v.ready = true
				}
				if !v.ready {
					b = false
				}
				blue = append(blue, k)
			}
			log.Printf("Round.run ready userid: %d b: %v red: %d blue: %d x: %d", s.UID(), b, len(r.red), len(r.blue), r.cap/2)
			if b && len(r.red) == len(r.blue) && len(r.red) == r.cap/2 {
				r.isRunning = true
				r.table.start(red, blue)
			}
		case s := <-r.chOver:
			b := true
			//log.Printf("Round.run over userid: %d", s.UID())
			r.table.over(s)
			for k, v := range r.red {
				if k == s.UID() {
					v.over = true
					r.rank = append(r.rank, v)
				}
				if !v.over {
					b = false
				}
			}
			for k, v := range r.blue {
				if k == s.UID() {
					v.over = true
					r.rank = append(r.rank, v)
				}
				if !v.over {
					b = false
				}
			}
			log.Printf("Round.run over userid: %d b: %v rank: %d", s.UID(), b, len(r.rank))
			if b {
				r.table.stopandsettle(r.rank)
				r.reset()
			}
		}
	}
}

func (t *Table) start(red, blue []int64) {
	count := 5
	scheduler.NewCountTimer(1*time.Second, count, func() {
		count--
		if count == 0 {
			log.Println("round start :)")
			t.group.Broadcast("onStart", &protocol.OnStart{
				Blue: blue,
				Red:  red,
			})
		} else {
			log.Printf("..............%d", count)
			t.group.Broadcast("onCountdown", &protocol.OnCountdown{Countdown: count})
		}
	})
}

func (t *Table) ready(s *session.Session) {
	//log.Printf("Table.ready userid: %d", s.UID())
	t.group.Broadcast("onReady", &protocol.OnReady{User: protocol.UserInfo{
		UId:     s.UID(),
		Name:    s.String("name"),
		Content: "user ready",
	}})
}

func (t *Table) cancelready(s *session.Session) {
	//log.Printf("Table.ready userid: %d", s.UID())
	t.group.Broadcast("onCancelReady", &protocol.OnCancelReady{User: protocol.UserInfo{
		UId:     s.UID(),
		Name:    s.String("name"),
		Content: "user cancelready",
	}})
}

func (t *Table) over(s *session.Session) {
	//log.Printf("Table.over userid: %d", s.UID())
	t.group.Broadcast("onOver", &protocol.OnOver{User: protocol.UserInfo{
		UId:     s.UID(),
		Name:    s.String("name"),
		Content: "user over",
	}})
}

func (t *Table) stopandsettle(rank []*User) {
	scheduler.NewAfterTimer(3*time.Second, func() {
		rl := []protocol.Reward{}
		wt := rank[len(rank)-1].team
		for _, r := range rank {
			if r.team == wt {
				is := []protocol.Item{}
				is = append(is, protocol.Item{
					ItemType: protocol.Diamond,
					Name:     "diamond",
					Desc:     "钻石",
					Count:    10,
				})
				rw := protocol.Reward{
					Items: is,
					UId:   r.uId,
				}
				rl = append(rl, rw)
			}
		}
		t.group.Broadcast("onStopAndSettle", &protocol.StopAndSettleBroadcast{
			Rewards: rl,
		})
	})
}

func (t *Table) join(s *session.Session) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.group.Count() >= t.cap {
		return fmt.Errorf("table is full %d", t.tableId)
	}
	if err := t.group.Add(s); err != nil {
		return errors.Trace(err)
	}

	s.Set(TableKey, t)
	// ensure success, need test
	t.round.add2recommend(s, false)
	log.Printf("Table.join userid: %d round: %+v", s.UID(), t.round)

	t.group.Broadcast("onJoinTable", &protocol.OnJoinTable{protocol.UserInfo{
		UId:     s.UID(),
		Name:    s.String("name"),
		Content: "user join table",
	}})
	return nil
}

func (t *Table) leave(s *session.Session) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	// leave table
	if err := t.group.Leave(s); err != nil {
		return errors.Trace(err)
	}
	s.Remove(TableKey)

	//leave round
	t.round.leave(s)

	log.Printf("Table.leave userid: %d round: %+v", s.UID(), t.round)

	t.group.Broadcast("onLeaveTable", &protocol.UserInfo{
		UId:     s.UID(),
		Name:    s.String("name"),
		Content: fmt.Sprintf("user: %s leave table: %d", s.String("name"), t.tableId),
	})

	return nil
}

func (t *Table) kick(s *session.Session) error {
	return t.leave(s)
}

func (t *Table) update(msg []byte) error {
	return t.group.Broadcast("onUpdate", msg)
}

func (t *Table) info() *protocol.TableInfo {
	return &protocol.TableInfo{
		TableId:   t.tableId,
		Name:      t.name,
		Desc:      t.desc,
		Owner:     1000,
		OwnerName: "name",
	}
}

//table service
type TableService struct {
	component.Base
	room *Room
	cap  int
}

func NewTableService(opts ...Option) *TableService {
	opt := options{}
	for _, option := range opts {
		option(&opt)
	}
	return &TableService{
		room: opt.room,
		cap:  opt.tablecap,
	}
}

func (t *TableService) AfterInit() {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for range ticker.C {
			for k, v := range t.room.tables {
				log.Printf("table: %d info: %+v round: %+v", k, v.group.Members(), *v.round)
			}
		}
	}()
}

func (ts *TableService) Create(s *session.Session, msg *protocol.CreateTableRequest) error {
	t := newTable(
		WithCap(ts.cap),
		WithName(msg.Name),
		WithDesc(msg.Desc),
		WithRoom(ts.room),
	)
	if err := t.join(s); err != nil {
		return errors.Trace(err)
	}
	// create table set room table
	ts.room.setTable(t.tableId, t)
	return s.Response(&protocol.CreateTableResponse{
		Code:      0,
		TableInfo: *t.info(),
	})
}

func (ts *TableService) Join(s *session.Session, msg *protocol.JoinTableRequest) error {
	t, err := ts.room.getTable(msg.TableId)
	if err != nil {
		return errors.Trace(err)
	}
	if err := t.join(s); err != nil {
		return errors.Trace(err)
	}
	return s.Response(&protocol.JoinTableResponse{
		Code:      0,
		TableInfo: *t.info(),
	})
}

func (ts *TableService) Leave(s *session.Session, msg *protocol.TableLeaveRequest) error {
	t, err := SessionTable(s)
	if err != nil {
		return errors.Trace(err)
	}
	if err := t.leave(s); err != nil {
		return errors.Trace(err)
	}
	return s.Response(&protocol.TableLeaveResponse{})
}

func (ts *TableService) Ready(s *session.Session, msg *protocol.ReadyRequest) error {
	r, err := SessionRound(s)
	if err != nil {
		return errors.Trace(err)
	}
	if err := r.ready(s); err != nil {
		return errors.Trace(err)
	}
	log.Printf("TableService.Ready userid: %d.........", s.UID())
	return s.Response(&protocol.ReadyResponse{})
}

func (ts *TableService) CancelReady(s *session.Session, msg *protocol.CancelReadyRequest) error {
	r, err := SessionRound(s)
	if err != nil {
		return errors.Trace(err)
	}
	if err := r.cancelready(s); err != nil {
		return errors.Trace(err)
	}
	log.Printf("TableService.CancelReady userid: %d.........", s.UID())
	return s.Response(&protocol.CancelReadyResponse{})
}

func (ts *TableService) Over(s *session.Session, msg *protocol.OverRequest) error {
	r, err := SessionRound(s)
	if err != nil {
		return errors.Trace(err)
	}
	r.over(s)
	return s.Response(&protocol.OverResponse{})
}

func (ts *TableService) Update(s *session.Session, msg []byte) error {
	t, err := SessionTable(s)
	if err != nil {
		return errors.Trace(err)
	}
	return t.update(msg)
}
