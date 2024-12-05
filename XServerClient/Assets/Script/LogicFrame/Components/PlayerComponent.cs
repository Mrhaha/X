using System;
using Unity.VisualScripting;
using UnityEditor.U2D;

namespace Script.LogicFrame.Components
{
    public class PlayerComponent: IComponent
    {
        public int PlayerID;

        public string Name
        {
            get;
            set;
        }

        public Int32 EntityID
        {
            get;
            set;
        }
        
        
    }
}