// Package models
//
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>

package models

import (
	"Nozzle/log"
	publicModels "Nozzle/src/app/models"
	"Nozzle/src/internal/io"
	"Nozzle/src/translation"
	"sync"
)

type OutboxMessage struct {
	File     File
	Status   string
	TempPath string
}

type Outbox struct {
	wg            sync.WaitGroup
	mu            sync.Mutex
	Q             []OutboxMessage
	c             chan OutboxMessage
	StatusChannel chan publicModels.StatusChannelEnum
}

func NewOutbox(statusChannel chan publicModels.StatusChannelEnum) *Outbox {
	return &Outbox{
		Q:             make([]OutboxMessage, 0),
		mu:            sync.Mutex{},
		c:             make(chan OutboxMessage, 1),
		StatusChannel: statusChannel,
	}
}

func (o *Outbox) Run() {
	o.wg.Wait()
}

// Compose appends a message to the Outbox
func (o *Outbox) Compose(m OutboxMessage) {
	o.wg.Add(1)
	go ComposeInBackground(o, m)
}

func ComposeInBackground(o *Outbox, m OutboxMessage) {
	defer o.wg.Done()

	o.mu.Lock()
	defer o.mu.Unlock()

	for _, v := range o.Q {
		if v.File.md5 == m.File.md5 && v.File.path == m.File.path {
			log.Info(translation.InfoMessagesDuplicate, v.File)
			return
		}
	}
	o.Q = append(o.Q, m)

	// Create binary representation of OutboxMessage
	messageBytes, err := m.File.MarshalBinary()
	if err != nil {
		log.Error(translation.InfoMessagesCannotSerialize, m.File)
		return
	}

	// Save with .noz extension
	err = io.SaveToFile("./temp/", m.File.uuid+".noz", messageBytes)
	if err != nil {
		log.Error(translation.InfoMessagesCannotSaveTempFile, m.File)
		return
	}

	log.Info(translation.InfoMessagesComposed, m.File.path)
}

func (o *Outbox) Consume(messages []OutboxMessage) *Outbox {
	// You can process messages here if needed
	return o
}

func (o *Outbox) Channel() chan OutboxMessage {
	return o.c
}
