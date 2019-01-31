package external

// MockExternalServer - mock external server
type MockExternalServer struct {
}

// Connect - mock connect
func (s *MockExternalServer) Connect() error {
	return nil
}

// SendMessage - mock ssend message, default received:msg
func (s *MockExternalServer) SendMessage(msg string) (string, error) {
	return "received:" + msg + "\n", nil
}

// Close - mock close external server connection
func (s *MockExternalServer) Close() {
	return
}
