using System.Security.Cryptography;

namespace Script.Network.Crypto
{
    
    public interface ICrypto
    {
        public void Encrypt(byte[] src, out byte[] dst);
        public string Decrypt(byte[] src, out byte[] dst);
    }
    
    
    
    public class AesCrypto : ICrypto
    {
        private static readonly byte[] DefaultAesKey = new byte[]
        {
            0x46, 0x72, 0x45, 0x6b, 0x55, 0x50, 0x37, 0x78,
            0x61, 0x4e, 0x3f, 0x26, 0x72, 0x65, 0x51, 0x3d,
            0x6a, 0x45, 0x66, 0x72, 0x61, 0x74, 0x68, 0x65,
            0x77, 0x35, 0x65, 0x47, 0x35, 0x51, 0x45, 0x63
        };
        
        private static readonly byte[] DefaultAesIV = new byte[]
        {
            0x73, 0x65, 0x42, 0x37, 0x24, 0x46, 0x35, 0x53,
            0x23, 0x75, 0x66, 0x61, 0x6d, 0x55, 0x6d, 0x41
        };

        private static Aes aesInstance;
        static AesCrypto()
        {
            InitializeAes();
        }
        
        private static void InitializeAes()
        {
            aesInstance = Aes.Create();
            aesInstance.Key = DefaultAesKey;
            aesInstance.IV = DefaultAesIV;
        }
        
        public void Encrypt(byte[] src, out byte[] dst)
        {
            ICryptoTransform encryptor = aesInstance.CreateEncryptor(aesInstance.Key, aesInstance.IV);

            using (var ms = new System.IO.MemoryStream())
            {
                using (var cs = new CryptoStream(ms, encryptor, CryptoStreamMode.Write))
                {
                    cs.Write(src, 0, src.Length);
                    cs.FlushFinalBlock();
                }

                dst = ms.ToArray();
            }
        }
        
        public string Decrypt(byte[] src, out byte[] dst)
        {
            ICryptoTransform decryptor = aesInstance.CreateDecryptor(aesInstance.Key, aesInstance.IV);

            using (var ms = new System.IO.MemoryStream())
            {
                using (var cs = new CryptoStream(ms, decryptor, CryptoStreamMode.Write))
                {
                    cs.Write(src, 0, src.Length);
                    cs.FlushFinalBlock();
                }

                dst = ms.ToArray();
            }

            return "";
        }
        
    }
}