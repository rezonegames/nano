package scheduler

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"nano/internal/env"
	"nano/internal/log"
)

const (
	messageQueueBacklog = 1 << 10
	sessionCloseBacklog = 1 << 8
)

// LocalScheduler schedules task to a customized goroutine
type LocalScheduler interface {
	Schedule(Task)
}

type Task func()

type Hook func()

var (
	chDie   = make(chan struct{})
	chExit  = make(chan struct{})
	chTasks = make(chan Task, 1<<8)
	started int32
	closed  int32
)

func try(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Handle message panic: %+v\n%s", err, debug.Stack()))
		}
	}()
	f()
}

func Sched() {
	if atomic.AddInt32(&started, 1) != 1 {
		return
	}

	ticker := time.NewTicker(env.TimerPrecision)
	defer func() {
		ticker.Stop()
		close(chExit)
	}()

	for {
		select {
		case <-ticker.C:
			cron()

		case f := <-chTasks:
			try(f)

		case <-chDie:
			return
		}
	}
}

func Close() {
	if atomic.AddInt32(&closed, 1) != 1 {
		return
	}
	close(chDie)
	<-chExit
	log.Println("Scheduler stopped")
}

func PushTask(task Task) {
	chTasks <- task
}
