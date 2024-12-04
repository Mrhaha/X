using System;
using Script.FrameSync;
using Script.LogicFrame.Entity;
using Script.LogicFrame.System;
using Script.ManagerController;
using UnityEngine;
using Script.Network.NetManager;
using XFramework;


public static class FrameSyncDefine
{
    public static Int64 ClientPlayerID = 1;
    public const Int32 Speed = 4;
    public const float LogicFrameInterval = 0.033f;
    public const  Int32 MaxFrameIndex = 4096;
    public static GameState GameState = GameState.None;
    
    public static float LogicFrameTimer;
    public static Int32 FrameIndex;
}


public class Main : MonoBehaviour
{
    private int inputX;
    private int inputY;
    
    // Start is called before the first frame update
    void Start()
    {
        //Init
        SystemManager.InitSystem();
        ManagerController.Init();
        NetManager.Init("127.0.0.1",5000,null);
        
        //后续绑定个button发送
        var req = new ReqReadyBattle {PlayerID = FrameSyncDefine.ClientPlayerID};
        NetManager.TcpConnect.SendMsg(req);
    }

    void ProcessLogicFrame()
    {
        FrameSyncDefine.FrameIndex = ++FrameSyncDefine.FrameIndex%FrameSyncDefine.MaxFrameIndex;
        var syncFrame = new SyncFrame
            { PlayerID = FrameSyncDefine.ClientPlayerID, Frame = FrameSyncDefine.FrameIndex, X = inputX, Y = inputY };
        var req = new ReqSyncFrame{Frame = syncFrame};
        inputX = 0;
        inputY = 0;
        NetManager.TcpConnect.SendMsg(req);
    }

    // Update is called once per frame
    void Update()
    {
        if (FrameSyncDefine.GameState != GameState.Start)
        {
            return;
        }
        
        ManagerController.Update(Time.deltaTime);
    }
}
