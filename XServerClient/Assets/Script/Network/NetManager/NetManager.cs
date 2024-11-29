using System;
using Google.Protobuf;
using Script.Network.Connect;
using Script.Network.Crypto;
using Script.Network.MsgProcessor;
using Script.Network.util;

namespace Script.Network.NetManager
{
    public static class NetManager
    {
        private static TcpConnect _tcpConnect;
        private static MsgHandler _msgHandler;

        static NetManager()
        {
            _msgHandler = new MsgHandler();
            _tcpConnect = new TcpConnect();
        }

        public static TcpConnect TcpConnect => _tcpConnect;
        
        public static void Init(string hostname, int port, ICrypto crypto)
        {
            _tcpConnect.NewTcpConnect(hostname,port,crypto);
            _tcpConnect.ReadLoop();
        }

        public static MsgHandlerDelegate GetMsgHandler(UInt32 msgID)
        {
            return _msgHandler.GetMsgHandler(msgID);
        }

        public static IMessage GetMsgProtoTypeByMsgID(UInt32 msgID)
        {
            return _msgHandler.GetProtoMsgTypeByMsgID(msgID);
        }
        
        public static void RegisterMsgHandler(IMessage protoMsg,MsgHandlerDelegate msgHandler)
        {
            _msgHandler.RegisterMsgHandler(protoMsg,msgHandler);
        }
    }
}