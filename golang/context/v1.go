package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 1. ok
	go handle(ctx, 500*time.Millisecond)
	// 2. not sure
	// go handle(ctx, 1500*time.Millisecond) // 这里举例不当，超过Main的时间，一定概率main先执行后直接退出
	// fmt.Println("hello world") // 这里不阻塞执行
	select {
	case <-ctx.Done(): // 这里阻塞执行，等待ctx超时后，错误保存在ctx.Err()中
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done(): // 无机会执行到这里，因为先duration超时
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration): // 优先500ms超时，然后才是ctx超时
		fmt.Println("process request with", duration)
	}
}
