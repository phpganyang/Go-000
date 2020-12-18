package main

import (
	"Week04/internal/di"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-c
	fmt.Println("receive signal:", s)
	closeFunc()
	fmt.Println("exit ...")
}
