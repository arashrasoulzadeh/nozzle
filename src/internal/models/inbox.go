package models

import (
	"Nozzle/log"
	publicModels "Nozzle/src/app/models"
	"Nozzle/src/internal/io"
	"Nozzle/src/translation"
	"encoding/base64"
	"os"
	"syscall"
)

type Inbox struct {
	File          File
	Rc            chan OutboxMessage
	statusChannel chan publicModels.StatusChannelEnum
}

func NewInbox(statusChannel chan publicModels.StatusChannelEnum) *Inbox {
	return &Inbox{
		Rc:            make(chan OutboxMessage, 10),
		statusChannel: statusChannel,
	}
}
func (i *Inbox) WriteFile(msg OutboxMessage) error {
	decodedPayload, err := base64.StdEncoding.DecodeString(string(msg.File.payload))
	if err != nil {
		log.Error("Failed to decode payload:", err, string(msg.File.payload))
		return err
	}
	err = io.SaveToFile(msg.File.path, "", decodedPayload)
	if err != nil {
		log.Error(translation.InfoMessagesCannotSaveFile, err)
		return err
	}
	err = i.DeleteTemp(msg)
	if err != nil {
		log.Error(translation.InfoMessagesCannotDeleteFile)
		return err
	}
	i.statusChannel <- publicModels.StatusChannelFileWritten

	return nil
}

func (i *Inbox) DeleteTemp(msg OutboxMessage) error {
	return io.DeleteFile(msg.TempPath)
}

func (i *Inbox) Run() {
	for {
		select {
		case msg := <-i.Rc:
			file, err := os.OpenFile(msg.TempPath, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				log.Error(translation.InfoMessagesCannotOpenFile, err)
				continue
			}

			if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
				log.Error(translation.InfoMessagesCannotLockFile, err)
				continue
			}
			err = i.WriteFile(msg)
			if err != nil {
				log.Error(translation.InfoMessagesCannotSaveFile, err)
				continue
			}
			if err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
				log.Error(translation.InfoMessagesCannotUnLockFile, err)
				continue
			}
			//case <-time.After(2 * time.Second):
			//	fmt.Println("Closing Inbox Channel")
		}
	}
}
