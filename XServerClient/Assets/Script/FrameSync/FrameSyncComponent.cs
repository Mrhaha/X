using System;
using UnityEngine;
using XFramework;

namespace Script.FrameSync
{
    public class FrameSyncComponent : MonoBehaviour
    {

        // public Int32 _PlayerID;
        // private float _accumulateDeltaTime = 0f;
        // private Vector3 _currentPosition;
        // private Vector3 _updatePosition;
        //
        // public void Start()
        // {
        // }
        //
        // public void ProcessRenderFrame(RspSyncFrame nexLogicFrame, float deltaTime)
        // {
        //     if (nexLogicFrame != null)
        //     {
        //         _accumulateDeltaTime = 0;
        //         SyncFrame syncFrame = null;
        //         var betweenPosition = FrameSyncDefine.Speed * FrameSyncDefine.LogicFrameInterval * Vector3.one;
        //         foreach (var frame in nexLogicFrame.ServerFrame)
        //         {
        //             if (frame.PlayerID != _PlayerID)
        //             {
        //                 continue;
        //             }
        //
        //             syncFrame = frame;
        //         }
        //
        //         if (syncFrame != null)
        //         {
        //             betweenPosition.x *= syncFrame.X;
        //             betweenPosition.y *= syncFrame.Y;
        //             betweenPosition.z *= 0;
        //             _currentPosition = gameObject.transform.position;
        //             _updatePosition = _currentPosition + betweenPosition;
        //         }
        //     }
        //
        //     _accumulateDeltaTime += deltaTime;
        //     var t = _accumulateDeltaTime / FrameSyncDefine.LogicFrameInterval;
        //     gameObject.transform.position = Vector3.Lerp(_currentPosition, _updatePosition, t);
        // }
    }
}