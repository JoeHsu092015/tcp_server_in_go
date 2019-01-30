package local

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"sync"
	"tcp_server_in_go/api/external"
	"time"
)

const defaultMessage = `=================================
HELLO TCP server
send 'quit' for close connection
=================================
`

// IndexHandler - TCP server handler
func IndexHandler(monitorQueue chan MonitorMetric, client net.Conn, server external.ExternalServer) {

	var clientQuery string
	var err error
	var serverResponse string
	var w sync.WaitGroup

	// monitor metric
	var processReq int64
	var waitingProcess bool
	var startTime time.Time

	clientConnecting := false

	if client == nil || server == nil {
		return
	}
	log.Println("connect start:", client.LocalAddr().String())

	// external API connection establish
	err = server.Connect()
	if err != nil {
		client.Write([]byte("external server connect fail"))
		log.Println("external server connect fail")
		client.Close()
		return
	}
	defer server.Close()

	// record client connection status
	clientConnecting = true
	defer client.Close()

	// record client connect time
	startTime = time.Now()

	// start a monitor and send metric to monitor queue every second
	w.Add(1)
	go func(q chan MonitorMetric) {
		defer w.Done()
		for {
			q <- MonitorMetric{Addr: client.RemoteAddr().String(), StartTime: startTime,
				ProcessedReq: processReq, WaitProcess: waitingProcess, Alive: true}
			time.Sleep(1 * time.Second)
			// if client disconnect stop gathering data
			if !clientConnecting {
				q <- MonitorMetric{Addr: client.RemoteAddr().String(), Alive: false, ProcessedReq: processReq}
				break
			}
		}
	}(monitorQueue)

	sendDefaultBanner(client)
	tp := textproto.NewReader(bufio.NewReader(client))
	// start listen client requests
	for {
		client.Write([]byte("> "))
		clientQuery, err = tp.ReadLine()
		if err != nil {
			break
		}

		// client finish query jobs
		if clientQuery == "quit" {
			client.Write([]byte("bye~"))
			break
		}

		waitingProcess = true
		serverResponse, err = server.SendMessage(clientQuery)
		waitingProcess = false
		processReq++
		if err != nil {
			client.Write([]byte("send query to external server failed"))
			break
		}
		client.Write([]byte(serverResponse))
	}
	clientConnecting = false
	log.Println(client.LocalAddr().String(), " disconnected")
	//wait client's monitor finish jobs
	w.Wait()
}

// sendDefaultBanner - send default message when client connection established
func sendDefaultBanner(conn net.Conn) error {
	conn.Write([]byte(defaultMessage))
	return nil
}
