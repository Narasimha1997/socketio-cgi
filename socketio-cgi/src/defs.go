package main

import (
	"io"
	"os/exec"
)

/*
 Standard POSIX based PIPE interface &
 definition of Process type
*/

//ProcessType definition
type ProcessType struct {
	cmd     *exec.Cmd
	ipPipe  io.WriteCloser
	opPipe  io.ReadCloser
	errPipe io.ReadCloser
}

//BoundedProcessType : A process in its running state
type BoundedProcessType struct {
	process           *ProcessType
	name              string
	pid               int
	supportEncryption bool
	output            chan []byte
	canReadBinary     bool
}

//ProcessMap : A map between socket-id and BoundedProcessType
var ProcessMap map[string]*BoundedProcessType

func getProcessPerID(socketID string) *BoundedProcessType {
	return ProcessMap[socketID]
}

func addProcessBySocketID(socketID string, process *BoundedProcessType) {
	ProcessMap[socketID] = process
}
