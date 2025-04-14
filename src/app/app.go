package app

import (
	"Nozzle/src/internal/models"
)

type nozzle struct {
	StatusChannel chan int
	i             *models.Inbox
	o             *models.Outbox
	fw            *models.FileWatcher
}

func NewNozzle(i *models.Inbox, o *models.Outbox, fw *models.FileWatcher) *nozzle {
	return &nozzle{
		i:  i,
		o:  o,
		fw: fw,
	}
}

type Nozzle interface {
	Write(path string, paylaod []byte)
}

func (n nozzle) Write(path string, paylaod []byte) {
	n.o.Compose(models.OutboxMessage{
		TempPath: "",
		File:     models.CreateFile("test", path, paylaod),
		Status:   "test",
	})
}
