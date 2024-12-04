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
	chanLen       = 10240
	contentLen    = 65530
	MaxFrameIndex = 4096
	content       = make([]byte, contentLen)

	ClientFrame    = make(map[int64][]*frame.ReqSyncFrame)
	ClientFrameIdx = make(map[int64]int32)

	syncClientFrameIndexMu = sync.RWMutex{}
	syncClientFrameMu      = sync.RWMutex{}

	ServerFrame            = make([]*frame.RspSyncFrame, 0)
	ServerFrameIndex       = 0
	StartGameNeedClientNum = 1
	GameStart              = frame.GameState_GameState_None
	PlayerID2Conn          = make(map[int64]connection.Connection)
)

func serverConnect(conn connection.Connection) {

}

func ReadyFrameSync() {
	for {
		isReady := true
		syncClientFrameIndexMu.Lock()
		for _, v := range ClientFrameIdx {
			if v < int32(ServerFrameIndex+1) || v == 0 {
				isReady = false
				break
			}
		}
		syncClientFrameIndexMu.Unlock()

		if isReady {
			rsp := &frame.RspSyncFrame{}
			syncClientFrameMu.Lock()
			for _, v := range ClientFrame {
				rsp.ServerFrame = append(rsp.ServerFrame, v[ServerFrameIndex].GetFrame())
			}
			syncClientFrameMu.Unlock()
			ServerFrameIndex++

			for _, conn := range PlayerID2Conn {
				err := conn.WriteMsg(nil, rsp)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"errorMsg": err.Error(),
					}).Error("Handler SyncFrameReqHandler error")
				}
			}
		}
	}
}

func ReadyBattleReqHandler(conn connection.Connection, ext proto.Message, msg proto.Message) {
	req := msg.(*frame.ReqReadyBattle)
	logrus.WithFields(logrus.Fields{
		"PlayerID": req.GetPlayerID(),
	}).Infof("receive Client Msg")

	//初始化客户端帧序列
	ClientFrame[req.GetPlayerID()] = make([]*frame.ReqSyncFrame, 0)

	syncClientFrameIndexMu.Lock()
	ClientFrameIdx[req.GetPlayerID()] = 0
	syncClientFrameIndexMu.Unlock()

	PlayerID2Conn[req.GetPlayerID()] = conn
	if len(ClientFrame) == StartGameNeedClientNum {
		GameStart = frame.GameState_GameState_Start
		for _, v := range PlayerID2Conn {
			notify := &frame.RspNotifyGameStart{}
			v.WriteMsg(nil, notify)
		}
		go ReadyFrameSync()
	}

	resp := &frame.RspReadyBattle{PlayerID: req.GetPlayerID()}
	err := conn.WriteMsg(nil, resp)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"errorMsg": err.Error(),
		}).Error("Handler SyncFrameReqHandler error")
	}
}

func SyncFrameReqHandler(conn connection.Connection, ext proto.Message, msg proto.Message) {
	req := msg.(*frame.ReqSyncFrame)
	logrus.WithFields(logrus.Fields{
		"PlayerID": req.GetFrame().PlayerID,
		"Frame":    req.GetFrame().GetFrame(),
		"X":        req.GetFrame().GetX(),
		"Y":        req.GetFrame().GetY(),
	}).Infof("receive Client Msg")
	if _, exist := ClientFrame[req.GetFrame().GetPlayerID()]; !exist {
		return
	}
	syncClientFrameMu.Lock()
	ClientFrame[req.GetFrame().GetPlayerID()] = append(ClientFrame[req.GetFrame().GetPlayerID()], req)
	syncClientFrameMu.Unlock()

	syncClientFrameIndexMu.Lock()
	ClientFrameIdx[req.GetFrame().GetPlayerID()]++
	ClientFrameIdx[req.GetFrame().GetPlayerID()] %= int32(MaxFrameIndex)
	syncClientFrameIndexMu.Unlock()

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
	serverMsgHandlerGetter.AddHandler((*frame.ReqReadyBattle)(nil), ReadyBattleReqHandler)

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
