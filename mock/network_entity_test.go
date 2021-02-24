package mock_test

import (
	"testing"

	. "github.com/pingcap/check"
	"nano/mock"
)

type networkEntitySuite struct{}

func TestNetworkEntity(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&networkEntitySuite{})

func (s *networkEntitySuite) TestNetworkEntity(c *C) {
	entity := mock.NewNetworkEntity()

	c.Assert(entity.LastResponse(), IsNil)
	c.Assert(entity.LastMid(), Equals, uint64(1))
	c.Assert(entity.Response("hello"), IsNil)
	c.Assert(entity.LastResponse().(string), Equals, "hello")

	c.Assert(entity.FindResponseByMID(1), IsNil)
	c.Assert(entity.ResponseMid(1, "test"), IsNil)
	c.Assert(entity.FindResponseByMID(1).(string), Equals, "test")

	c.Assert(entity.FindResponseByRoute("t.tt"), IsNil)
	c.Assert(entity.Push("t.tt", "test"), IsNil)
	c.Assert(entity.FindResponseByRoute("t.tt").(string), Equals, "test")

	c.Assert(entity.RemoteAddr().String(), Equals, "mock-addr")
	c.Assert(entity.Close(), IsNil)
}
