package http

import (
	"Week04/api"
	"Week04/internal/service"
	"context"
	"fmt"
	"net/http"
	"sync"
)

type HttpServer struct {
	server *http.Server
	ctx    context.Context
	cancel func()
	once   *sync.Once
}

func New(svr *service.Service) (hsvr *HttpServer, cf func(), err error) {
	ctx, cancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc(api.HelloUrl, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			resp, err := svr.SayHello(ctx, &api.Demo{Id: 111})
			if err == nil {
				w.Write([]byte(resp.Content))
				w.WriteHeader(http.StatusOK)
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})
	hsvr = &HttpServer{
		server: &http.Server{
			Addr:    ":18000",
			Handler: mux,
		},
		ctx:    ctx,
		cancel: cancel,
		once:   &sync.Once{},
	}
	cf = hsvr.Close
	go hsvr.Start()
	return
}

func (hs *HttpServer) Start() {
	err := hs.server.ListenAndServe()
	if err != nil {
		fmt.Println("error:ListenAndServe():", err)
	}
}

func (hs *HttpServer) Close() {
	hs.once.Do(func() {
		hs.cancel()
		hs.server.Close()
	})
}
