package metricserver

import (
	"fmt"
	"net/http"
	"sync"
)

type MetricServer struct {
	sync.Mutex
	progress map[string]uint
	addr     string
}

func NewMetricServer(addr string) *MetricServer {
	return &MetricServer{
		progress: make(map[string]uint),
		addr:     addr,
	}
}

func (m *MetricServer) Start() {
	http.HandleFunc("/metrics", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, "# TYPE plot_progress counter")

		m.Lock()
		progressMap := m.progress
		m.Unlock()

		for plotName, progress := range progressMap {
			fmt.Fprintf(writer, "plot_progress{plot=\"%s\"} %d\n", plotName, progress)
		}
	})

	go func() {
		err := http.ListenAndServe(m.addr, nil)
		if err != nil {
			panic(fmt.Sprintf("error: MetricServer.Start(): %s", err))
		}
	}()
}

func (m *MetricServer) SetProgress(progress map[string]uint) {
	m.Lock()
	m.progress = progress
	m.Unlock()
}
