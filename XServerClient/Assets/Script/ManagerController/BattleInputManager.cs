using XFramework;

namespace Script.ManagerController
{
    public class BattleInputManager:IManager
    {
        private string _name;
        public string GetStringName()
        {
            return _name;
        }

        public BattleInputManager(string name)
        {
            _name = name;
        }

        public void Init()
        {
            
        }
        
        public void Update(float dt)
        {
            var frameSyncManager = (FrameSyncManager)ManagerController.GetManagerByStringName("FrameSync");
            switch (frameSyncManager.SyncType)
            {
                case FrameSyncType.Local:
                    break;
                case FrameSyncType.Server:
                    break;
            }
        }
        
    }
}