package main

import (
	"XServer/serverproto/frame"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	TCPMsgHeadSize = 4
)

type server struct {
	frame.UnimplementedGreeterServer
	serverFrame             int64
	playerNum               int32
	inputNotify             chan struct{}
	streamSyncServerHandler map[int64]*frame.Greeter_SyncInfoServer
}

func (s *server) SyncInfo(stream frame.Greeter_SyncInfoServer) error {
	req, err := stream.Recv()
	if err == io.EOF || err != nil {
		return err
	}
	return err
}

func handleRequest(conn net.Conn) {
	// 读取客户端发来的消息
	msgHeadBuf := make([]byte, TCPMsgHeadSize)
	_, err := io.ReadFull(conn, msgHeadBuf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	defer listener.Close()

	for {
		// 等待客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Received connection from", conn.RemoteAddr().String())
		// 处理连接
		go handleRequest(conn)
	}
}
