using System;

namespace Script.LogicFrame.Components
{
    public interface IComponent
    {
        public string Name { get; set; }
        public Int32 EntityID { get; set; }
    }
}