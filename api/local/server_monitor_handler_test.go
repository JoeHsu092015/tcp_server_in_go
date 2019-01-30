package local_test

import (
	"net/http"
	"net/http/httptest"
	"tcp_server_in_go/api/local"
	"testing"
)

func TestServerStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	local.ServerStatusHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "========Client INFO========\n========Server INFO========\ncurrent connections: 0\ncurrent remaining jobs: 0\ntotal processed jobs: 0\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got \n%v want \n%v",
			rr.Body.String(), expected)
	}
}
