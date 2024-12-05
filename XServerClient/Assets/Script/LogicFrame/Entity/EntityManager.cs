using System;
using System.Collections.Generic;
using System.Linq;
using Script.LogicFrame.Components;
using Unity.VisualScripting;
using UnityEditor.SceneManagement;
using UnityEngine;

namespace Script.LogicFrame.Entity
{
    public static class EntityManager
    {
        public static Dictionary<string,bool> _ComponentNameDictionary;
        public static Dictionary<Int32, Entity> UniqueID2Entity => _uniqueID2Entity;
        private static Int32 UniqueID = 0;
        private static Dictionary<Int32, Entity> _uniqueID2Entity;
        private static Dictionary<Int32, Int32> _gameInstanceID2EntityID; 
        private static Dictionary<string, List<IComponent>> _name2Componets;

        public static Entity CreateEntity()
        {
            var idx = UniqueID++;
            var entity = new Entity(idx);
            _uniqueID2Entity[idx] = entity;
            return entity;
        }

        static EntityManager()
        {
            _gameInstanceID2EntityID = new Dictionary<Int32, Int32>();
            _name2Componets = new Dictionary<string, List<IComponent>>();
            _ComponentNameDictionary = new Dictionary<string, bool>();

            _ComponentNameDictionary.Add("InputComponent",true);
            _ComponentNameDictionary.Add("PlayerComponent",true);
        }

        //将组件的创建直接内置化，方便后续做池化管理

        public static IComponent GetTargetComByEntityIDAndStringName(Int32 entityID,string componentName)
        {
            var entity = GetEntityByID(entityID);
            return entity.GetComponentByName(componentName);
        }


        public static Entity GetEntityByID(Int32 entityID)
        {
            var exist = _uniqueID2Entity.TryGetValue(entityID, out var entity);
            return exist ? entity : null;
        }
        
        
        public static bool EntityAddComponent(Int32 entityID,string componentName)
        {
            var comExist = _ComponentNameDictionary.ContainsKey(componentName);
            if (comExist)
            {
                return false;
            }
            
            var exist = _uniqueID2Entity.TryGetValue(entityID,out var entity);
            if (!exist)
            {
                return false;
            }
            var type = Type.GetType(componentName);
            if (type == null)
            {
                return false;

            }
            
            var com = (IComponent)Activator.CreateInstance(type);
            com.EntityID = entityID;
            com.Name = componentName;
            
            entity.AddComponent(com);
            if (!_name2Componets.ContainsKey(com.Name))
            {
                _name2Componets[com.Name] = new List<IComponent>();
            }
            _name2Componets[com.Name].Add(com);
            return true;
        }
        
        public static List<IComponent> GetComponentsByStringName(string comName)
        {
            var exist = _name2Componets.TryGetValue(comName, out var components);
            if (!exist)
            {
                return null;
            }
            return components;
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