package anet

import (
	"ainx/ainterface"
	"ainx/utils"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]ainterface.IRouter //存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32                        //业务工作Worker池的数量
	TaskQueue      []chan ainterface.IRequest    //Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ainterface.IRouter),
		WorkerPoolSize: utils.GlobalSetting.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan ainterface.IRequest, utils.GlobalSetting.WorkerPoolSize),
	}
}

// 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ainterface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
// msgId
func (mh *MsgHandle) AddRouter(msgId uint32, router ainterface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

// 启动一个Woeker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ainterface.IRequest) {
	fmt.Println("Work ID =", workerID, "is started.")
	// 不断的等待队列消息
	for {
		select {
		// 从消息取出队列的Request，比执行绑定的业务方法
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

/*
启动workpool
*/
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		////给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ainterface.IRequest, utils.GlobalSetting.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

/*
将消息交给TaskQueue,由worker进行处理
目前采用轮询法则
*/
func (mh *MsgHandle) SendMsgToTaskQueue(request ainterface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
