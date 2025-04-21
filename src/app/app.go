package app

import (
	publicModels "github.com/arashrasoulzadeh/nozzle/src/app/models"
	"github.com/arashrasoulzadeh/nozzle/src/internal/io"
	"github.com/arashrasoulzadeh/nozzle/src/internal/models"
	"github.com/google/uuid"
)

type nozzle struct {
	StatusChannel chan publicModels.StatusChannelEnum
	i             *models.Inbox
	o             *models.Outbox
	fw            *models.FileWatcher
}

func createNozzle(i *models.Inbox, o *models.Outbox, fw *models.FileWatcher, StatusChannel chan publicModels.StatusChannelEnum) *nozzle {
	return &nozzle{
		i:             i,
		o:             o,
		fw:            fw,
		StatusChannel: StatusChannel,
	}
}

func (n *nozzle) Write(path string, payload []byte) {

	n.o.Compose(models.OutboxMessage{
		TempPath: "",
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

	// keep process running
	for {
		select {}
	}
}
