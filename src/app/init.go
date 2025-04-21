package app

import (
	"github.com/arashrasoulzadeh/nozzle/log"
	publicModels "github.com/arashrasoulzadeh/nozzle/src/app/models"
	"github.com/arashrasoulzadeh/nozzle/src/internal/models"
	"github.com/arashrasoulzadeh/nozzle/src/translation"
	"time"
)

func StartDaemon(path string) (n *Nozzle, err error) {

	statusChannel := make(chan publicModels.StatusChannelEnum)

	outbox := models.NewOutbox(statusChannel)
	inbox := models.NewInbox(statusChannel)
	fileWatcher := models.NewFileWatcher(statusChannel, path, 5*time.Millisecond, inbox.Rc)

pendingFiles:
	for {
		log.Info(translation.InfoMessagesProcessingPendingFiles)
		err := fileWatcher.SendPendingToChannel(path, statusChannel)
		if err != nil {
			log.Error(translation.InfoMessagesCannotProcessPendingFiles, err)
			break pendingFiles
		}

		break pendingFiles
	}

	go outbox.Run()
	go inbox.Run()
	go fileWatcher.Start()

	n = NewNozzle(inbox, outbox, fileWatcher, statusChannel)
	go exit(n)

	return n, err
}

func exit(n *Nozzle) {
	time.Sleep(time.Second * 3)
	n.StatusChannel <- publicModels.StatusChannelExit
}
