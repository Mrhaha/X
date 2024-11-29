using System;
using Script.Network.Connect;
using UnityEngine;

namespace Script.Network.MsgProcessor
{
    public class TcpMsgProcessor : IMsgProcessor
    {
        public void OnConnection(TcpConnect connect)
        {
            
        }

        public void OnMessage(TcpConnect connect, UInt32 msgID, byte[] msgData)
        {
            Debug.Log("receiveMsg " + msgID);
            var resp =  XFramework.RspSyncFrame.Parser.ParseFrom(msgData);
            Debug.Log("parser end: "+resp.Frame+" : "+resp.X+" : "+resp.Y);
        }

        public void OnClose(TcpConnect connect)
        {
            
        }
    }
}