package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/obicons/plot_exporter/logwatch"
	"github.com/obicons/plot_exporter/metricserver"
)

func main() {
	directoryName := flag.String("directory", "", "directory where plot logs are stored")
	serverAddr := flag.String("addr", "0.0.0.0:10001", "address for metrics server")
	flag.Parse()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL)

	logWatcher := logwatch.NewLogWatcher(*directoryName)
	logWatcher.Start()

	ms := metricserver.NewMetricServer(*serverAddr)
	ms.Start()

	ticker := time.NewTicker(time.Second * 30)
	shouldShutdown := false
	for !shouldShutdown {
		select {
		case <-signalChan:
			shouldShutdown = true

		case <-ticker.C:
			ms.SetProgress(logWatcher.GetProgress())
		}
	}
	logWatcher.Shutdown()
}
