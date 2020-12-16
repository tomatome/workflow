// exec.go
package workflow

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ACT int

const (
	ACT_KILL = iota + 1
	ACT_STOP
	ACT_RESUME
)

type UserToken struct {
	Domain string
	User   string
	Passwd string
	//Id     ID    //windows:sid  linux:uid
	//token  TOKEN //windows:token  linux:gid
	Home string
}

type Cmd struct {
	Token  *UserToken
	Input  string
	output string
	errput string
	*exec.Cmd
	done   chan bool
	Action chan ACT
}

func newCmd(s string, token *UserToken) *Cmd {
	c := &Cmd{
		Token: token,
		done:  make(chan bool),
	}
	args := strings.Fields(s)
	//args := append([]string{}, "cmd", s)
	c.Cmd = exec.Command(args[0], args[1:]...)
	//c.SysProcAttr = SysProcAttr(token)

	return c
}

func (c *Cmd) init(inFile, outFile, errFile, dir string, env []string) *Cmd {
	c.Input = inFile
	c.output = outFile
	c.errput = errFile
	c.Dir = dir
	c.Env = append(c.Env, env...)

	return c
}

func (c *Cmd) start() (err error) {
	if c.Input != "" {
		fdin, err := os.Open(c.Input)
		if err != nil {
			return err
		}
		c.Stdin = fdin
	}
	if c.output != "" {
		fdout, err := os.OpenFile(c.output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		c.Stdout = fdout
	} else {
		c.Stdout = os.Stdout
	}
	if c.errput != "" {
		fderr, err := os.OpenFile(c.errput, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		c.Stderr = fderr
	} else {
		c.Stderr = os.Stdout
	}

	err = c.Start()
	if err != nil {
		fmt.Println(err)
		return err
	}

	c.wait()

	return nil
}

func (c *Cmd) wait() {
	go func(c *Cmd) {
		c.Wait()
		c.done <- true
	}(c)
}

func (c *Cmd) Run() error {
	if err := c.start(); err != nil {
		return err
	}

	for {
		select {
		case act := <-c.Action:
			switch act {
			case ACT_KILL:
				//kill
				//c.kill()
			case ACT_STOP:
				//stop
				//c.stop()
			case ACT_RESUME:
				//resume
				//c.resume()
			default:
				fmt.Println("Unknown action:", act)
			}

		case <-c.done:
			code := c.ProcessState.ExitCode()
			cputime := c.ProcessState.UserTime() + c.ProcessState.SystemTime()
			fmt.Println("exit:", code, "with cputime:", cputime, "on time:", time.Now().String())
			return nil
		default:
		}

		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

/*func (c *Command) kill() {
	if c.Process != nil {
		err := KillChilds(-c.exec.Process.Pid, syscall.SIGTERM)
		if err != nil {
			common.Error("kill process group is failed:", err)
			return
		}
		time.AfterFunc(10*time.Second, func() {
			if c.exec.ProcessState.Exited() {
				return
			}
			err := KillChilds(-c.exec.Process.Pid, syscall.SIGKILL)
			if err != nil {
				common.Error("kill process group is failed:", err)
			}
		})
	}
}

func (c *Command) stop() {
	if c.exec.Process != nil {
		err := StopChilds(-c.exec.Process.Pid)
		if err != nil {
			common.Errorf("kill process group<%d> is failed: %v", c.exec.Process.Pid, err)
			return
		}
	}
}
func (c *Command) resume() {
	if c.exec.Process != nil {
		err := ResumeChilds(-c.exec.Process.Pid)
		if err != nil {
			common.Errorf("resume process group<%d> is failed: %v", c.exec.Process.Pid, err)
			return
		}
	}
}

func (c *Command) send(data []byte) {
	_, err := c.stdin.Write(data)
	if err != nil {
		common.Error("send:", err)
	}
	//c.stdin.Write([]byte("\n"))
	//c.stdin.Close()
}
*/
