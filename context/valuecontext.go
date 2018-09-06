package context

import "fmt"

type valueCtx struct {
	Context
	key, val interface{}
}

func WithValue(parent Context, key interface{}, val interface{}) Context {
	return &valueCtx{parent, key, val}
}

func (c *valueCtx) String() string {
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", c.Context, c.key, c.val)
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}
