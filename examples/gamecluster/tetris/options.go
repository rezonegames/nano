package tetris

import "nano"

type (
	options struct {
		name	string
		desc	string
		cap		int
		tablecap	int
		room 	*Room
		filter 	nano.SessionFilter
	}

	// Option used to customize handler
	Option func(options *options)
)

func WithName(name string) Option {
	return func(opt *options) {
		opt.name = name
	}
}

func WithDesc(desc string) Option {
	return func(opt *options) {
		opt.desc = desc
	}
}

func WithCap(cap int) Option {
	return func(opt *options) {
		opt.cap = cap
	}
}

func WithTableCap(tablecap int) Option {
	return func(opt *options) {
		opt.tablecap = tablecap
	}
}

func WithRoom(room *Room) Option {
	return func(opt *options) {
		opt.room = room
	}
}

func WithSessionFilter(filter nano.SessionFilter) Option {
	return func(opt *options) {
		opt.filter = filter
	}
}