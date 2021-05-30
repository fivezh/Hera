package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 1500*time.Millisecond)
	time.Sleep(time.Millisecond * 600)
	select {
	// 在这种情况下，ctx传递到多个函数中，超时时会同时触发，顺序如何呢？
	case <-ctx.Done(): // 这里阻塞执行，等待ctx超时后，错误保存在ctx.Err()中
		fmt.Println("main", ctx.Err())
	}

	// 注意，这里多个goroutine中ctx.Done()顺序是无保证的，按goroutine调度情况
	/*
		[Running] go run "/Users/zhangxiaowu/workspace/github.com/Hera/golang/context/v2.go"
		main context deadline exceeded
		handle context deadline exceeded

		[Running] go run "/Users/zhangxiaowu/workspace/github.com/Hera/golang/context/v2.go"
		handle context deadline exceeded
		main context deadline exceeded
	*/

	// 这里需要阻塞等待子goroutine执行完成
	time.Sleep(time.Millisecond * 1000)
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done(): // 优先1s超时，ctx超时在父子协程中触发无顺序保证
		fmt.Println("handle", ctx.Err())
		time.Sleep(5 * time.Millisecond)
		fmt.Println("hello world in goroutine")
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
