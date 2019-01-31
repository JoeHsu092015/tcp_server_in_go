package external

import (
	"bufio"
	"net"
	"net/textproto"
	"time"
)

// ExternalServer - external server interface
type ExternalServer interface {
	Connect() error
	SendMessage(msg string) (string, error)
	Close()
}

// MyExternalServer - connect to external server
type MyExternalServer struct {
	Addr         string
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	conn         net.Conn
	reader       *textproto.Reader
}

// Connect - connect to external server
func (s *MyExternalServer) Connect() error {
	var err error
	s.conn, err = net.Dial("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.reader = textproto.NewReader(bufio.NewReader(s.conn))
	return nil
}

// SendMessage - send client's query to message
func (s *MyExternalServer) SendMessage(msg string) (string, error) {

	s.conn.SetWriteDeadline(time.Now().Add(s.WriteTimeOut))
	_, err := s.conn.Write([]byte(msg))
	if err != nil {
		return "", err
	}

	s.conn.SetReadDeadline(time.Now().Add(s.ReadTimeOut))
	line, err := s.reader.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

// Close - close external server connection
func (s *MyExternalServer) Close() {
	s.conn.Close()
}
