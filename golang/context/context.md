# context使用说明

## context.Done()在多个goroutine中顺序

多个协程中都有`ctx.Done()`时，是无法保证哪个优先执行的，按照goroutine调度情况执行。因此，不要假设哪个协程优先执行。


## AfterFunc作用

## WithDeadline() 实现

## parentCancelCtx() 作用

```golang
// parentCancelCtx returns the underlying *cancelCtx for parent.
// It does this by looking up parent.Value(&cancelCtxKey) to find
// the innermost enclosing *cancelCtx and then checking whether
// parent.Done() matches that *cancelCtx. (If not, the *cancelCtx
// has been wrapped in a custom implementation providing a
// different done channel, in which case we should not bypass it.)
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	p.mu.Lock()
	ok = p.done == done
	p.mu.Unlock()
	if !ok {
		return nil, false
	}
	return p, true
}
```

- 向上找到首个可取消的ctx