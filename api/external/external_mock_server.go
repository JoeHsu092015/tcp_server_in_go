package external

type MockExternalServer struct {
}

func (s *MockExternalServer) Connect() error {
	return nil
}

func (s *MockExternalServer) SendMessage(msg string) (string, error) {
	return "received:" + msg + "\n", nil
}

func (s *MockExternalServer) Close() {
	return
}
