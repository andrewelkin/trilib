package utils

import (
	"context"
	"time"
)

type PeriodicalFunction struct {
	ctx            context.Context
	done           chan struct{}
	expiredF       func(string)
	doneF          func(string)
	messageExpired string
	messageDone    string
	t              int64
}

func NewPeriodicalFunction(ctx context.Context, timeout int64, messageExpired string, expiredF func(string)) *PeriodicalFunction {
	wd := &PeriodicalFunction{
		ctx:            ctx,
		done:           make(chan struct{}),
		expiredF:       expiredF,
		messageExpired: messageExpired,
		t:              timeout,
	}
	go wd.worker()
	return wd
}

func NewPeriodicalFunctionFull(ctx context.Context, timeout int64, messageExpired string, messageDone string, expiredF func(string), doneF func(string)) *PeriodicalFunction {
	wd := &PeriodicalFunction{
		ctx:            ctx,
		done:           make(chan struct{}),
		expiredF:       expiredF,
		doneF:          doneF,
		messageExpired: messageExpired,
		messageDone:    messageDone,
		t:              timeout,
	}
	go wd.worker()
	return wd
}

func (wd *PeriodicalFunction) Stop() {
	wd.done <- struct{}{}
}

func (wd *PeriodicalFunction) worker() {
	ticker := time.NewTicker(time.Duration(wd.t) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-wd.ctx.Done(): // user cancellation
			if wd.doneF != nil {
				wd.doneF(wd.messageDone)
			}
			return
		case <-ticker.C: // expired
			if wd.expiredF != nil {
				wd.expiredF(wd.messageExpired)
			}
		case <-wd.done: // full stop
			if wd.doneF != nil {
				wd.doneF(wd.messageDone)
			}
			return
		}
	}
}
