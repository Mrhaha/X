using Google.Protobuf;
using Script.ManagerController;
using XFramework;
using Script.Network.NetManager;
using UnityEngine;

namespace Script.FrameSync
{
    public class FrameSyncManager:IManager
    {
        //********回调方法注册**********//
        public void Init()
        {
            NetManager.RegisterMsgHandler(new RspSyncFrame(),HandlerRspSyncFrame);
        }

        private void HandlerRspSyncFrame(IMessage msg)
        {
            var rsp = (RspSyncFrame)msg;
            Debug.Log("1111parser end: "+rsp.Frame+" : "+rsp.X+" : "+rsp.Y);
        }
    }
}