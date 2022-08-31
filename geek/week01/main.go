package main

import (
	"bytes"
	"context"
	"encoding/json"
	"geek/week01/server"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	cache := &server.SyncCache{
		Cache: make(server.Cache, 1),
		DB:    make(server.DB, 1),
	}

	s1 := server.NewServer("biz", "localhost:8080")
	s1.Handle("/work", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		workHandler(writer, request, cache)
	}))
	s2 := server.NewServer("admin", "localhost:8081")
	app := server.NewApp(
		[]*server.Server{s1, s2},
		server.WithShutdownCallback(func(ctx context.Context) {
			StoreCacheToDBCallback(ctx, cache)
		}),
		server.WithShutdownTimeout(server.DefaultShutdownTimeout),
		server.WithWaitTimeout(server.DefaultWaitTimeout),
		server.WithCallbackTime(server.DefaultCbTimeout),
	)

	for i := 1; i <= 3; i++ {
		go func(i int) {
			time.Sleep(time.Second * 3 * time.Duration(i))
			work := "work" + strconv.Itoa(i)
			log.Printf("准备建立工作 %v", work)
			resp, err := http.Post("http://localhost:8080/work", "", bytes.NewBuffer([]byte(work)))
			if err != nil {
				log.Printf("建立工作 %v 失败，err: %v", work, err)
			} else {
				var body []byte
				body, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("建立工作 %v 失败，err: %v", work, err)
				} else {
					log.Printf("工作 %v 结束, body: %v", work, string(body))
				}
			}
		}(i)
	}

	app.StartAndServe()
}

// StoreCacheToDBCallback 关闭时刷新缓存到DB
func StoreCacheToDBCallback(ctx context.Context, cache *server.SyncCache) {
	done := make(chan int, 1)
	go func() {
		log.Printf("刷新缓存到DB\n")
		time.Sleep(time.Second * 2)
		cnt := cache.SyncAll()
		time.Sleep(time.Second * 2)
		done <- cnt
	}()
	select {
	case <-ctx.Done():
		log.Printf("关闭刷新缓存超时\n")
	case cnt := <-done:
		log.Printf("关闭刷新缓存处理结束，刷入缓存 %v\n", cnt)
	}
}

func workHandler(writer http.ResponseWriter, request *http.Request, cache *server.SyncCache) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		_, err = writer.Write([]byte("未接收到工作数据，请重新提交"))
		if err != nil {
			log.Printf("响应数据失败, err: %v", err)
		}
		return
	}
	log.Printf("获取到工作数据 %v，开始工作", string(body))
	time.Sleep(time.Second * 5)
	val, ok := cache.Get("key")
	if !ok {
		cache.Set("key", &server.Data{Name: string(body), Age: 18})
		val, _ = cache.Get("key")
	}
	data, _ := json.Marshal(val)
	log.Printf("工作结束 %v, 结算中", string(body))
	time.Sleep(time.Second * 5)
	_, err = writer.Write([]byte("结算工作 1000, val: " + string(data)))
	if err != nil {
		log.Printf("响应数据失败 %v, err: %v", body, err)
	}
}
