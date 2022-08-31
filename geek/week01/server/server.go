package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Option func(*App)
type ShutdownCallback func(ctx context.Context)

func WithShutdownCallback(cbs ...ShutdownCallback) Option {
	return func(app *App) {
		app.cbs = cbs
	}
}
func WithShutdownTimeout(d time.Duration) Option {
	return func(app *App) {
		app.shutdownTimeout = d
	}
}

func WithWaitTimeout(d time.Duration) Option {
	return func(app *App) {
		app.waitTimeout = d
	}
}

func WithCallbackTime(d time.Duration) Option {
	return func(app *App) {
		app.cbTimeout = d
	}
}

// App 基础应用结构
type App struct {
	// servers 所有启动的服务注册
	servers []*Server
	// shutdownTimeout 优雅退出整体超时时间，默认30s
	shutdownTimeout time.Duration
	// waitTimeout 拒绝请求的等待超时时间，默认10s。当前拒绝请求策略使用的是等待一段时间
	waitTimeout time.Duration
	// cbTimeout 回调任务执行超时时间，默认3s。
	cbTimeout time.Duration
	cbs       []ShutdownCallback
}

// NewApp 创建新的应用
func NewApp(servers []*Server, opts ...Option) *App {
	app := &App{
		servers:         servers,
		shutdownTimeout: DefaultShutdownTimeout,
		waitTimeout:     DefaultWaitTimeout,
		cbTimeout:       DefaultCbTimeout,
	}

	// 使用 opts 进一步处理自定义的应用操作
	for _, opt := range opts {
		opt(app)
	}

	log.Printf("创建新应用，服务数量 %v\n", len(app.servers))
	return app
}

// shutdown 应用优雅退出逻辑
func (app *App) shutdown() {
	log.Printf("应用退出中，准备关闭所有服务，已停止接收新请求, 等待%v\n", app.waitTimeout)
	// 关闭前拒绝所有新请求
	for _, server := range app.servers {
		server.reject()
	}
	// 等待已有请求处理完毕，使用固定时间策略
	time.Sleep(app.waitTimeout)

	// 关闭所有服务器
	log.Printf("开始关闭服务器")
	ctx := context.Background()
	// 使用异步等待组，确保全部服务关闭后在继续
	var wg sync.WaitGroup
	for _, server := range app.servers {
		wg.Add(1)
		go func(server *Server) {
			if err := server.stop(ctx); err != nil {
				log.Printf("服务 %v 停止异常，后续会进行强制停止, err: %v\n", server.name, err)
			}
			wg.Done()
		}(server)
	}
	wg.Wait()

	// 执行关闭回调
	for _, cb := range app.cbs {
		wg.Add(1)
		go func(cb ShutdownCallback) {
			cbContext, cancel := context.WithTimeout(ctx, app.cbTimeout)
			defer cancel()
			cb(cbContext)
			wg.Done()
		}(cb)
	}
	wg.Wait()
}

// shutdownNow 应用立即退出逻辑
func (app *App) shutdownNow() {
	os.Exit(1)
}

// StartAndServe 启动服务
func (app *App) StartAndServe() {
	// 异步启动多个服务
	for _, server := range app.servers {
		// 服务注意闭包取值
		go func(server *Server) {
			if err := server.Start(); err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器 %v 正常关闭退出\n", server.name)
				} else {
					log.Printf("服务器 %v 异常退出, err: %v\n", server.name, err)
				}
			} else {
				log.Printf("服务器 %v 正常启动\n", server.name)
			}
		}(server)
	}

	// 退出信号监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, Signals...)

	// 处理退出信号
	select {
	case s := <-quit:
		log.Printf("接收到退出信号 %v 应用将开始进行优雅退出\n", s)
		// 启动异步任务监听超时或者强制退出
		go func() {
			select {
			case <-quit:
				// 又收到退出信号，进行强制退出
				log.Printf("第二次收到退出信号，应用进行强制退出\n")
				app.shutdownNow()
			case <-time.After(app.shutdownTimeout):
				log.Printf("优雅退出超时，应用进行强制退出\n")
				// 超时，进行强制退出
				app.shutdownNow()
			}
		}()
		// 第一次推出信号，进行关闭应用优雅退出
		app.shutdown()
	}

}

type Server struct {
	srv  *http.Server
	name string
	mux  *serverMux
}

func NewServer(name string, addr string) *Server {
	mux := &serverMux{ServeMux: http.NewServeMux()}
	return &Server{
		name: name,
		mux:  mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("服务 %v 启动成功\n", s.name)
	return s.srv.ListenAndServe()
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	log.Printf("服务 %v 处理路由 %v\n", s.name, pattern)
	s.mux.Handle(pattern, handler)
}

func (s *Server) stop(context context.Context) error {
	log.Printf("服务器 %v 正在关闭中 ...\n", s.name)
	return s.srv.Shutdown(context)
}

func (s *Server) reject() {
	s.mux.reject = true
}

type serverMux struct {
	reject bool
	*http.ServeMux
}

func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭\n"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}
