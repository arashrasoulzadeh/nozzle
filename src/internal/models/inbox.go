package models

import (
	"Nozzle/log"
	"Nozzle/src/io"
	"Nozzle/src/translation"
	"encoding/base64"
	"fmt"
)

type Inbox struct {
	File File
	Rc   chan OutboxMessage
}

func NewInbox(statusChannel chan int) *Inbox {
	return &Inbox{
		Rc: make(chan OutboxMessage, 10),
	}
}
func (i *Inbox) WriteFile(msg OutboxMessage) error {
	decodedPayload, err := base64.StdEncoding.DecodeString(string(msg.File.payload))
	if err != nil {
		log.Error("Failed to decode payload:", err)
		return err
	}
	err = io.SaveToFile(msg.File.path, "", decodedPayload)
	if err != nil {
		log.Error(translation.InfoMessagesCannotSaveFile)
		return err
	}
	err = i.DeleteTemp(msg)
	if err != nil {
		log.Error(translation.InfoMessagesCannotDeleteFile)
		return err
	}
	return nil
}

func (i *Inbox) DeleteTemp(msg OutboxMessage) error {
	return io.DeleteFile(msg.TempPath)
}

func (i *Inbox) Run() {
	for {
		select {
		case msg := <-i.Rc:
			fmt.Println("Inbox received:", msg.File.path)
			i.WriteFile(msg)
			//case <-time.After(2 * time.Second):
			//	fmt.Println("Closing Inbox Channel")
		}
	}
}
