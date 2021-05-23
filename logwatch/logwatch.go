package logwatch

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/obicons/plot_exporter/plotinfo"
)

type LogWatcher struct {
	sync.Mutex
	directoryName  string
	completedLogs  map[string]bool
	shutdownChan   chan bool
	ackChan        chan bool
	progressStatus map[string]uint
}

const (
	refreshRate = time.Second * 15
)

/*
 * Returns a new log watcher to monitor directoryName.
 * Post: !running(NewLogWatcher(directoryName))
 */
func NewLogWatcher(directoryName string) *LogWatcher {
	return &LogWatcher{
		directoryName:  directoryName,
		completedLogs:  make(map[string]bool),
		progressStatus: make(map[string]uint),
		shutdownChan:   make(chan bool),
		ackChan:        make(chan bool),
	}
}

/*
 * Begins w's work.
 * Pre: !running(w)
 * Post: running(w)
 */
func (w *LogWatcher) Start() {
	go w.work()
}

/*
 * Stops w's work.
 * Pre: running(w)
 * Post: !running(w)
 */
func (w *LogWatcher) Shutdown() {
	w.shutdownChan <- true
	<-w.ackChan
}

/*
 * Returns a read-only view of each log file's progress.
 */
func (w *LogWatcher) GetProgress() map[string]uint {
	w.Lock()
	defer w.Unlock()
	return w.progressStatus
}

func (w *LogWatcher) work() {
	ticker := time.NewTicker(refreshRate)
	shouldShutdown := false
	for !shouldShutdown {
		select {
		case <-w.shutdownChan:
			shouldShutdown = true
		case <-ticker.C:
			w.updateProgress()
		}
	}
	ticker.Stop()
	w.ackChan <- true
}

func (w *LogWatcher) updateProgress() {
	tmpMap := make(map[string]uint)

	dirFile, err := os.Open(w.directoryName)
	if err != nil {
		panic("error: updateProgress(): cannot read log directory")
	}
	defer dirFile.Close()

	files, err := dirFile.Readdir(-1)
	if err != nil {
		log.Printf("error: updateProgress(): Readdir(): %s\n", err)
		return
	}

	for _, file := range files {
		// skip completed logs
		if w.completedLogs[file.Name()] {
			continue
		}

		reader, err := os.Open(file.Name())
		if err != nil {
			log.Printf("error: updateProgress(): Open(): %s\n", file.Name())
			continue
		}
		defer reader.Close()

		progress, done := plotinfo.GetPlotProgress(reader)
		if done {
			w.completedLogs[file.Name()] = true
		} else {
			tmpMap[file.Name()] = progress
		}
	}

	w.Lock()
	w.progressStatus = tmpMap
	w.Unlock()
}
