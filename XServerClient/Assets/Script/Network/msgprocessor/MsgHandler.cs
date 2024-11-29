using System;
using System.Collections.Generic;
using Google.Protobuf;
using Script.Network.util;

namespace Script.Network.MsgProcessor
{
    public delegate void MsgHandlerDelegate(IMessage rsp);
    public class MsgHandler
    {
        private static Dictionary<UInt32, MsgHandlerDelegate> _msgHandlerDictionary;
        private static Dictionary<UInt32, IMessage> _msgID2ProtoMsg;

        public  MsgHandler()
        {
            _msgHandlerDictionary = new Dictionary<uint, MsgHandlerDelegate>();
            _msgID2ProtoMsg = new Dictionary<uint, IMessage>();
        }
        
        public void RegisterMsgHandler(IMessage protoMsg,MsgHandlerDelegate msgHandler)
        {
            var msgID = ProtoUtil.ProtoMsg2MsgID(protoMsg);
            _msgHandlerDictionary[msgID] = msgHandler;
            _msgID2ProtoMsg[msgID] = protoMsg;
        }

        public  IMessage GetProtoMsgTypeByMsgID(UInt32 msgID)
        {
            return _msgID2ProtoMsg[msgID];
        }

        public  MsgHandlerDelegate GetMsgHandler(UInt32 msgID)
        {
            return _msgHandlerDictionary.TryGetValue(msgID, out var msgHandler) ?  msgHandler:null;
        }
    }
}