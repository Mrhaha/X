using System.Collections.Generic;

namespace Script.LogicFrame.System
{
    public static class SystemManager
    {
        private static List<ISystem> _systems = new List<ISystem>();
        
        public static void InitSystem()
        {
            _systems.Add(new InputSystem());
        }

        public static void LogicUpdate()
        {
            foreach (var system in _systems)
            {
                system.LogicUpdate();
            }
        }
    }
}