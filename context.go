// context.go
package workflow

import (
	"context"
)

type Context struct {
	context.Context
	context.CancelFunc
	data interface{}
}

func (c *Context) SetData(data interface{}) {
	c.data = data
}

func (c *Context) Data() interface{} {
	return c.data
}
