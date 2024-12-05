using UnityEngine;
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


        public PlayerInput PackInput()
        {
            var input = new PlayerInput();
            input.X = Input.GetAxis("Horizontal");
            input.Y = Input.GetAxis("Vertical");
            input.IsJump = Input.GetKeyDown(KeyCode.Space);
            return input;
        }
        
        public void Init()
        {
            
        }
        
        public void Update(float dt)
        {
            
        }
        
    }
}