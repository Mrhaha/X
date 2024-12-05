﻿using System;
using Unity.VisualScripting;
using UnityEditor.U2D;

namespace Script.LogicFrame.Components
{
    public class InputComponent: IComponent
    {
        public float InputX;
        public float InputY;
        public bool IsJump;

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