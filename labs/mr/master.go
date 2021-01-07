package mr

import (
	"log"
	"sync"
)
import "net"
import "os"
import "net/rpc"
import "net/http"


var mapTaskLock = new(sync.Mutex)
var reduceTaskLock = new(sync.Mutex)

type TaskStatus struct {
	index int // task index
	status int // value 0 is un-started, 1 is started
}

type Master struct {
	// Your definitions here.
	nReduce int
	nMapper int
	// map tasks
	mapTask map[string]TaskStatus // key is for filename;
	// reduce tasks
	reduceTask map[int]TaskStatus // key is for index of reduce task
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) AskForTask(args *TaskArgs, reply *TaskReply) error {
	hasMapTask := m.getMapTask(args, reply)
	if hasMapTask {
		// continue with map task
		return nil
	}
	hasReduceTask := m.getReduceTask(args, reply)
	if hasReduceTask {
		// continue with reduce task
		return nil
	}
	reply.TaskType = "exit"
	return nil
}

func (m *Master) TaskFinish(args *TaskFinishArgs, reply *TaskFinishReply) error {
	if args.TaskType == "map" {
		fileName := args.Filename
		mapTaskLock.Lock()
		delete(m.mapTask, fileName)
		mapTaskLock.Unlock()
	} else {
		reduceIndex := args.Index
		reduceTaskLock.Lock()
		delete(m.reduceTask, reduceIndex)
		reduceTaskLock.Unlock()
	}
	return nil
}

func (m *Master) getMapTask(args * TaskArgs, reply * TaskReply) bool {
	var mapFileName = ""
	var mapIndex = -1
	mapTaskLock.Lock()
	defer mapTaskLock.Unlock()
	// the un-started task is first order
	// to check any un-started(status is 0) task
	// assign un-started task to worker
	for key, taskStatus := range m.mapTask{
		if taskStatus.status == 0 {
			mapFileName, mapIndex = key, taskStatus.index
			m.mapTask[mapFileName] = TaskStatus{status: 1,
																					index: taskStatus.index}
			break
		}
	}
	if mapFileName == "" {
		// then consider started task,
		// the task has been assigned, but the worker is slow
		// so assign it to another worker
		for key, taskStatus := range m.mapTask{
			if taskStatus.status == 1 {
				mapFileName, mapIndex = key, taskStatus.index
				break
			}
		}
	}
	// no any task remain, just return false
	if mapFileName == "" {
		return false
	}
	// init the task parameters
	reply.TaskType = "map"
	reply.FileName = mapFileName
	reply.Index = mapIndex
	reply.ReduceN = m.nReduce
	return true
}

func (m *Master) getReduceTask(args *TaskArgs, reply *TaskReply) bool {
	var reduceIndex = -1
	reduceTaskLock.Lock()
	defer reduceTaskLock.Unlock()
	// the strategy to assign map task is same to reduce task
	for key, taskStatus := range m.reduceTask{
		if taskStatus.status == 0 {
			reduceIndex = key
			m.reduceTask[key] = TaskStatus{index: reduceIndex, status: 1}
			break
		}
	}
	if reduceIndex == -1{
		for key, taskStatus := range m.reduceTask {
			if taskStatus.status == 1 {
				reduceIndex = key
				break
			}
		}
	}
	if reduceIndex == -1 {
		return false
	}
	reply.Index = reduceIndex
	reply.TaskType = "reduce"
	reply.MapperN = m.nMapper
	return true
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	// Your code here.
	mapTaskLock.Lock()
	reduceTaskLock.Lock()
	defer mapTaskLock.Unlock()
	defer reduceTaskLock.Unlock()
	return len(m.mapTask) == 0 && len(m.reduceTask) == 0
}
//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	mapTask := make(map[string]TaskStatus)
	for i,file := range files {
		mapTask[file] = TaskStatus{index: i, status: 0}
	}
	reduceTask := make(map[int]TaskStatus)
	for i := 0; i < nReduce; i++ {
		reduceTask[i] = TaskStatus{index: i, status: 0}
	}
	m := Master{
		nMapper: len(files),
		nReduce: nReduce,
		mapTask: mapTask,
		reduceTask: reduceTask,
	}
	m.server()
	return &m
}
