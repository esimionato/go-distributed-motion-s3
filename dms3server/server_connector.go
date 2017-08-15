package dms3server

import (
	"fmt"
	"go-distributed-motion-s3/dms3libs"
	"net"
	"strconv"
)

// Init configs the library and configuration for dms3server
func Init() {

	dms3libs.LoadLibConfig("/etc/distributed-motion-s3/dms3libs/dms3libs.toml")
	LoadServerConfig("/etc/distributed-motion-s3/dms3server/dms3server.toml")

	cfg := ServerConfig.Logging
	dms3libs.CreateLogger(cfg.LogLevel, cfg.LogDevice, cfg.LogLocation, cfg.LogFilename)
	StartServer(ServerConfig.ServerPort)

}

// StartServer starts the TCP server
func StartServer(ServerPort int) {

	if listener, error := net.Listen("tcp", ":"+fmt.Sprint(ServerPort)); error != nil {
		dms3libs.LogFatal(error.Error())
	} else {
		defer listener.Close()
		serverLoop(listener)
	}

}

// serverLoop starts a loop to listen for clients, spawning a separate processing thread on
// dms3client connect
//
func serverLoop(listener net.Listener) {

	for {

		if conn, err := listener.Accept(); err != nil {
			dms3libs.LogFatal(err.Error())
		} else {
			dms3libs.LogInfo("OPEN connection from: " + conn.RemoteAddr().String())
			go processClientRequest(conn)
		}

	}

}

// processClientRequest passes motion detector application state to all dms3client listeners
func processClientRequest(conn net.Conn) {

	dms3libs.LogDebug(dms3libs.GetFunctionName())
	state := DetermineMotionDetectorState()

	if _, err := conn.Write([]byte(strconv.Itoa(int(state)))); err != nil {
		dms3libs.LogInfo(err.Error())
	}

	dms3libs.LogInfo("CLOSE connection from: " + conn.RemoteAddr().String())
	conn.Close()

}
