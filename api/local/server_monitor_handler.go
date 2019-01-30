package local

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// MonitorMetric - client info metric
type MonitorMetric struct {
	Addr         string
	StartTime    time.Time
	ProcessedReq int64
	WaitProcess  bool
	Alive        bool
}

// ServerInfo - server info metric
type ServerInfo struct {
	ClientList   map[string]MonitorMetric
	Lock         sync.RWMutex
	ProcessedJob int64
}

var (
	// MonitorQueue - send client metric to server monitor from this channel
	MonitorQueue chan MonitorMetric
	// serverStatus - current server status with client info
	serverStatus ServerInfo
)

func init() {
	MonitorQueue = make(chan MonitorMetric, 5)
	serverStatus.ClientList = make(map[string]MonitorMetric)
	serverStatus.ProcessedJob = 0
	go monitorListener()
}

// ServerStatusHandler - show server metric gauge
func ServerStatusHandler(w http.ResponseWriter, r *http.Request) {

	var s strings.Builder
	fmt.Fprint(&s, "========Client INFO========\n")
	clientCount := 0

	var totalRemainingJobs int64
	var remaingJob int

	serverStatus.Lock.RLock()
	totalProcessedJobs := serverStatus.ProcessedJob
	// clients status statistic
	for _, info := range serverStatus.ClientList {
		clientCount++
		remaingJob = 0
		if info.WaitProcess {
			remaingJob = 1
			totalRemainingJobs++
		}

		fmt.Fprint(&s, "client: ", clientCount, "\n")
		fmt.Fprint(&s, "client addr: ", info.Addr, "\n")
		fmt.Fprint(&s, "client connect time: ", info.StartTime.Format(time.RFC3339), "\n")
		fmt.Fprint(&s, "client connect duration: ", int(time.Now().Sub(info.StartTime).Seconds()), "s\n")
		fmt.Fprint(&s, "processed requests: ", info.ProcessedReq, "\n")

		totalProcessedJobs += info.ProcessedReq
		fmt.Fprint(&s, "request rate: ",
			fmt.Sprintf("%.2f",
				float64(info.ProcessedReq+int64(remaingJob))/
					time.Now().Sub(info.StartTime).Seconds()), "/s\n")
		fmt.Fprint(&s, "\n")
	}
	serverStatus.Lock.RUnlock()
	// server status statistic
	fmt.Fprint(&s, "========Server INFO========\n")
	fmt.Fprint(&s, "current connections: ", clientCount, "\n")
	fmt.Fprint(&s, "current remaining jobs: ", totalRemainingJobs, "\n")
	fmt.Fprint(&s, "total processed jobs: ", totalProcessedJobs, "\n")
	w.Write([]byte(s.String()))
}

// monitorListener - gather client's status
func monitorListener() {
	for {
		select {
		case data := <-MonitorQueue:
			updateClientInfo(data)
		}
	}
}

// updateClientInfo - update client's status
func updateClientInfo(data MonitorMetric) {
	serverStatus.Lock.Lock()
	defer serverStatus.Lock.Unlock()
	if data.Alive {
		serverStatus.ClientList[data.Addr] = data
	} else {
		if _, ok := serverStatus.ClientList[data.Addr]; ok {
			// gather finished clients processed jobs
			serverStatus.ProcessedJob += data.ProcessedReq
			delete(serverStatus.ClientList, data.Addr)
		}
	}
}
