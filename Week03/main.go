package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type myHandler struct{}

func main() {
	sig := make(chan os.Signal, 1) //go中信号通知机制可以通过channel发送实现
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	group, _ := errgroup.WithContext(context.Background())
	//启动一个http服务
	group.Go(func() error {
		err := StartHttpServer()
		if err != nil {
			return err
		}
		return nil
	})

	//监听信号
	group.Go(func() error {
		for s := range sig {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit....", s)
				return errors.New("haha,exit")
			case syscall.SIGUSR1:
				fmt.Println("sigusr1", s)
			case syscall.SIGUSR2:
				fmt.Println("sigusr2", s)
			default:
				fmt.Println("other sig")
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println("app start success")
	} else {
		fmt.Println("app fail" + err.Error())
	}
}
func StartHttpServer() error {
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/bye", bye)
	server := &http.Server{
		Addr:         ":9999",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	return server.ListenAndServe()
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to week03"))
}

func bye(w http.ResponseWriter, r *http.Request) {
	time.Sleep(4 * time.Second)
	w.Write([]byte("bye bye,homework"))
}
