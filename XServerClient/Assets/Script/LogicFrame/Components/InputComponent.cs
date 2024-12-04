using Unity.VisualScripting;
using UnityEditor.U2D;

namespace Script.LogicFrame.Components
{
    public class InputComponent: IComponent
    {
        private float _inputX;
        private float _inputY;
        private bool _isJump;

        public float InputX
        {  
            get { return _inputX; }
            set { _inputX = value; }
        }
        
        public float InputY
        {  
            get { return _inputY; }
            set { _inputY = value; }
        }
        
        public bool IsJump
        {  
            get { return _isJump; }
            set { _isJump = value; }
        }
        
        public string GetName()
        {
            return "Input";
        }
    }
}