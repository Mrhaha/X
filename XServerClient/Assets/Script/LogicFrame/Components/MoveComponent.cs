using System;

namespace Script.LogicFrame.Components
{
    public class MoveComponent
    {
        public float SpeedX;
        public float SpeedY;
        public float SpeedZ;

        public float PosX;
        public float PosY;
        public float PosZ;
        
        public float ScaleX;
        public float ScaleY;
        public float ScaleZ;

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