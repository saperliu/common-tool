package socket_io

import (
	"github.com/googollee/go-socket.io"
	"github.com/saperliu/common-tool/logger"
	"net/http"
	"strconv"
)

type SocketServer struct {
	Path       string
	Port       int
	socketConn []socketio.Conn
}

func (socketServer *SocketServer) StartServer() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		logger.Error("start  SocketIO Server error  %s", err)
	}

	server.OnConnect(socketServer.Path, func(s socketio.Conn) error {
		s.SetContext("")
		logger.Info(" SocketIO Server connected %v    %s", s.RemoteAddr(), s.ID())
		s.Emit("realtime_data", "point:234")
		return nil
	})

	//server.OnEvent("/iot/websocket", "realtime_data", func(s socketio.Conn, msg string) {
	//	fmt.Println("notice:", msg)
	//	s.Emit("reply", "have "+msg)
	//	for {
	//
	//		s.Emit("reply", "have "+msg)
	//	}
	//})

	//server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	//	s.SetContext(msg)
	//	return "recv " + msg
	//})

	//server.OnEvent("/", "bye", func(s socketio.Conn) string {
	//	last := s.Context().(string)
	//
	//	s.Emit("bye", last)
	//	s.Close()
	//	return last
	//})

	server.OnError(socketServer.Path, func(s socketio.Conn, e error) {
		logger.Info(" SocketIO Server OnError %v    %s", s.RemoteAddr(), e)
	})

	server.OnDisconnect(socketServer.Path, func(s socketio.Conn, reason string) {
		logger.Info(" SocketIO Server OnDisconnect  %v  %s", s.RemoteAddr(), reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle(socketServer.Path, server)
	logger.Info(" SocketIO Server Port  %v  %s", socketServer.Port, socketServer.Path)
	err = http.ListenAndServe(":"+strconv.Itoa(socketServer.Port), nil)
	logger.Error("start  SocketIO Server error  %s", err)
}
