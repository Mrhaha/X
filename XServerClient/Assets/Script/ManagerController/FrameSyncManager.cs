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
        private readonly float  _frameDelta = 0.033f; 
        private Int32 _clientFrame = -1;                 //客户端帧（处理服务器帧更新）
        private Int32 _clientSendFrame = -1;             //客户端已经发送帧号 (发送完成更新)
        private Int32 _clientReceiveFrame = -1;          //客户端已经收到帧号 (收到服务器帧更新)


        public float FrameDelta => _frameDelta;
        //考虑一共需要几个帧号来进行区分
        //是否能发送该帧，当前发送的客户端帧号的前一帧号服务器的确认是不是已经收到了        
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
                _clientFrame++;
                return _logicServerFrame.Dequeue();
            }
            return null;
        }
        
        private bool SendCheck()
        {
            return _clientSendFrame > _clientReceiveFrame;
        }
        
        
        public void Update(float dt)
        {
            switch (SyncType)
            {
                case FrameSyncType.Local:
                    break;
                case FrameSyncType.Server:
                    if (!SendCheck())
                    {
                        return;
                    }
                    var battleInputManager = (BattleInputManager)ManagerController.GetManagerByStringName("BattleInput");
                    var syncFrame = new SyncFrame
                        { PlayerID = FrameSyncDefine.ClientPlayerID, Frame = _clientSendFrame, Input = battleInputManager.PackInput()};
                    var req = new ReqSyncFrame{Frame = syncFrame};
                    NetManager.TcpConnect.SendMsg(req);
                    _clientSendFrame++;
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
            _clientReceiveFrame++;
            Debug.Log("HandlerRspSyncFrame");
        }
    }
}