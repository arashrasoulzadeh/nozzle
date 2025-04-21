package app

import (
	"github.com/arashrasoulzadeh/nozzle/log"
	publicModels "github.com/arashrasoulzadeh/nozzle/src/app/models"
	"github.com/arashrasoulzadeh/nozzle/src/internal/io"
	"github.com/arashrasoulzadeh/nozzle/src/internal/models"
	"github.com/arashrasoulzadeh/nozzle/src/translation"
	"github.com/google/uuid"
)

type nozzle struct {
	StatusChannel chan publicModels.StatusChannelEnum
	i             *models.Inbox
	o             *models.Outbox
	fw            *models.FileWatcher
	tempPath      string
}

func createNozzle(i *models.Inbox, o *models.Outbox, fw *models.FileWatcher, StatusChannel chan publicModels.StatusChannelEnum, tempPath string) *nozzle {
	return &nozzle{
		i:             i,
		o:             o,
		fw:            fw,
		StatusChannel: StatusChannel,
		tempPath:      tempPath,
	}
}

func (n *nozzle) Write(path string, payload []byte) {

	n.o.Compose(models.OutboxMessage{
		TempPath: n.tempPath,
		File:     models.CreateFile(uuid.New().String(), path, payload),
		Status:   "test",
	})
}
func (n *nozzle) Read(path string) ([]byte, error) {
	return io.LoadFromFile(path)
}

func (n *nozzle) Start() {
	go n.o.Run()
	go n.i.Run()
	go n.fw.Start()
	go n.Pending()

	// keep process running
	for {
		select {}
	}
}

func (n *nozzle) Pending() {
pendingFiles:
	for {
		log.Info(translation.InfoMessagesProcessingPendingFiles)
		err := n.fw.SendPendingToChannel(n.tempPath, n.o.InternalChannel)
		if err != nil {
			log.Error(translation.InfoMessagesCannotProcessPendingFiles, err)
			break pendingFiles
		}

		break pendingFiles
	}
}
