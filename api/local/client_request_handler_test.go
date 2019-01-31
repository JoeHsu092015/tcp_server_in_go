package local_test

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"tcp_server_in_go/api/external"
	"tcp_server_in_go/api/local"
	"testing"
)

// TestIndexHandler - test client index
func TestIndexHandler(t *testing.T) {

	mockServer := &external.MockExternalServer{}
	tcpAddr := "localhost:4455"
	testQueue := make(chan local.MonitorMetric, 5)

	go func() {
		l, err := net.Listen("tcp", tcpAddr)
		if err != nil {
			log.Panic(err)
		}
		defer l.Close()
		for {
			client, err := l.Accept()
			if err != nil {
				log.Panic(err)
			}

			go local.IndexHandler(testQueue, client, mockServer)
		}
	}()

	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		t.Errorf("client connect fail")
		return
	}

	reader := textproto.NewReader(bufio.NewReader(conn))

	clientQuery, err := reader.ReadLine()
	expected := "================================="
	if expected != clientQuery {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			clientQuery, expected)
		return
	}
	clientQuery, err = reader.ReadLine()
	expected = "HELLO TCP server"
	if expected != clientQuery {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			clientQuery, expected)
		return
	}
	clientQuery, err = reader.ReadLine()
	expected = "send 'quit' for close connection"
	if expected != clientQuery {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			clientQuery, expected)
		return
	}
	clientQuery, err = reader.ReadLine()
	expected = "================================="
	if expected != clientQuery {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			clientQuery, expected)
		return
	}

	if err != nil {
		t.Errorf("client read fail")
		return
	}
	fmt.Println("write hello")
	conn.Write([]byte("hello\n"))
	clientQuery, err = reader.ReadLine()
	expected = "> received:hello"
	if expected != clientQuery {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			clientQuery, expected)
		return
	}
}
