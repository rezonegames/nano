package chat

import (
	"nano/component"
	"nano/session"
)

var (
	// All services in master server
	Services = &component.Components{}

	roomService = newRoomService()
)

func init() {
	Services.Register(roomService)
}

func OnSessionClosed(s *session.Session) {
	roomService.userDisconnected(s)
}
