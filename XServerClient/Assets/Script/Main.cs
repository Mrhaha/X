using Script.ManagerController;
using UnityEngine;
using Script.Network.NetManager;
using XFramework;

public class Main : MonoBehaviour
{
    public float speed = 5.0f;
    public float logicFrameInterval = 0.33f;

    // Start is called before the first frame update
    void Start()
    {
        ManagerController.Init();
        NetManager.Init("127.0.0.1",5000,null);
        var req = new ReqSyncFrame{Frame = 1,X = 2,Y = 3};
        NetManager.TcpConnect.SendMsg(req);
    }

    async void ProcessLogicFrame()
    {
        float moveHorizontal = Input.GetAxis("Horizontal");
        float moveVertical = Input.GetAxis("Vertical");

        
    }
    

    // Update is called once per frame
    void Update()
    {
        // logicFrameTimer += Time.deltaTime;
        // while (logicFrameTimer >= logicFrameInterval)
        // {
        //     logicFrameTimer -= logicFrameInterval;
        //     ProcessLogicFrame(); // 执行逻辑帧
        // }
    }
    
    
    
    
    
}
