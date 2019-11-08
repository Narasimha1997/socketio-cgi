# socketio-cgi
A cgi layer that allows any process / program to stream data across world wide web.


### How it works ?
The core idea behind this project is to use POSIX Pipes to grab `stdout` from a child process and streaming them over a socket.io interface.
Clients can connect to the server using socket.io client and grab the streamed outputs.

### Endpoints:

`/probe/` - Main socket.io route
`/probe/  event name = stdout message = any` - Make a ping to /probe/ with stdout as event to start streaming.
`event listner = probes` = All the stdout data is streamed using this event.

### How to compile ?

To compile the program, provide +x permission to `build.sh`. Now run :
`./build.sh or $sh build.sh`

### How to run a program with socketio-cgi support ?
Run the program inside socketio-cgi following command : 
`./build_object program_file [...args]`

For example : Run sample test.sh program
`./build_object sh test.sh`

If you encounter any issues while running along with python program, its because of the output buffering of the python interpteter.
Disable buffering or call `sys.stdout.flush()` to manually flush the buffer. One best solution is to use `print()` as it automatically
clears the buffer.

#### This is the initial version and might have bugs. More functionalities and bugs will be fixed later.
