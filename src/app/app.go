package app

import (
	publicModels "Nozzle/src/app/models"
	"Nozzle/src/internal/io"
	"Nozzle/src/internal/models"
	"github.com/google/uuid"
)

type Nozzle struct {
	StatusChannel chan publicModels.StatusChannelEnum
	i             *models.Inbox
	o             *models.Outbox
	fw            *models.FileWatcher
}

func NewNozzle(i *models.Inbox, o *models.Outbox, fw *models.FileWatcher, StatusChannel chan publicModels.StatusChannelEnum) *Nozzle {
	return &Nozzle{
		i:             i,
		o:             o,
		fw:            fw,
		StatusChannel: StatusChannel,
	}
}

func (n *Nozzle) Write(path string, payload []byte) {

	n.o.Compose(models.OutboxMessage{
		TempPath: "",
		File:     models.CreateFile(uuid.New().String(), path, payload),
		Status:   "test",
	})
}
func (n *Nozzle) Read(path string) ([]byte, error) {
	return io.LoadFromFile(path)
}
