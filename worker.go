// job.go
package workflow

import (
	"fmt"
	"sync"
	"time"
)

type State int

const (
	STATE_PEND = iota
	STATE_RUN
	STATE_STOP
	STATE_EXIT
	STATE_DONE
)

func (s State) isPend() bool {
	return s == STATE_PEND
}
func (s State) isRun() bool {
	return s == STATE_RUN
}
func (s State) isStop() bool {
	return s == STATE_STOP
}
func (s State) isExit() bool {
	return s == STATE_EXIT
}
func (s State) isDone() bool {
	return s == STATE_DONE
}
func (s State) isFinish() bool {
	return s.isExit() || s.isDone()
}

const (
	WORKER_LOAD_JOB = iota
	WORKER_INIT_FILE
	WORKER_START_CMD
)

type Worker struct {
	index        int
	Name         string
	wId          string
	Sync         bool
	state        State
	Env          []string
	wf           *WorkFlow
	DependsOn    []*Worker
	ExecFunc     func(*Context) error
	FailureFunc  func(error, *Context) error
	FailureEnd   bool
	failureTimes int
	RepeatTimes  int
	c            *Cmd
	Cmd          string
	Input        string
	Output       string
	Errput       string
	Dir          string
}

func NewWorker(wf *WorkFlow) *Worker {
	index := wf.len() + 1
	w := &Worker{
		index: index,
		wId:   fmt.Sprintf("work-%d-%s", index, time.Now().String()),
		wf:    wf,
		state: STATE_PEND,
		Env:   make([]string, 0, 50),
	}

	return w
}

func (w *Worker) init(wf *WorkFlow) {
	w.index = wf.len() + 1
	w.wId = fmt.Sprintf("work-%d-%s", w.index, time.Now().String())
	w.wf = wf
	w.state = STATE_PEND
	w.Env = make([]string, 0, 50)
	w.SetExecFunc(w.RunCmd)
}

func (w *Worker) RunCmd(c *Context) error {
	w.c = newCmd(w.Cmd, nil)
	w.c.init(w.Input, w.Output, w.Errput, w.Dir, w.Env)
	return w.c.Run()
}

func (w *Worker) workFlowId() string {
	return w.wf.fId
}

func (w *Worker) SetExecFunc(fn func(*Context) error) {
	if w.ExecFunc == nil {
		w.ExecFunc = fn
	}
}

func (w *Worker) SetFailureFunc(fn func(error, *Context) error) {
	if w.FailureFunc == nil {
		w.FailureFunc = fn
	}
}

func (w *Worker) doExecFunc() error {
	if w.ExecFunc != nil {
		return w.ExecFunc(&w.wf.c)
	}
	return nil
}

func (w *Worker) doFailureFunc(err error) error {
	if w.FailureFunc != nil {
		return w.FailureFunc(err, &w.wf.c)
	}
	return nil
}

func (w *Worker) start(wg *sync.WaitGroup) error {
	fn := func() error {
		err := w.doExecFunc()
		if err != nil {
			w.state = STATE_EXIT
			w.doFailureFunc(err)
		} else {
			w.state = STATE_DONE
		}

		return err
	}
	return w.run(wg, fn)
}

func (w *Worker) run(wg *sync.WaitGroup, fn func() error) error {
	w.state = STATE_RUN
	w.wf.phase = w.index
	if !w.Sync {
		err := fn()
		w.wf.err <- Error{w.index, err}
		wg.Done()
		return err
	}

	go func(w *Worker, wg *sync.WaitGroup) {
		defer func() {
			if err := recover(); err != nil {
				w.wf.err <- Error{w.index, err.(error)}
			}
		}()
		err := fn()
		fmt.Println(w.Name, err)
		w.wf.err <- Error{w.index, err}
		wg.Done()
	}(w, wg)

	return nil
}

func (w Worker) IsFailureEnd() bool {
	if !w.state.isExit() || !w.FailureEnd {
		return false
	}
	return true
}

func (w *Worker) IsFailureRepeat() bool {
	if w.RepeatTimes == 0 || w.failureTimes >= w.RepeatTimes {
		return false
	}
	w.failureTimes++
	return true
}
