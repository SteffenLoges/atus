package scheduler

import (
	"context"
	"time"

	"github.com/tevino/abool/v2"
)

type Task func(context.Context)

type Scheduler struct {
	closeChan  chan bool
	task       Task
	interval   time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
	isStopping *abool.AtomicBool
}

// Simple task scheduler
// runs the given task every interval
// will, by design, stall if the task takes longer than the interval
func New(interval time.Duration, task Task) *Scheduler {

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		closeChan:  make(chan bool),
		task:       task,
		interval:   interval,
		ctx:        ctx,
		cancel:     cancel,
		isStopping: abool.NewBool(false),
	}

}

func (w *Scheduler) Run(immediate bool) {
	go func() {
		defer func() {
			w.cancel()
			w = nil
		}()

		if immediate {
			w.task(w.ctx)
		}

		ticker := time.NewTicker(w.interval)
		for {
			select {
			case <-ticker.C:
				w.task(w.ctx)
			case <-w.closeChan:
				return
			}
		}
	}()
}

func (w *Scheduler) Stop() {
	if w.cancel != nil {
		w.cancel()
	}

	// BugFix: closeChan will block if the task is currently running
	if w.isStopping.IsSet() {
		return
	}
	w.isStopping.Set()

	w.closeChan <- true
}
