package main

import (
	"fmt"
	"net/http"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

var sockets []socketio.Conn = make([]socketio.Conn, 0, 100)

func customPainc(err error) {
	if err != nil {
		panic(err)
	}
}

func customBroadCast(data string) {
	for i := 0; i < len(sockets); i++ {
		socket := sockets[i]
		socket.Emit("probes", data)
	}
}

func beginStream(process *BoundedProcessType) {
	for {
		channel := process.output
		data := string(<-channel)

		customBroadCast(data)
	}
}

func initServer(argv []string) {
	if len(argv) == 0 {
		fmt.Println("No process specified! Exiting")
		return
	} else if len(argv) == 1 {
		fmt.Printf("Attempting to run %s without arguments\n", argv[0])
	}

	server, err := socketio.NewServer(nil)
	customPainc(err)

	//as of now allow the socket to just stream the output :
	fmt.Print(argv[1:])
	processPipe := createChildProcessByCommand(argv[0], argv[1:])

	boundedProcess := createNewBoundedProcess(argv[0], 2, false, processPipe, false)

	boundedProcess.startProbing()

	defer fmt.Printf("Process %s terminated itself.\n", boundedProcess.name)

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Printf("A new client " + s.ID() + " connected")
		s.Emit("connect_"+s.ID(), "socket_id="+s.ID()+";pid="+string(boundedProcess.pid))
		sockets = append(sockets, s)
		return nil
	})

	server.OnEvent("/", "stdout", func(s socketio.Conn, message string) {

		//connect to channel and enit the output :
		fmt.Printf("Started probing for %s with PID : %d\n", s.ID(), boundedProcess.pid)
		go beginStream(boundedProcess)
		defer fmt.Printf("Finished streaming for %s\n", s.ID())
	})

	server.OnEvent("/", "stdin", func(s socketio.Conn, message string) {
		stdinString := strings.Split(message, "=")
		if len(stdinString) == 0 || len(stdinString) == 1 {
			s.Emit("probes", "Invalid stdin")
		} else {
			//append a \n
			boundedProcess.writeToStdin([]byte(stdinString[1]))
			s.Emit("probes", "Wrote stdin")
		}
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/probe/", server)
	http.Handle("/", http.FileServer(http.Dir("./test")))
	http.ListenAndServe(":8000", nil)
}
