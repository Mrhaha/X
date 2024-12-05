using System.Collections.Generic;
using System.Reflection;
using Script.LogicFrame.Components;
using Script.LogicFrame.System;
using XFramework;

namespace Script.ManagerController
{
    public class LogicFrameManager : IManager
    {
        private string _name;
        private float _accumulateDeltaTime;
        private Dictionary<string, ISystem> _name2System;
        private RspSyncFrame _curFrame;

        public string GetStringName()
        {
            return _name;
        }

        public void Init()
        {
            _name2System.Add("InputSystem",new InputSystem());
        }

        public void Update(float dt)
        {
            _accumulateDeltaTime += dt;
            var frameSyncManager = (FrameSyncManager)ManagerController.GetManagerByStringName("FrameSync");
            while (_accumulateDeltaTime >= frameSyncManager.FrameDelta)
            {
                _curFrame = frameSyncManager.GetLogicFrame();
                foreach (var pair in _name2System)
                {
                    pair.Value.LogicUpdate(_curFrame);
                }       
                _accumulateDeltaTime -= frameSyncManager.FrameDelta;
            }
        }
    }
}