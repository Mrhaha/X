using System;
using System.Collections.Generic;
using UnityEngine;

namespace Script.LogicFrame.Entity
{
    public static class EntityManager
    {
        private static Int32 UniqueID = 0;
        private static Dictionary<Int32, Entity> _uniqueID2Entity;
        public static Dictionary<Int32, Entity> UniqueID2Entity => _uniqueID2Entity;
        public static Dictionary<int, Int32> _gameInstanceID2EntityID; 

        public static Entity CreateEntity()
        {
            var idx = UniqueID++;
            var entity = new Entity(idx);
            _uniqueID2Entity[idx] = entity;
            return entity;
        }

        public static void BindGameObject(GameObject obj,Int32 entityID)
        {
            _gameInstanceID2EntityID[obj.GetInstanceID()] = entityID;
        }

        public static Int32 GetEntityIDByGameObject(GameObject obj)
        {
            if (_gameInstanceID2EntityID.TryGetValue(obj.GetInstanceID(), out var entityID))
            {
                return entityID;
            }
            return 0;
        }
    }
}