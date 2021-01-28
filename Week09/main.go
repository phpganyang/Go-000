package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("链接失败")
	}

	//处理链接
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Print(err)
			return
		}
		go handleConn(conn)
		fmt.Print("执行成功")
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	msgChan := make(chan string)
	//先读再写
	go handleWrite(conn, msgChan)
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			close(msgChan)
			return
		}
		msgChan <- msg
	}
}

func handleWrite(conn net.Conn, msgChan chan string) {
	//循环处理
	for {
		msg, ok := <-msgChan
		if !ok {
			fmt.Print("读取失败")
			return
		}
		fmt.Print(msg)
		_, err := io.WriteString(conn, msg)
		if err != nil {
			return
		}
	}
}
