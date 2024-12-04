using Script.LogicFrame.Components;
using Script.LogicFrame.Entity;
using UnityEngine;

namespace Script.LogicFrame.System
{
    public class InputSystem:ISystem
    {
        public void LogicUpdate()
        {
            var moveHorizontal = Input.GetAxis("Horizontal");
            var moveVertical = Input.GetAxis("Vertical");
            var  inputJump = Input.GetKeyDown(KeyCode.Space);

            foreach (var entityPair in EntityManager.UniqueID2Entity)
            {
                var inputComp = (InputComponent)entityPair.Value.GetComponentByName("Input");
                if (moveHorizontal != 0)
                {
                    inputComp.InputX  = moveHorizontal == 0 ? 0 : (moveHorizontal > 0 ? 1 : -1);
                }
                if (moveVertical != 0)
                {
                    inputComp.InputY = moveVertical == 0 ? 0 : (moveVertical > 0 ? 1 : -1);
                }

                inputComp.IsJump = inputJump;
            }
        }
    }
}