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

func main() {

	g, _ := errgroup.WithContext(context.Background())
	sig := make(chan os.Signal, 1)
	stop := make(chan struct{})

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	//处理信号
	g.Go(func() error {
		sig := <-sig
		fmt.Println(sig)
		if sig == syscall.SIGINT {
			close(stop)
		}
		fmt.Println("信号监听结束.")
		return errors.New("sigal error")
	})
	server := http.Server{Addr: "127.0.0.1:8080"}
	g.Go(func() error {
		go func() {
			<-stop
			sig <- syscall.SIGTERM
			fmt.Println("服务中断")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			err := server.Shutdown(ctx)
			fmt.Printf("服务关闭原因：%v\n", err)
		}()
		fmt.Println("服务开始")
		return server.ListenAndServe()
	})

	go func() {
		fmt.Println("模拟请求数据")
		time.Sleep(time.Second * 3)
		close(stop)
	}()

	if err := g.Wait(); err != nil {
		fmt.Printf("gorutine退出原因:%v\n", err)
	}
}
