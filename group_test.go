package nano

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/lonng/nano/session"
)

func TestChannel_Add(t *testing.T) {
	c := NewGroup("test_add")

	var paraCount = 100
	w := make(chan bool, paraCount)
	for i := 0; i < paraCount; i++ {
		go func(id int) {
			s := session.New(nil)
			uid := strconv.Itoa(id + 1)
			s.Bind(uid)
			c.Add(s)
			w <- true
		}(i)
	}

	for i := 0; i < paraCount; i++ {
		<-w
	}

	if c.Count() != paraCount {
		t.Fatalf("count expect: %d, got: %d", paraCount, c.Count())
	}

	n := strconv.Itoa(rand.Intn(paraCount) + 1)
	if !c.Contains(n) {
		t.Fail()
	}

	// leave
	c.LeaveAll()
	if c.Count() != 0 {
		t.Fail()
	}
}
