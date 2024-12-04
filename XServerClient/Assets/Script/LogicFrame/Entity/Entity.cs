using System;
using System.Collections.Generic;
using Google.Protobuf.Collections;
using Script.LogicFrame.Components;

namespace Script.LogicFrame.Entity
{
    public class Entity
    {
        private Int32 _uniqueID;
        private Dictionary<string, IComponent> _name2Component;

        public Entity()
        {}

        public Entity(Int32 uniqueID)
        {
            _uniqueID = uniqueID;
        }

        public void AddComponent(IComponent com)
        {
            _name2Component[com.GetName()] = com;
        }
        
        public IComponent GetComponentByName(string name)
        {
            if (_name2Component.TryGetValue(name, out var com))
            {
                return com;
            }
            return null;
        }
    }
}