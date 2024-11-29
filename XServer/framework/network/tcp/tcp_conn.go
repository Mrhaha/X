package tcp

import (
	"errors"
	"io"
	"net"
	"sync/atomic"
	"time"

	"XServer/framework/network/crypto"
	"XServer/framework/network/msgpackager"
	"XServer/framework/network/msgprocessor"
	"XServer/framework/util"
	gxbytes "github.com/dubbogo/gost/bytes"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var (
	writeMsgChanTimeout   = time.Second * 3
	tcpConnDeadlineSecond = time.Second * 1
	tcpConnWaitSecond     = 3
)

var (
	// ErrClosed 已关闭
	ErrClosed = errors.New("ErrClosed")
	// ErrWriteChanFull 写管道满了
	ErrWriteChanFull = errors.New("ErrWriteChanFull")
	// ErrMsgIsNil 空消息
	ErrMsgIsNil = errors.New("ErrMsgIsNil")
)

type writemsg struct {
	ext   interface{}
	msgid uint32
	msg   interface{}
}

// Conn ...
// 连接到服务器的链接信息
type Conn struct {
	name     string
	conn     net.Conn
	isServer bool

	writeChan           chan msgprocessor.Msg
	closed              int32
	isCloseWhenBuffFull bool

	msgPackager  msgpackager.IMsgPackager
	msgProcessor msgprocessor.IMsgProcessor
	crypto       crypto.Crypto

	localAddr  net.Addr
	remoteAddr net.Addr

	reader  io.Reader
	writer  io.Writer
	rBuffer *gxbytes.Buffer
	wBuffer *gxbytes.Buffer
}

func newTCPConn(name string, conn net.Conn, isServer bool, tcpConnWriteChanLen int, bCloseBuffFull bool, msgPackager msgpackager.IMsgPackager, msgProcessor msgprocessor.IMsgProcessor, crypto crypto.Crypto) *Conn {
	tcpConn := new(Conn)

	tcpConn.name = name
	tcpConn.conn = conn
	tcpConn.isServer = isServer

	tcpConn.writeChan = make(chan msgprocessor.Msg, tcpConnWriteChanLen)
	tcpConn.isCloseWhenBuffFull = bCloseBuffFull

	tcpConn.msgPackager = msgPackager
	tcpConn.msgProcessor = msgProcessor
	tcpConn.crypto = crypto

	tcpConn.localAddr = conn.LocalAddr()
	tcpConn.remoteAddr = conn.RemoteAddr()

	tcpConn.reader = io.Reader(conn)
	tcpConn.writer = io.Writer(conn)
	tcpConn.rBuffer = gxbytes.NewBuffer(nil)
	tcpConn.wBuffer = gxbytes.NewBuffer(nil)

	return tcpConn
}

func (tcpConn *Conn) log(args ...interface{}) {
	if tcpConn.isServer {
		logrus.WithField("conn", tcpConn.Name()).Debug(args...)
	} else {
		logrus.WithField("conn", tcpConn.Name()).Info(args...)
	}
}

func (tcpConn *Conn) logError(err error, args ...interface{}) {
	logrus.WithField("conn", tcpConn.Name()).WithError(err).Error(args...)
}

func (tcpConn *Conn) isClosed() bool {
	return atomic.LoadInt32(&tcpConn.closed) != 0
}

func (tcpConn *Conn) close() bool {
	if atomic.CompareAndSwapInt32(&tcpConn.closed, 0, 1) {
		tcpConn.log("Conn close succeed")
		return true
	}
	return false
}

func (tcpConn *Conn) closeConn(waitSec int) {
	conn, ok := tcpConn.conn.(*net.TCPConn)
	if ok {
		now := time.Now()
		_ = conn.SetReadDeadline(now.Add(tcpConnDeadlineSecond))
		_ = conn.SetWriteDeadline(now.Add(tcpConnDeadlineSecond))
		_ = conn.SetLinger(waitSec)
		err := conn.Close()
		if err != nil {
			tcpConn.logError(err, "Conn closeConn error")
		} else {
			tcpConn.log("Conn closeConn succeed")
		}
	}
}

// Close ...
func (tcpConn *Conn) Close() {
	if tcpConn.isClosed() {
		return
	}
	tcpConn.log("Conn Close start")
	err := tcpConn.doWrite(nil)
	if err == ErrWriteChanFull {
		tcpConn.log("Conn Close start force close conn")
		tcpConn.closeConn(0)
	}
	if tcpConn.close() {
		tcpConn.log("Conn Close start succeed")
	}
}

func (tcpConn *Conn) postWriteMsg(msg msgprocessor.Msg) bool {
	select {
	case tcpConn.writeChan <- msg:
		return true
	case <-time.After(writeMsgChanTimeout):
		return false
	}
}

func (tcpConn *Conn) doWrite(msg msgprocessor.Msg) error {
	// 需要控制的话，使用select+timeout 方式
	if tcpConn.isCloseWhenBuffFull && len(tcpConn.writeChan) == cap(tcpConn.writeChan) {
		tcpConn.logError(ErrWriteChanFull, "Conn doWrite error")
		tcpConn.closeConn(0)
		if tcpConn.close() {
			close(tcpConn.writeChan)
		}
		return ErrWriteChanFull
	}

	ok := tcpConn.postWriteMsg(msg)
	if !ok {
		return ErrWriteChanFull
	}
	return nil
}

// Name ...
func (tcpConn *Conn) Name() string {
	return tcpConn.name
}

// WriteMsg ...
func (tcpConn *Conn) WriteMsg(ext proto.Message, msg proto.Message) error {
	if msg == nil {
		return ErrMsgIsNil
	}
	if tcpConn.isClosed() {
		return ErrClosed
	}
	return tcpConn.doWrite(&msgprocessor.ProtoMsg{
		ExtMsg: ext,
		Msg:    msg,
	})
}

// WriteBytes ...
func (tcpConn *Conn) WriteBytes(ext proto.Message, msgid uint32, bytes []byte) error {
	if tcpConn.isClosed() {
		return ErrClosed
	}
	return tcpConn.doWrite(&msgprocessor.BytesMsg{
		ExtMsg: ext,
		MsgID:  msgid,
		Data:   bytes,
	})
}

// LocalAddr ...
func (tcpConn *Conn) LocalAddr() net.Addr {
	return tcpConn.localAddr
}

// RemoteAddr ...
func (tcpConn *Conn) RemoteAddr() net.Addr {
	return tcpConn.remoteAddr
}

// WriteLoop ...
func (tcpConn *Conn) WriteLoop() {
	defer util.Recover()
	defer func() {
		tcpConn.closeConn(tcpConnWaitSecond)
	}()

	for msg := range tcpConn.writeChan {
		if msg == nil {
			tcpConn.log("Conn WriteLoop exit loop")
			break
		}

		ext, msgid, data, err1 := msg.Marshal()
		if err1 != nil {
			tcpConn.logError(err1, "Marshal message error")
			continue
		}

		// 加密消息
		extData, err2 := tcpConn.msgPackager.OnEncodeExt(ext)
		if err2 != nil {
			tcpConn.logError(err2, "Encode ext error, msgid: ", msgid)
			break
		}

		// 打包消息
		err3 := tcpConn.msgPackager.WriteMsg(tcpConn.writer, tcpConn.wBuffer, tcpConn.crypto, msgid, extData, data)
		if err3 != nil {
			if err3 != io.EOF {
				tcpConn.logError(err3, "Conn write message, msgid: ", msgid)
			}
			continue
		}
	}

	tcpConn.log("Conn exit write loop")
}

// ReadLoop ...
func (tcpConn *Conn) ReadLoop() {
	defer util.Recover()
	defer func() {
		tcpConn.msgProcessor.OnClose(tcpConn)
		tcpConn.Close()
	}()

	tcpConn.msgProcessor.OnConnect(tcpConn)

	for {
		msgid, extData, msgData, err1 := tcpConn.msgPackager.ReadMsg(tcpConn.reader, tcpConn.rBuffer, tcpConn.crypto)
		if err1 != nil {
			if err1 != io.EOF {
				logger.WithFields(logrus.Fields{
					"name": tcpConn.Name(),
				}).WithError(err1).Error("Conn read message error")
			}
			break
		}

		ext, err2 := tcpConn.msgPackager.OnDecodeExt(extData)
		if err2 != nil {
			logger.WithFields(logrus.Fields{
				"conn":  tcpConn.Name(),
				"msgid": msgid,
			}).WithError(err2).Error("Decode ext error")
			break
		}

		err3 := tcpConn.msgProcessor.OnMessage(tcpConn, ext, msgid, msgData)
		if err3 != nil {
			logger.WithFields(logrus.Fields{
				"conn":  tcpConn.Name(),
				"msgid": msgid,
			}).WithError(err3).Error("OnMessage error")

			//消息处理出错。这里不断开连接
			continue
		}
	}

	tcpConn.log("Conn exit read loop")
}
