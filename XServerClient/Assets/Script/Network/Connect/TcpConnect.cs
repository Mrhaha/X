using System;
using System.IO;
using System.Net.Sockets;
using System.Reflection;
using System.Threading.Tasks;
using Google.Protobuf;
using Script.Network.Crypto;
using Script.Network.MsgProcessor;
using Script.Network.util;
using UnityEngine;


namespace Script.Network.Connect
{
    public class TcpConnect
    {
        private IMsgProcessor _msgProcessor;
        private MsgPackager.MsgPackager _msgPackager;
        private TcpClient _tcpClient;
        private ICrypto _crypto;

        public void NewTcpConnect(string hostname, int port,ICrypto crypto)
        {
            TcpClient tcpClient = new TcpClient();
            tcpClient.Connect(hostname,port);

            _tcpClient = tcpClient;
            _msgPackager = new MsgPackager.MsgPackager();
            _msgPackager.Init(2,0,null);
            _msgProcessor = new TcpMsgProcessor();
            _crypto = crypto;
        }
        
        
        public async void ReadLoop()
        {
            while(true)
            {
                var ret =  await ReceiveMsg();
                if (ret.errorCode != 0)
                {
                    Debug.Log("readLoop Error: " + ret.errorCode);
                    return;
                }
                
                var msgHandler = NetManager.NetManager.GetMsgHandler(ret.msgID);
                var msgProto = NetManager.NetManager.GetMsgProtoTypeByMsgID(ret.msgID);
                var messageType = msgProto.Descriptor.ClrType; 
                if (null != msgHandler && messageType != null)
                {
                    var descriptor = messageType.GetProperty("Descriptor")?.GetValue(null);
                    var parser = descriptor?.GetType().GetProperty("Parser")?.GetValue(descriptor);
                    if (null != parser)
                    {
                        var parseMethod = parser.GetType().GetMethod("ParseFrom", new[] { typeof(CodedInputStream) });
                        if (null != parseMethod)
                        {
                            var codedInputStream = new CodedInputStream(new MemoryStream(ret.msgData));
                            var parsedMessage = parseMethod.Invoke(parser, new object[] { codedInputStream });
                            // 处理反射返回的消息实例
                            msgHandler((IMessage)parsedMessage);
                        }
                    }
                }
                _msgProcessor.OnMessage(this,ret.msgID,ret.msgData);
            }
        }
        
        private async Task<(UInt32 msgID,byte[] extData,byte[] msgData,int errorCode)> ReceiveMsg()
        {
            var ret = await _msgPackager.ReadMsg(_tcpClient.GetStream(), _crypto);
            return ret;
        }

        private void SendMsgByte(UInt32 msgID,byte[] extData,byte[] msgData)
        {
            _msgPackager.WriteMsg(_tcpClient.GetStream(), _crypto, msgID, extData, msgData);
        }

        public void SendMsg(IMessage dataMsg)
        {
            var msgID = ProtoUtil.ProtoMsg2MsgID(dataMsg);
            var extMsgBytes = Array.Empty<byte>();
            var dataMsgBytes = dataMsg?.ToByteArray() ?? Array.Empty<byte>();
            SendMsgByte(msgID,extMsgBytes,dataMsgBytes);
        }
    }
}