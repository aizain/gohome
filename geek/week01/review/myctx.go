package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
	"time"
)

//func main() {
//	ctx := context.Background()
//	//WithValueDemo(ctx)
//	//WithCancelDemo(ctx)
//	//WithDeadlineDemo(ctx)
//	//WithTimeoutDemo(ctx)
//	//AfterFuncDemo()
//	ErrorGroupDemo(ctx)
//}

// WithValueDemo 安全数据传递
func WithValueDemo(ctx context.Context) {
	vCtx := context.WithValue(ctx, "traceid", "uihvabnwkkljasd")
	vVal := vCtx.Value("traceid")
	fmt.Printf("val: %v\n", vVal)
	dl, ok := vCtx.Deadline()
	fmt.Printf("df: %v, ok %v\n", dl, ok)
	err := vCtx.Err()
	fmt.Printf("err: %v\n", err)
	ch := make(chan any, 1)
	go func(vCtx context.Context) {
		v1Val := vCtx.Value("traceid")
		fmt.Printf("v1: %v\n", v1Val)
		v1Ctx := context.WithValue(vCtx, "async", "v1")
		time.Sleep(time.Second * 1)
		fmt.Printf("v1Val: %v\n", v1Ctx.Value("async"))

		ch <- v1Ctx.Value("async")
	}(vCtx)
	go func(vCtx context.Context) {
		bCtx := context.Background()
		v2Val := bCtx.Value("traceid")
		fmt.Printf("v2: %v\n", v2Val)
		v2Ctx := context.WithValue(vCtx, "async", "v2")
		time.Sleep(time.Second * 1)
		fmt.Printf("v2Val: %v\n", v2Ctx.Value("async"))

		ch <- v2Ctx.Value("async")
	}(vCtx)

	select {
	case <-vCtx.Done():
		fmt.Printf("done: %v", vCtx.Err())
	case data := <-ch:
		fmt.Printf("get async %v\n", vCtx.Value("async"))
		fmt.Printf("get data %v\n", data)
	}
}

// WithCancelDemo 控制-取消
func WithCancelDemo(ctx context.Context) {
	cCtx, cancel := context.WithCancel(ctx)
	ch := make(chan any, 1)
	go func() {
		ch <- 1
		_, cCancel := context.WithCancel(cCtx)
		cCancel()
	}()
	go func() {
		ch <- 1
		time.Sleep(time.Second * 2)
		ch <- 2
		cancel()
	}()

	for {
		select {
		case data := <-ch:
			fmt.Printf("val %v \n", data)
		case <-cCtx.Done():
			fmt.Printf("done \n")
			os.Exit(1)
		}
	}
}

// WithDeadlineDemo 控制-固定时间过期
func WithDeadlineDemo(ctx context.Context) {
	dlCtx, cancel := context.WithDeadline(ctx, time.UnixMilli(time.Now().UnixMilli()+(time.Second*5).Milliseconds()))
	defer cancel()
	t, ok := dlCtx.Deadline()
	fmt.Printf("dl: %v:%v\n", ok, t)
	time.Sleep(time.Second * 5)
	select {
	case <-dlCtx.Done():
		fmt.Printf("done %v\n", time.Now())
	}
}

// WithTimeoutDemo 控制-超时过期
func WithTimeoutDemo(ctx context.Context) {
	tCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	ch := make(chan any, 1)

	g := sync.WaitGroup{}
	g.Add(2)

	go func() {
		defer cancel()
		time.Sleep(time.Second * 1)
		ch <- 1
		time.Sleep(time.Second * 1)
		ch <- 1
		time.Sleep(time.Second * 1)
		ch <- 1
		g.Done()
	}()

	go func(tCtx context.Context) {
	a:
		for {
			select {
			case <-tCtx.Done():
				fmt.Printf("done\n")
				break a
			case val := <-ch:
				fmt.Printf("val %v\n", val)
			}
		}
		g.Done()
	}(tCtx)
	g.Wait()
}

// AfterFuncDemo 定时执行
func AfterFuncDemo() {
	g := sync.WaitGroup{}
	g.Add(2)
	fmt.Printf("begin %v\n", time.Now())
	time.AfterFunc(time.Second, func() {
		fmt.Printf("run timmer %v\n", time.Now())
		g.Done()
	})
	time.AfterFunc(time.Second*10, func() {
		fmt.Printf("run timmer %v\n", time.Now())
		g.Done()
	})
	g.Wait()
}

// ErrorGroupDemo 异常组
func ErrorGroupDemo(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-ctx.Done()
		fmt.Printf("done %v\n", ctx.Err())
		return ctx.Err()
	})

	g.Go(func() error {
		fmt.Printf("run\n")
		time.Sleep(time.Second)
		return errors.New("err")
	})

	g.Wait()
}

type MyContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
	A(key any) any
}

type myValueCtx struct {
	// 嵌入了该接口的方法
	MyContext
	key, val any
}

func WithMyValue(ctx MyContext, key string, val string) MyContext {
	return myValueCtx{
		ctx,
		key,
		val,
	}
}

func testMyContext() {
	var ctx *myEmptyCtx
	myCtx := WithMyValue(ctx, "", "")
	fmt.Printf("a: %v", myCtx.Value("key"))
}

type myEmptyCtx int

func (*myEmptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*myEmptyCtx) Done() <-chan struct{} {
	return nil
}

func (*myEmptyCtx) Err() error {
	return nil
}

func (*myEmptyCtx) Value(key any) any {
	return nil
}

func (*myEmptyCtx) A(key any) any {
	return nil
}
