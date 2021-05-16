package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 500*time.Millisecond)
	time.Sleep(time.Millisecond * 600)
	fmt.Println("hello world") // 这里不阻塞执行
	select {
	case <-ctx.Done(): // 这里阻塞执行，等待ctx超时后，错误保存在ctx.Err()中
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
