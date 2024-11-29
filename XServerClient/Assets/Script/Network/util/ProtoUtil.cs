using System;
using Google.Protobuf;

namespace Script.Network.util
{
    public static class ProtoUtil
    {
        private static bool _crcTableInitialized = false;
        private const uint CrcPoly = 0x04c11db7;
        private static readonly uint[] CrcTable = new uint[256];

        
        // 初始化 CRC32 查找表
        private static void InitCRCTable()
        {
            if (_crcTableInitialized)
            {
                return;
            }

            uint c;
            for (uint i = 0; i < 256; i++)
            {
                c = i << 24;
                for (uint j = 8; j > 0; j--)
                {
                    if ((c & 0x80000000) != 0)
                    {
                        c = (c << 1) ^ CrcPoly;
                    }
                    else
                    {
                        c = (c << 1);
                    }
                }
                CrcTable[i] = c;
            }

            _crcTableInitialized = true;
        }
        
        // 将字符串转换为32位CRC32哈希值
        private static uint StringHash(string s)
        {
            InitCRCTable();

            UInt32 hash = 0;
            foreach (char c in s)
            {
                UInt32 b = c;
                hash = ((hash >> 8) & 0x00FFFFFF) ^ CrcTable[(hash ^ b) & 0x000000FF];
            }
            return hash;
        }
        
        
        public static string GetProtoFullStringName(IMessage msg)
        {
             return msg.Descriptor.FullName;
        }
        
        public static UInt32 ProtoMsg2MsgID(IMessage msg)
        {
            var msgName = GetProtoFullStringName(msg);
            return StringHash(msgName);
        }
        
    }
}