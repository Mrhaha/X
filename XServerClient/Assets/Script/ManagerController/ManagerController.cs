using System.Collections.Generic;
using System.Linq;
using System.Linq.Expressions;
using Script.FrameSync;
using Unity.Profiling;
using Unity.VisualScripting;

namespace Script.ManagerController
{

    public interface IManager
    {
        public void Init();
        public string GetStringName();

        public void Update(float dt);
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
            RegisterManager(new FrameSyncManager("FrameSync"));
            RegisterManager(new BattleInputManager("BattleInput"));
            
            foreach (var manager in _managers)
            {
                manager.Init();
            }
        }

        //********注册顺序决定 Update顺序*****//
        public static void Update(float dt)
        {
            foreach (var manager in _managers)
            {
                manager.Update(dt);
            }
        }
        
        
        private static void RegisterManager(IManager manager)
        {
            _managers.Add(manager);
        }

        public static IManager GetManagerByStringName(string name)
        {
            foreach (var manager in _managers)
            {
                if (manager.GetStringName() == name)
                {
                    return manager;
                }
            }

            return null;
        }
        
    }
}