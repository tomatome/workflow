// workflow_test.go
package workflow

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type te struct {
	num int
}

func Test(t *testing.T) {
	w1 := &Worker{
		Name:     "pre",
		ExecFunc: pre,
		//Sync:     true,
	}
	w2 := &Worker{
		Name:     "pre2",
		ExecFunc: pre2,
		//Sync:        true,
		//RepeatTimes: 2,
		//FailureEnd: true,
	}
	w3 := &Worker{
		Name:     "post",
		ExecFunc: post,
		//Sync:     true,
	}
	w4 := &Worker{
		Name:     "post2",
		ExecFunc: post2,
		//Sync:     true,
	}
	w5 := &Worker{
		Name: "cmd",
		//FailureEnd: true,
		//ExecFunc:   cmd,
		DependsOn: []*Worker{w3},
		Cmd:       "cmd /C ls -t&ls -l",
		//Output:    "out.txt",
		//Sync:      true,
	}

	wf := newWorkFlow("job")
	wf.phase = 1
	wf.c.SetData(&te{3})
	wf.AddWorker(w1)
	wf.AddWorker(w2)
	wf.AddWorker(w3)
	wf.AddWorker(w4)
	wf.AddWorker(w5)
	wf.Start()

	time.Sleep(10 * time.Second)
}

func pre(c *Context) error {
	fmt.Println("pre:")
	//fmt.Println(c.Data())
	t := c.Data().(*te)
	t.num = 9
	return errors.New("failed")
}
func pre2(c *Context) error {
	fmt.Println("pre2:")
	//fmt.Println(c.Data())
	t := c.Data().(*te)
	t.num = 9
	return errors.New("failed")
}
func cmd(c *Context) error {
	fmt.Println("cmd:")
	//fmt.Println(c.Data())
	t := c.Data().(*te)
	t.num = 19
	return nil
}

func post(c *Context) error {
	fmt.Println("post:")
	//fmt.Println(c.Data())
	return nil
}

func post2(c *Context) error {
	fmt.Println("post2:")

	return nil
}
