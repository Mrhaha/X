using Script.LogicFrame.Components;
using Script.LogicFrame.Entity;
using UnityEngine;
using XFramework;

namespace Script.LogicFrame.System
{
    public class InputSystem:ISystem
    {
        public void LogicUpdate(RspSyncFrame curFrame)
        {
            var iComponents = EntityManager.GetComponentsByStringName("CInput");
            foreach (var iComponent in iComponents)
            {
                var inputComponent = (InputComponent)iComponent;
                var playerComponent = (PlayerComponent)EntityManager.GetTargetComByEntityIDAndStringName(inputComponent.EntityID,"PlayerComponent");
                var playerID = playerComponent.PlayerID;
                var exist = curFrame.ServerFrame.TryGetValue(playerID,out var curInput);
                if (!exist)
                {
                    continue;
                }

                inputComponent.InputX = curFrame.ServerFrame[playerID].Input.X;
                inputComponent.InputY = curFrame.ServerFrame[playerID].Input.Y;
                inputComponent.IsJump = curFrame.ServerFrame[playerID].Input.IsJump;
            }
        }
    }
}