using System;
using System.IO;
using System.Linq;
using System.Net.Sockets;
using System.Numerics;
using System.Runtime.Serialization;
using System.Threading.Tasks;
using Google.Protobuf.Reflection;
using Script.Network.Crypto;
using Unity.VisualScripting;
using UnityEditor.PackageManager;
using UnityEngine;
using UnityEngine.Windows;

namespace Script.Network.MsgPackager
{
    public static class MsgPackagerDefine
    {
        public const Int32 MessageIDSize = 4;
        public const Int32 MessageLenSize = 2;
        public const Int32 MessageMaxLen = 1024000;
    }
    
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
    
    public class MsgPackager
    {
        private Int32 HeaderLen { get; set; }    //头部长度
        private Int32 DataLenSize { get; set; }      //数据长度占用
        private Int32 ExtLenSize { get; set; }       //额外数据长度占用
        private UInt32 ExtMaxLen { get; set; }    //额外数据最大长度
        public UInt32 MsgMaxLen { get; set; }    //消息体的最大长度
        public MessageDescriptor ExtType { get; set; } //额外消息类型

        public int Init(Int32 dataLenSize,Int32 extLenSize,MessageDescriptor extType)
        {
            HeaderLen = extLenSize + dataLenSize;           //头部只有两个字段(extLenSize,dataLenSize)
            DataLenSize = dataLenSize;
            ExtLenSize = extLenSize;
            
            switch (extLenSize)
            {
                case 0:
                    ExtMaxLen = 0;
                    break;
                case 1:
                    ExtMaxLen = byte.MaxValue;
                    break;
                case 2:
                    ExtMaxLen = ushort.MaxValue;
                    break;
                default:
                    return -1;
            }
            
            switch (dataLenSize)
            {
                case 1:
                    MsgMaxLen = byte.MaxValue;
                    break;
                case 2:
                    MsgMaxLen = ushort.MaxValue;
                    break;
                case 4:
                    MsgMaxLen = UInt32.MaxValue;
                    break;
                default:
                    return -1;
            }

            var msgLen = HeaderLen + MsgPackagerDefine.MessageIDSize + MsgMaxLen + ExtMaxLen;
            if (msgLen > MsgPackagerDefine.MessageMaxLen)
            {
                return -1;
            }

            return 0;
        }

        private int _encodeHeader(out MemoryStream buff,UInt32 el,UInt32 ml)
        {
            
            buff = new MemoryStream();
            switch (ExtLenSize)
            {
                case 1:
                    byte[] elByte = BitConverter.GetBytes((byte)el);
                    if (BitConverter.IsLittleEndian)
                    {
                        Array.Reverse(elByte);
                    }
                    buff.Write(elByte, 0, sizeof(byte));
                    break;
                case 2:
                    byte[] elShort = BitConverter.GetBytes((ushort)el);
                    if (BitConverter.IsLittleEndian)
                    {
                        Array.Reverse(elShort);
                    }
                    buff.Write(elShort, 0, sizeof(ushort));
                    break;
            }
            
            switch (DataLenSize)
            {
                case 1:
                    byte[] mlByte = BitConverter.GetBytes((byte)ml);
                    if (BitConverter.IsLittleEndian)
                    {
                        Array.Reverse(mlByte);
                    }
                    buff.Write(mlByte, 0, sizeof(byte));
                    break;
                case 2:
                    byte[] mlShort = BitConverter.GetBytes((byte)ml);
                    if (BitConverter.IsLittleEndian)
                    {
                        Array.Reverse(mlShort);
                    }
                    buff.Write(mlShort, 0, sizeof(ushort));
                    break;
                default:
                    return -1;
            }

            return 0;
        }

        private int  _decodeHeader(MemoryStream buff,out UInt32 el, out UInt32 ml)
        {
            el = 0;
            ml = 0;
            
            switch (ExtLenSize)
            {
                case 1:
                    el = (byte)buff.ReadByte();
                    break;
                case 2:
                    byte[] buffer = new byte[2];
                    buff.Read(buffer, 0, 2);
                    if (BitConverter.IsLittleEndian)
                    {
                        Array.Reverse(buffer);
                    }
                    el = BitConverter.ToUInt16(buffer);
                    break;
            }
            
            
            switch (DataLenSize)
            {
                case 1:
                    ml = (ushort)buff.ReadByte();
                    break;
                case 2:
                    byte[] buffer = new byte[2];
                    buff.Read(buffer, 0, 2);
                    if (BitConverter.IsLittleEndian)
                    {
                        Array.Reverse(buffer);
                    }
                    ml = BitConverter.ToUInt16(buffer);
                    break;
                default:
                    return -1;
            }
            return 0;
        }

        private int _encodeBody(out MemoryStream buff, ICrypto crypto, UInt32 id, byte[] extData, byte[] msgData)
        {
            var idBuffer = BitConverter.GetBytes(id);
            buff = new MemoryStream();
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(idBuffer);
            }
            buff.Write(idBuffer,0,sizeof(UInt32));

            if (extData.Length > 0)
            {
                buff.Seek(0, SeekOrigin.End);
                buff.Write(extData,0,extData.Length);
            }
            
            if (msgData.Length > 0)
            {
                buff.Seek(0, SeekOrigin.End);
                buff.Write(msgData,0,msgData.Length);
            }

            if (crypto != null)
            {
                crypto.Encrypt(buff.ToArray(),out var encryptBuff);
                buff.Seek(0, SeekOrigin.Begin);
                buff.Write(encryptBuff,0,encryptBuff.Length);
            }

            return 0;
        }

        private static (UInt32 msgID,byte[] extData,byte[] msgData,int errorCode) _decodeBody(MemoryStream buff, ICrypto crypto, UInt32 el)
        {
            var extData = new byte[el];
            var msgData = new byte[buff.Length-el-MsgPackagerDefine.MessageIDSize];
            byte[] srcBuff;
            if (crypto != null)
            {
                crypto.Decrypt(buff.ToArray(), out var decryptBuff);
                srcBuff = decryptBuff;
            }
            else
            {
                srcBuff = buff.ToArray();
            }

            var msgIDBuff = new byte[MsgPackagerDefine.MessageIDSize];
            Array.Copy(srcBuff, 0, msgIDBuff, 0, MsgPackagerDefine.MessageIDSize);
                if (BitConverter.IsLittleEndian)
                {
                    Array.Reverse(msgIDBuff);
                }
                var msgID = BitConverter.ToUInt32(msgIDBuff);
                
                Array.Copy(srcBuff,MsgPackagerDefine.MessageIDSize,extData,0,el);
                Array.Copy(srcBuff,MsgPackagerDefine.MessageIDSize+el,msgData,0,msgData.Length);
           
            return (msgID, extData, msgData, 0);
        }

        
        //都是针对网络字节流的处理
        private static async Task<int> _recvByte(NetworkStream netStream, byte[] recvBuf, Int32 length)
        {
            //从网络流中读取数据，要卡住for循环
            var totalByte = 0;
            while (totalByte < length)
            {
                try
                {
                    var cnt = await netStream.ReadAsync(recvBuf,0,length);
                    if (cnt == 0)
                    {
                        return totalByte;
                    }
                    totalByte += cnt;
                }
                catch (Exception e)
                {
                    Console.WriteLine(e);
                    return -1;
                }

            }
            return totalByte;
        }
        private static void _sendByte(NetworkStream netStream, byte[] sendBuf)
        {
            netStream.Write(sendBuf,0,sendBuf.Length);
        }
        
        //上层接口调用
        public async Task<(UInt32 msgID,byte[] extData,byte[] msgData,int errorCode)> ReadMsg(NetworkStream netStream,ICrypto crypto)
        {
            //read Header
            var headerBuf = new byte[HeaderLen];
            var headerCnt = await _recvByte(netStream, headerBuf, HeaderLen);
            if (headerCnt < HeaderLen)
            {
                Debug.LogError("ReadMsg _recvByte error");
                return (0,null,null,-1);
            }
            
            //decode Header
            var headerMemoryStream = new MemoryStream(headerBuf);
            var decodeHeaderRet = _decodeHeader(headerMemoryStream, out var el, out var ml);
            if (decodeHeaderRet != 0)
            {
                Debug.LogError("ReadMsg _decodeHeader error");
                return (0,null,null,-1);
            }
            
            Debug.Log("extDataLen: " + el + "msgDataLen: " + ml);

            if (el > ExtMaxLen)
            {
                Debug.LogError("ReadMsg ExtMaxLen error");
                return (0,null,null,-1);
            }

            if (ml > MsgMaxLen)
            {
                Debug.LogError("ReadMsg DataLenSize error");
                return (0,null,null,-1);
            }
            
            //read body
            var msgBodySize = (Int32)(MsgPackagerDefine.MessageIDSize + el + ml);
            var bodyBuf = new byte[msgBodySize];
            var bodyCnt = await _recvByte(netStream, bodyBuf, msgBodySize);
            if (bodyCnt < msgBodySize)
            {
                Debug.LogError("ReadMsg _recvByte error");
                return (0,null,null,-1);
            }
            
            var memoryStream = new MemoryStream(bodyBuf);
            //var aesCrypto = new AesCrypto();
            var retTuple = _decodeBody(memoryStream, null,el);
            return retTuple;
        }

        public Int32 WriteMsg(NetworkStream networkStream,ICrypto crypto,UInt32 msgID,byte[] extData,byte[] msgData)
        {
            var hRet = _encodeHeader(out var headerMemoryStream, (UInt32)extData.Length, (UInt32)msgData.Length);
            if (hRet != 0)
            {
                Debug.LogError("WriteMsg _encodeHeader: " + hRet);
                return hRet;
            }
            var bRet = _encodeBody(out var bodyMemoryStream, crypto, msgID, extData, msgData);
            if (bRet != 0)
            {
                Debug.LogError("WriteMsg _encodeBody: " + bRet);
                return bRet;
            }
            
            var mergeStream = new MemoryStream();
            headerMemoryStream.Seek(0, SeekOrigin.Begin);
            headerMemoryStream.CopyTo(mergeStream);
            bodyMemoryStream.Seek(0, SeekOrigin.Begin);
            bodyMemoryStream.CopyTo(mergeStream);
            
            _sendByte(networkStream,mergeStream.ToArray());
            return 0;
        }
    }
}