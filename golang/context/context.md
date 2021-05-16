# context使用说明

## context.Done()在多个goroutine中顺序

多个协程中都有`ctx.Done()`时，是无法保证哪个优先执行的，按照goroutine调度情况执行。因此，不要假设哪个协程优先执行。
