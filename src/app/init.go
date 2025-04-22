package app

import (
	publicModels "github.com/arashrasoulzadeh/nozzle/src/app/models"
	"github.com/arashrasoulzadeh/nozzle/src/internal/models"
	"time"
)

func Nozzle(path string) (n *NozzleStruct, err error) {

	internalChannel := make(chan publicModels.StatusChannelEnum)
	statusChannel := make(chan publicModels.StatusChannelEnum)

	outbox := models.NewOutbox(internalChannel, statusChannel)
	inbox := models.NewInbox(statusChannel)
	fileWatcher := models.NewFileWatcher(inbox.ReceiveChannel, statusChannel, path, 5*time.Millisecond)

	n = createNozzle(inbox, outbox, fileWatcher, statusChannel, path)

	return n, err
}
