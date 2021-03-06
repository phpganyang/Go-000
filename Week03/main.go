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

	group, _ := errgroup.WithContext(context.Background())
	sig := make(chan os.Signal, 1)
	stop := make(chan struct{})
	start := make(chan struct{})
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	//处理信号
	group.Go(func() error {
		sig := <-sig
		fmt.Println(sig)
		if sig == syscall.SIGINT {
			close(stop)
		}
		fmt.Println("信号监听结束.")
		return errors.New("signal error")
	})
	server := http.Server{Addr: "127.0.0.1:9999"}
	group.Go(func() error {
		go func() {
			<-stop
			sig <- syscall.SIGTERM
			fmt.Println("服务中断")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			err := server.Shutdown(ctx)
			close(start)
			fmt.Printf("服务关闭原因：%v\n", err)
		}()

		fmt.Println("服务开始")
		close(start)
		return server.ListenAndServe()
	})

	go func() {
		<-start
		fmt.Println("模拟请求数据") //首先执行
		time.Sleep(time.Second * 3)
		close(stop)

	}()

	if err := group.Wait(); err != nil {
		fmt.Printf("gorutine退出原因:%v\n", err)
	}
}
