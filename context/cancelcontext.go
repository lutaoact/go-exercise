package context

import (
	"fmt"
	"sync"
)

type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

type cancelCtx struct {
	Context
	done     chan struct{}
	mu       sync.Mutex
	children map[canceler]bool
	err      error
}

func (c *cancelCtx) Done() <-chan struct{} {
	return c.done
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (c *cancelCtx) String() string {
	return fmt.Sprintf("%v.WithCancel", c.Context)
}

//核心是关闭c.done
//同时会设置c.err = err, c.children = nil
//依次遍历c.children，每个child分别cancel
//如果设置了removeFromParent，则将c从其parent的children中删除
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return //already canceled
	}
	c.err = err
	close(c.done)
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()
	if removeFromParent {
		removeChild(c.Context, c) // 从此处可以看到 cancelCtx的Context项是一个类似于parent的概念
	}
}

type CancelFunc func()

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := newCancelCtx(parent)
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, Canceled) }
}

// WithCancel方法返回一个继承自parent的Context对象，同时返回的cancel方法可以用来关闭返回的Context当中的Done channel
// 其将新建立的节点挂载在最近的可以被cancel的父节点下（向下方向）
// 如果传入的parent是不可被cancel的节点，则直接只保留向上关系
func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{
		Context: parent,
		done:    make(chan struct{}),
	}
}

// 传递cancel
// 从当前传入的parent开始（包括该parent），向上查找最近的一个可以被cancel的parent
// 如果找到的parent已经被cancel，则将方才传入的child树给cancel掉
// 否则，将child节点直接连接为找到的parent的children中（Context字段不变，即向上的父亲指针不变，但是向下的孩子指针变直接了）
//
// 如果没有找到最近的可以被cancel的parent，即其上都不可被cancel，则启动一个goroutine等待传入的parent终止，则cancel传入的child树，或者等待传入的child终结。
func propagateCancel(parent Context, child canceler) {
	if parent.Done() == nil {
		return // parent is never canceled
	}
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]bool)
			}
			p.children[child] = true
		}
		p.mu.Unlock()
	} else {
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}

// 从传入的parent对象开始，依次往上找到一个最近的可以被cancel的对象，即cancelCtx或者timerCtx
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	for {
		switch c := parent.(type) {
		case *cancelCtx:
			return c, true
		case *timerCtx:
			return &c.cancelCtx, true
		case *valueCtx:
			parent = c.Context //循环查找
		default:
			return nil, false
		}
	}
}

//从parent开始往上找到最近的一个可以cancel的父对象
// 从父对象的children map中删除这个child
func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}
