package main

import (
	"context"
	"net/http"
)

type Option func(*App)
type ShutdownCallback func(ctx context.Context)

func WithShutdownCallback(cbs ...ShutdownCallback) Option {
	return func(app *App) {
		app.cbs = cbs
	}
}

type App struct {
	servers *[]Server
	cbs     []ShutdownCallback
}

func (app *App) Start(addr string) {
	server := &serverMux{}
	server.Start(addr)
	servers := append(*app.servers, server)
	app.servers = &servers
}

func (app *App) Stop() {
	for _, server := range *app.servers {
		server.Stop()
	}
}

type Server interface {
	Start(addr string)
	Stop()
}

type serverMux struct {
	reject bool
	*http.ServeMux
}

func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}

func (s *serverMux) Start(addr string) {
	server := &http.Server{Addr: addr, Handler: nil}
	err := server.ListenAndServe()
	if err != nil {
		panic("启动服务失败，addr：" + addr)
	}
	s.ServeMux = http.NewServeMux()
}

func (s *serverMux) Stop() {
	//err := s.service.Close()
	//if err != nil {
	//	panic("关闭服务失败，addr：")
	//}
}

func main() {

}
