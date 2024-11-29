package msgpackager

import (
	"XServer/framework/network/crypto"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"

	gxbytes "github.com/dubbogo/gost/bytes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// MessageIDSize = 4个字节长度, 不能改
	MessageIDSize = 4
	// MessageLenSize 消息头中表示消息长度的字节的大小
	MessageLenSize = 2
	// MessageMaxLen 消息最大长度
	MessageMaxLen = 1024000
)

var (
	// BigEndian ...
	BigEndian = binary.ByteOrder(binary.BigEndian)
	// LittleEndian ...
	LittleEndian = binary.ByteOrder(binary.LittleEndian)
)

//
// msg struct/msg packet
// ----------------------------------------
// | extlen | msglen | id | ext | msg |
// ----------------------------------------
// |      head       |       body     |
// |  none encrypted |    encrypted   |
// ----------------------------------------
// head 是 可根据不同的protocol来设定
// headerLen = extLenSize + dataLenSize
// msgHead = 包含两个部分extLen+msgLen
// msgBody = id + ext + msg

// IMsgPackager 管理协议的组织
type IMsgPackager interface {
	ReadMsg(reader io.Reader, buf *gxbytes.Buffer, crypto crypto.Crypto) (id uint32, extdata []byte, msgdata []byte, error error)
	WriteMsg(writer io.Writer, buf *gxbytes.Buffer, crypto crypto.Crypto, id uint32, extdata []byte, msgdata []byte) error

	// OnDecodeExt on read gorutine decode扩展数据.
	OnDecodeExt(extData []byte) (ext proto.Message, err error)
	// OnEncodeExt on write gorutine encode扩展数据.
	OnEncodeExt(ext proto.Message) (extData []byte, err error)
}

type msgPackager struct {
	byteOrder   binary.ByteOrder
	headerLen   int
	dataLenSize int
	extLenSize  int
	extMaxLen   uint32 // 扩展数据最大长度
	msgMaxLen   uint32 // 数据最大长度
	extType     protoreflect.MessageType
}

// NewMsgPackager Create a {| extLen | dataLen | id | ext | msg |} msg.
// The extLenSize 是 extlen 的字节数. extLenSize must is 0、1、2、4
// The dataLen 是 msglen 的字节数. dataLen must is 1、2、4
func NewMsgPackager(byteOrder binary.ByteOrder, dataLenSize int, extLenSize int, ext proto.Message) IMsgPackager {
	packager := &msgPackager{}
	err := packager.init(byteOrder, dataLenSize, extLenSize, ext)
	if err != nil {
		panic("init NewMsgPackager message ext type err, required pointer")
	}
	return packager
}

func (p *msgPackager) init(byteOrder binary.ByteOrder, dataLenSize int, extLenSize int, ext proto.Message) error {
	p.byteOrder = byteOrder
	p.headerLen = extLenSize + dataLenSize
	p.dataLenSize = dataLenSize
	p.extLenSize = extLenSize

	if extLenSize == 0 {
		p.extMaxLen = 0
	} else if extLenSize == 1 {
		p.extMaxLen = math.MaxUint8
	} else if extLenSize == 2 {
		p.extMaxLen = math.MaxUint16
	} else if extLenSize == 4 {
		p.extMaxLen = math.MaxUint32
	} else {
		return errors.New("unsupported packet ext len size")
	}
	if p.extMaxLen > MessageMaxLen {
		p.extMaxLen = MessageMaxLen
	}

	if dataLenSize == 1 {
		p.msgMaxLen = math.MaxUint8
	} else if dataLenSize == 2 {
		p.msgMaxLen = math.MaxUint16
	} else if dataLenSize == 4 {
		p.msgMaxLen = math.MaxUint32
	} else {
		return errors.New("unsupported packet msg len size")
	}

	if p.msgMaxLen > MessageMaxLen {
		p.msgMaxLen = MessageMaxLen
	}

	if ext != nil {
		extType := reflect.TypeOf(ext)
		if extType.Kind() != reflect.Ptr || extType.Elem() == nil {
			panic("init msgPackager message ext type err, required pointer")
		}
		p.extType = ext.ProtoReflect().Type()
	}

	return nil
}

func (p *msgPackager) encodeHeader(buf *gxbytes.Buffer, el uint32, ml uint32) error {
	switch p.extLenSize {
	case 1:
		err := binary.Write(buf, p.byteOrder, byte(el))
		if err != nil {
			return err
		}
	case 2:
		err := binary.Write(buf, p.byteOrder, uint16(el))
		if err != nil {
			return err
		}
	case 4:
		err := binary.Write(buf, p.byteOrder, uint32(el))
		if err != nil {
			return err
		}
	}

	switch p.dataLenSize {
	case 1:
		err := binary.Write(buf, p.byteOrder, byte(ml))
		if err != nil {
			return err
		}
	case 2:
		err := binary.Write(buf, p.byteOrder, uint16(ml))
		if err != nil {
			return err
		}
	case 4:
		err := binary.Write(buf, p.byteOrder, uint32(ml))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *msgPackager) decodeHeader(buf *gxbytes.Buffer) (el uint32, ml uint32, err error) {
	switch p.extLenSize {
	case 1:
		var n byte
		err = binary.Read(buf, p.byteOrder, &n)
		if err != nil {
			return
		}
		el = uint32(n)
	case 2:
		var n uint16
		err = binary.Read(buf, p.byteOrder, &n)
		if err != nil {
			return
		}
		el = uint32(n)
	case 4:
		err = binary.Read(buf, p.byteOrder, &el)
		if err != nil {
			return
		}
	}

	switch p.dataLenSize {
	case 1:
		var n byte
		err = binary.Read(buf, p.byteOrder, &n)
		if err != nil {
			return
		}
		ml = uint32(n)
	case 2:
		var n uint16
		err = binary.Read(buf, p.byteOrder, &n)
		if err != nil {
			return
		}
		ml = uint32(n)
	case 4:
		err = binary.Read(buf, p.byteOrder, &ml)
		if err != nil {
			return
		}
	}

	return
}

func (p *msgPackager) encodeBody(buf *gxbytes.Buffer, crypto crypto.Crypto, id uint32, ext, msg []byte) error {
	// write msgid
	err := binary.Write(buf, p.byteOrder, uint32(id))
	if err != nil {
		return err
	}
	// write ext
	if len(ext) > 0 {
		err = binary.Write(buf, p.byteOrder, ext)
		if err != nil {
			return err
		}
	}
	// write msg
	if len(msg) > 0 {
		err = binary.Write(buf, p.byteOrder, msg)
		if err != nil {
			return err
		}
	}
	if crypto != nil {
		crypto.Encrypt(buf.Bytes()[p.headerLen:], buf.Bytes()[p.headerLen:])
	}
	return nil
}

func (p *msgPackager) decodeBody(buf *gxbytes.Buffer, crypto crypto.Crypto, extLen uint32) (uint32, []byte, []byte, error) {
	if crypto != nil {
		err := crypto.Decrypt(buf.Bytes(), buf.Bytes())
		if err != nil {
			return 0, nil, nil, fmt.Errorf("decode decrypt body error: %v", err)
		}
	}
	// read msgid
	var id uint32
	err := binary.Read(buf, p.byteOrder, &id)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("decode read msgid error: %v", err)
	}
	// write ext & msg
	body := buf.Next(buf.Len())
	if p.extMaxLen != 0 && extLen > 0 {
		return id, body[:extLen], body[extLen:], nil
	}
	return id, nil, body, nil
}

func (p *msgPackager) read(r io.Reader, buf *gxbytes.Buffer, length int) (int, error) {
	n, err := io.ReadFull(r, buf.WriteNextBegin(length))
	if err != nil {
		if !(err == io.EOF && length == n) {
			return 0, err
		}
	}
	return buf.WriteNextEnd(n)
}

func (p *msgPackager) write(w io.Writer, buf *gxbytes.Buffer) (int64, error) {
	return buf.WriteTo(w)
}

// OnDecodeExt decode扩展数据
func (p *msgPackager) OnDecodeExt(extData []byte) (ext proto.Message, err error) {
	if extData != nil && p.extType != nil {
		ext = p.extType.New().Interface()
		err = proto.Unmarshal(extData, ext)
	}
	return ext, err
}

// OnEncodeExt encode扩展数据
func (p *msgPackager) OnEncodeExt(ext proto.Message) (extData []byte, err error) {
	if ext == nil {
		return nil, nil
	}
	return proto.Marshal(ext)
}

// ReadMsg ...
func (p *msgPackager) ReadMsg(reader io.Reader, buf *gxbytes.Buffer, crypto crypto.Crypto) (uint32, []byte, []byte, error) {
	// read header
	if _, err := p.read(reader, buf, p.headerLen); err != nil {
		return 0, nil, nil, err
	}
	// decode header
	extLen, msgLen, err := p.decodeHeader(buf)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("decode header error: %v", err)
	}
	if extLen > p.extMaxLen {
		return 0, nil, nil, errors.New("read ext too max")
	}
	if msgLen > p.msgMaxLen {
		return 0, nil, nil, errors.New("read msg too max")
	}
	// read body
	if _, err := p.read(reader, buf, int(MessageIDSize+extLen+msgLen)); err != nil {
		return 0, nil, nil, err
	}
	// decode body
	msgID, ext, msg, err := p.decodeBody(buf, crypto, extLen)
	if err != nil {
		return 0, nil, nil, err
	}
	return msgID, ext, msg, nil
}

// WriteMsg ...
func (p *msgPackager) WriteMsg(w io.Writer, buf *gxbytes.Buffer, crypto crypto.Crypto, id uint32, extdata []byte, msgdata []byte) error {
	msgLen := uint32(len(msgdata))
	if msgLen > p.msgMaxLen {
		return fmt.Errorf("write msgdata too max msgid: %d, len: %d", id, msgLen)
	}
	var extLen uint32
	if p.extMaxLen != 0 && extdata != nil {
		extLen = uint32(len(extdata))
	}

	buf.Reset()
	// encode head
	err := p.encodeHeader(buf, extLen, msgLen)
	if err != nil {
		return err
	}
	// encode body
	err = p.encodeBody(buf, crypto, id, extdata, msgdata)
	if err != nil {
		return err
	}
	// write to io
	if _, err = p.write(w, buf); err != nil {
		return err
	}
	return nil
}
