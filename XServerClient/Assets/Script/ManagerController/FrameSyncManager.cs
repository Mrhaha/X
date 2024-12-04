using System;
using System.Collections.Generic;
using Google.Protobuf;
using Script.ManagerController;
using XFramework;
using Script.Network.NetManager;
using UnityEngine;

namespace Script.ManagerController
{
    public class FrameSyncManager:IManager
    {
        private string _name;
        private Queue<RspSyncFrame> _logicServerFrame;
        private Int32 _serverFrame;         //服务器帧
        private Int32 _clientFrame;         //客户端帧


        public FrameSyncType SyncType
        {
            get;
            set;
        }


        public string GetStringName()
        {
            return _name;
        }

        public FrameSyncManager(string name)
        {
            _name = name;
        }
        
        //**********帧同步相关处理*********//
        public RspSyncFrame GetLogicFrame()
        {
            if (_logicServerFrame.Count > 0)
            {
                return _logicServerFrame.Dequeue();
            }
            return null;
        }
        
        
        public void Update(float dt)
        {
            switch (SyncType)
            {
                case FrameSyncType.Local:
                    break;
                case FrameSyncType.Server:
                    break;
            }
        }
        
        
        //********回调方法注册**********//
        public void Init()
        {
            NetManager.RegisterMsgHandler(new RspReadyBattle(),HandlerRspReadyBattle);
            NetManager.RegisterMsgHandler(new RspNotifyGameStart(),HandlerRspNotifyGameStart);
            NetManager.RegisterMsgHandler(new RspSyncFrame(),HandlerRspSyncFrame);
            _logicServerFrame = new Queue<RspSyncFrame>();
            SyncType = FrameSyncType.Local;
        }

        private void HandlerRspReadyBattle(IMessage msg)
        {
            var rsp = (RspReadyBattle)msg;
            Debug.Log("HandlerRspReadyBattle");
        }
        
        private void HandlerRspNotifyGameStart(IMessage msg)
        {
            var rsp = (RspNotifyGameStart)msg;
            FrameSyncDefine.GameState = GameState.Start;
            Debug.Log("HandlerRspNotifyGameStart");
        }
        

        private void HandlerRspSyncFrame(IMessage msg)
        {
            var rsp = (RspSyncFrame)msg;
            _logicServerFrame.Enqueue(rsp);
            Debug.Log("HandlerRspSyncFrame");
        }
    }
}