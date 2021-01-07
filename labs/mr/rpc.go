package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type TaskArgs struct {
	X int
}

type TaskFinishArgs struct {
	TaskType string
	Filename string
	Index int
}

type TaskReply struct {
	TaskType string
	FileName string
	Index int
	ReduceN int
	MapperN int
}

type TaskFinishReply struct {
	Ack bool
}

// Add your RPC definitions here.


// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
