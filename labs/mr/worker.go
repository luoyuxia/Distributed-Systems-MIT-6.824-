package mr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)
import "log"
import "net/rpc"
import "hash/fnv"


//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}
// for sorting by key.
type ByKey []KeyValue
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the master.
	for true {
		hasWork := doTask(mapf, reducef)
		if !hasWork {
			fmt.Printf("Finish all work\n")
			break
		}
		time.Sleep(1000)
	}
}

func doTask(mapf func(string, string) []KeyValue,
						reducef func(string, []string) string) bool {
	reply, err := AskForTask()
	if err != nil || reply.TaskType == "exit" {
		return false
	}
	if reply.TaskType == "map" {
		// for map
		doMapFunc(mapf, reply.FileName, reply.Index, reply.ReduceN)
	} else {
		// for reduce
		doReduceFunc(reducef, reply.Index, reply.MapperN)
	}
	return true
}

//
// example function to show how to make an RPC call to the master.
//
// the RPC argument and reply types are defined in rpc.go.
//
func AskForTask() (TaskReply, error) {
	// declare an argument structure.
	args := TaskArgs{}
	// declare a reply structure.
	reply := TaskReply{}
	// send the RPC request, wait for the reply.
	err := call("Master.AskForTask", &args, &reply)
	// can't connect to master or the response is "exit"
	return reply, err
}



func doMapFunc(mapf func(string, string) []KeyValue,filename string,
							 mapIndex int, nReduce int) {
	var intermediate []KeyValue
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := mapf(filename, string(content))
	intermediate = append(intermediate, kva...)
	// write to files
	mapOutFilePattern := "tmr-%d-%d"
	keyValues := make(map[int][]KeyValue)
	for i := 0; i < nReduce; i ++ {
		keyValues[i] = []KeyValue{}
	}
	for _, keyval := range intermediate {
		reduceIndex := ihash(keyval.Key) % nReduce
		keyValues[reduceIndex] = append(keyValues[reduceIndex], keyval)
	}
	for i := 0; i < nReduce; i++ {
		outFile := fmt.Sprintf(mapOutFilePattern, mapIndex, i)
		fd, _:= os.OpenFile(outFile,os.O_RDWR|os.O_CREATE,0644)
		enc := json.NewEncoder(fd)
		for _, kva := range keyValues[i] {
			err := enc.Encode(&kva)
			if err != nil {
				log.Println("Fail to write kv to file.")
			}
		}
		fd.Close()
	}
	for i := 0; i < nReduce; i ++ {
		outFile := fmt.Sprintf(mapOutFilePattern, mapIndex, i)
		if _, err := os.Stat(outFile); err == nil {
			// from tmr-%d-%d to mr-%d-%d
			os.Rename(outFile, outFile[1:])
		}
	}
	finishMap(filename)
}

func finishMap(filename string)  {
	taskFishArgs := TaskFinishArgs{Filename: filename, TaskType: "map"}
	reply := TaskFinishReply{}
	err := call("Master.TaskFinish", &taskFishArgs, &reply)
	if err != nil {
		log.Printf("Fail to tell map task finish to master")
	}
}

func doReduceFunc(reducef func(string, []string) string, reduceIndex int, mapperN int)  {
	kva := []KeyValue{}
	for i := 0; i < mapperN; i ++ {
		fileName := fmt.Sprintf("mr-%d-%d", i, reduceIndex)
		fd, _:= os.Open(fileName)
		dec := json.NewDecoder(fd)
		for  {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			kva = append(kva, kv)
		}
	}
	sort.Sort(ByKey(kva))
	doReduce(kva, reducef, reduceIndex)
	finishReduce(reduceIndex)
}

func doReduce(kva []KeyValue, reducef func(string, []string) string, reduceIndex int)  {
	oname := fmt.Sprintf("tmr-out-%d", reduceIndex)
	ofile, _ := os.Create(oname)
	i := 0
	for i < len(kva) {
		j := i + 1
		for j < len(kva) && kva[j].Key == kva[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, kva[k].Value)
		}
		output := reducef(kva[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", kva[i].Key, output)
		i = j
	}
	ofile.Close()
	os.Rename(oname, oname[1:])
}

func finishReduce(reduceIndex int)  {
	taskFishArgs := TaskFinishArgs{Index: reduceIndex, TaskType: "reduce"}
	reply := TaskFinishReply{}
	err := call("Master.TaskFinish", &taskFishArgs, &reply)
	if err != nil {
		log.Printf("Fail to tell reduce task finish to master")
	}
}


//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) error {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Printf("dialing:%s", err)
		return err
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	return err
}
