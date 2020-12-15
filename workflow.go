// workflow.go
package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type FlowClient interface {
	Name() string
	Len() int
	Phase() string
	Start()
}

type Error struct {
	index int
	error
}

type WorkFlow struct {
	name  string
	fId   string
	phase int
	ws    []*Worker
	in    map[*Worker]bool
	s     State
	c     Context
	err   chan Error
}

func (wf WorkFlow) len() int {
	return len(wf.ws)
}

func newWorkFlow(name string) *WorkFlow {
	wf := &WorkFlow{
		name: name,
		fId:  "WF-" + time.Now().Format("20060102150405.000"),
		ws:   make([]*Worker, 0, 10),
		s:    STATE_PEND,
		in:   make(map[*Worker]bool, 10),
	}
	wf.c.Context, wf.c.CancelFunc = context.WithCancel(context.Background())

	return wf
}

func (wf *WorkFlow) Start() {
	go func() {
		wf.load()
	LOOP:
		wg := sync.WaitGroup{}
		wf.s = STATE_RUN
		wf.err = make(chan Error, 10)
		for i, w := range wf.ws {
			if w.state.isExit() {
				if w.FailureEnd &&
					w.failureTimes > w.RepeatTimes {
					break
				}
				if w.RepeatTimes == 0 {
					continue
				}
			} else if w.state.isDone() {
				continue
			}

			err := wf.doRunWorker(&wg, i, w)
			if err != nil &&
				w.state.isExit() &&
				w.FailureEnd {
				break
			}
		}

		wg.Wait()
		close(wf.err)
		for e := range wf.err {
			fmt.Println(e.index, e.error)
			if e.error != nil && e.index < wf.len() {
				if wf.Workers()[e.index].IsFailureRepeat() {
					fmt.Println("goto LOOP")
					goto LOOP
				}
			}
		}

		wf.s = STATE_DONE
	}()
}

func (wf WorkFlow) doRunWorker(wg *sync.WaitGroup, i int, w *Worker) error {
	wg.Add(1)
	wf.phase = i + 1
	fmt.Println(w.Name, ":start")
	err := w.start(wg)
	if err != nil &&
		w.state.isExit() &&
		w.IsFailureRepeat() {
		return wf.doRunWorker(wg, i, w)
	}
	//fmt.Println("failureTimes:", w.failureTimes)
	return err
}

func (wf WorkFlow) Workers() []*Worker {
	return wf.ws
}

func (wf WorkFlow) Cancel() {
	wf.c.CancelFunc()
}

func (wf WorkFlow) Phase() int {
	return wf.phase
}

func (wf *WorkFlow) AddWorker(w *Worker) {
	for _, w1 := range w.DependsOn {
		if w1 == nil || w1 == w {
			continue
		}
		wf.AddWorker(w1)
	}

	if !wf.in[w] {
		w.init(wf)
		wf.in[w] = true
		wf.ws = append(wf.ws, w)
	}
}

func (wf *WorkFlow) load() {

}
