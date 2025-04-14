package app

import (
	"Nozzle/src/internal/models"
	"time"
)

func StartDaemon(path string) (n Nozzle, err error) {

	var statusChannel chan int

	outbox := models.NewOutbox(statusChannel)
	inbox := models.NewInbox(statusChannel)
	fileWatcher := models.NewFileWatcher(statusChannel, path, 5*time.Millisecond, inbox.Rc)
	go outbox.Run()
	go inbox.Run()
	go fileWatcher.Start()

	n = nozzle{
		i:             inbox,
		o:             outbox,
		fw:            fileWatcher,
		StatusChannel: statusChannel,
	}

	return n, err
}
