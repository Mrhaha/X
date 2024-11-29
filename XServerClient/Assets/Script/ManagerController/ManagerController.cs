using System.Collections.Generic;
using System.Linq;
using System.Linq.Expressions;
using Script.FrameSync;
using Unity.Profiling;

namespace Script.ManagerController
{



    public interface IManager
    {
        public void Init();
    }



    public static class ManagerController
    {
        private static List<IManager> _managers;

        static ManagerController()
        {
            _managers = new List<IManager>();
        }
        public static void Init()
        {
            //********类型注册********//
            RegisterManager(new FrameSyncManager());
            
            foreach (var manager in _managers)
            {
                manager.Init();
            }
        }
        
        private static void RegisterManager(IManager manager)
        {
            _managers.Add(manager);
        }
        
    }
}