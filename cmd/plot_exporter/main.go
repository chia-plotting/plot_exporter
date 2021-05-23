package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/obicons/plot_exporter/logwatch"
)

func main() {
	directoryName := flag.String("directory", "", "directory where plot logs are stored")
	flag.Parse()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL)

	logWatcher := logwatch.NewLogWatcher(*directoryName)
	logWatcher.Start()

	ticker := time.NewTicker(time.Second * 30)
	shouldShutdown := false
	for !shouldShutdown {
		select {
		case <-signalChan:
			shouldShutdown = true

		case <-ticker.C:
			progress := logWatcher.GetProgress()
			for log, prog := range progress {
				fmt.Printf("%s: %d\n", log, prog)
			}
		}
	}
	logWatcher.Shutdown()
}
