
using System;
using Script.Network.Connect;

namespace Script.Network.MsgProcessor
{
    public interface IMsgProcessor
    {
        public void OnConnection(TcpConnect connect);
        public void OnMessage(TcpConnect connect, UInt32 msgID, byte[] msgData);
        public void OnClose(TcpConnect connect);
    }

    
}