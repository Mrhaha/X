package main

import (
	"XServer/framework/ioservice"
	"XServer/framework/network/connection"
	"XServer/framework/network/msgpackager"
	"XServer/framework/network/msgprocessor"
	"XServer/framework/network/tcp"
	"XServer/serverproto/frame"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"sync"
)

var (
	chanLen    = 10240
	contentLen = 65530
	content    = make([]byte, contentLen)
)

func serverConnect(conn connection.Connection) {

}

func SyncFrameReqHandler(conn connection.Connection, ext proto.Message, msg proto.Message) {
	logrus.WithFields(logrus.Fields{
		"info": "receiveMsg",
	}).Infof("receive Client Msg")

	resp := &frame.RspSyncFrame{Frame: 3, X: 2, Y: 1}
	err := conn.WriteMsg(nil, resp)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"errorMsg": err.Error(),
		}).Error("Handler SyncFrameReqHandler error")
	}
}

func serverBytes(conn connection.Connection, head proto.Message, msgID uint32, msgData []byte) {
	err := conn.WriteBytes(nil, msgID, msgData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"msgID":   msgID,
			"msgData": len(msgData),
		}).WithError(err).Error("serverBytes conn.WriteMsg")
	}
}

func main() {
	//newCrypto := func() crypto.Crypto {
	//	return crypto.NewAesCryptoUseDefaultKey()
	//}
	serverMsgHandlerGetter := msgprocessor.NewMsgHandlers()
	serverMsgHandlerGetter.AddHandler((*frame.ReqSyncFrame)(nil), SyncFrameReqHandler)
	ioService := ioservice.NewIOService("io-service", chanLen)
	msgPackage := msgpackager.NewMsgPackager(binary.BigEndian, 2, 0, nil)
	msgProcess := msgprocessor.NewMsgProcessor(ioService, serverConnect, serverConnect, serverBytes, serverMsgHandlerGetter)
	tcpServer := tcp.NewTCPServer("frameServer", "127.0.0.1:5000", 10, chanLen, false,
		msgPackage, msgProcess, nil)

	if tcpServer == nil {
		logrus.Fatal("port is listening")
		return
	}
	ioService.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			si := <-c
			switch si {
			case os.Interrupt, os.Kill:
				logrus.Info("tcpServer exiting")
				tcpServer.Close()
				logrus.Info("tcpServer exited")
				logrus.Info("serverIO exiting")
				ioService.Fini()
				logrus.Info("serverIO exited")
				wg.Done()
				logrus.Info("server exited")
				return
			}
		}
	}()
	wg.Wait()
}
