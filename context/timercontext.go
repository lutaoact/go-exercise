package context

import (
	"fmt"
	"time"
)

type timerCtx struct {
	cancelCtx //此处的封装为了继承来自于cancelCtx的方法，cancelCtx.Context才是父亲节点的指针
	timer     *time.Timer
	deadline  time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return fmt.Sprintf("%v.WithDeadline(%s [%s])", c.cancelCtx.Context, c.deadline, c.deadline.Sub(time.Now()))
}

// 与cencelCtx有所不同，其除了处理cancelCtx.cancel，还回对c.timer进行Stop()，并将c.timer=nil
func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		// Remove this timerCtx from its parent cancelCtx's children.
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	// 如果parent的deadline比新传入的deadline已经要早，则直接WithCancel，因为新传入的deadline没有效，父亲的deadline会先到期。
	if cur, ok := parent.Deadline(); ok && cur.Before(deadline) {
		return WithCancel(parent)
	}
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  deadline,
	}
	propagateCancel(parent, c)

	// 检查如果已经过期，则cancel新的子树
	d := deadline.Sub(time.Now())
	if d <= 0 {
		c.cancel(true, DeadlineExceeded)
		return c, func() { c.cancel(true, Canceled) }
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		// 还没有被cancel的话，就设置deadline之后cancel的计时器
		c.timer = time.AfterFunc(d, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, Canceled) }
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
