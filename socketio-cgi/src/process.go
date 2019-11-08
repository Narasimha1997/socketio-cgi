package main

import (
	"bufio"
	"os/exec"
)

func checkErrors(err error) {
	if err != nil {
		panic("Error : " + err.Error())
	}
}

func createChildProcessByCommand(command string, args []string) *ProcessType {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	checkErrors(err)

	stdin, err := cmd.StdinPipe()
	checkErrors(err)

	stderr, err := cmd.StderrPipe()
	checkErrors(err)

	checkErrors(cmd.Start())

	return &ProcessType{cmd, stdin, stdout, stderr}
}

func createNewBoundedProcess(name string, pid int, supportEncryption bool, process *ProcessType, binaryMode bool) *BoundedProcessType {
	return &BoundedProcessType{
		process:           process,
		name:              name,
		pid:               pid,
		supportEncryption: supportEncryption,
		output:            make(chan []byte),
		canReadBinary:     binaryMode,
	}
}

//readers Stdin, Stdout, Stderr :
func (p *BoundedProcessType) writeToStdin(data []byte) {
	p.process.ipPipe.Write(data)
}

func (p *BoundedProcessType) probeOutputs() {
	//TODO: Enable error logging later

	if !p.canReadBinary { //text stream reader, run if not binary mode
		for {
			bufferStream := bufio.NewReader(p.process.opPipe)
			byteData, error := bufferStream.ReadBytes('\n')
			if error == nil {
				p.output <- trimEOL(byteData)
			} else {
				break
			} // handle errors when Logger is developed
		}
		close(p.output)
	} else {
		byteBuffer := make([]byte, 1024*1024)
		for {
			size, err := p.process.opPipe.Read(byteBuffer)
			if err != nil {
				break
			}
			//copy the data :
			p.output <- append(make([]byte, 0, size), byteBuffer[:size]...)
		}

		close(p.output)
	}
}

func (p *BoundedProcessType) readOutput() chan []byte {
	return p.output
}

func (p *BoundedProcessType) startProbing() {
	go p.probeOutputs()
}

//from websocketd , check github
func trimEOL(b []byte) []byte {
	lns := len(b)
	if lns > 0 && b[lns-1] == '\n' {
		lns--
		if lns > 0 && b[lns-1] == '\r' {
			lns--
		}
	}
	return b[:lns]
}
