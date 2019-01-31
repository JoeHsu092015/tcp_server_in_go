package main

import (
	"log"
	"net"
	"net/http"
	"tcp_server_in_go/api/external"
	"tcp_server_in_go/api/local"
)

const (
	tcpAddr  = "localhost:7788"
	httpAddr = "localhost:7799"
)

func main() {

	l, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()
	go httpServer()

	log.Println("TCP server start listen at: ", tcpAddr)
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go local.IndexHandler(local.MonitorQueue, client, getExternalServer())
	}

}

// getExternalServer - create external server
func getExternalServer() external.ExternalServer {
	//return &external.MyExternalServer{Addr: "http://192.168.1.153:80/tag", ReadTimeOut: 10 * time.Second, WriteTimeOut: 10 * time.Second}
	return &external.MockExternalServer{}
}

// httpServer - server status http endpoint
func httpServer() {

	http.HandleFunc("/", local.ServerStatusHandler)
	log.Println("HTTP server start listen at: ", httpAddr)
	err := http.ListenAndServe(httpAddr, nil)
	if err != nil {
		log.Println(err)
	}
}
