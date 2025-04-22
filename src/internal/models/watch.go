// Package models
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/arashrasoulzadeh/nozzle/log"
	publicModels "github.com/arashrasoulzadeh/nozzle/src/app/models"
	"github.com/arashrasoulzadeh/nozzle/src/translation"
)

type FileWatcher struct {
	dir           string
	seen          map[string]os.FileInfo
	interval      time.Duration
	events        chan FileEvent
	stopChan      chan struct{}
	SendChannel   chan OutboxMessage
	statusChannel chan publicModels.StatusChannelEnum
}

type FileEvent struct {
	Type string // "create" or "delete"
	Name string
}

func NewFileWatcher(sendChannel chan OutboxMessage, statusChannel chan publicModels.StatusChannelEnum, dir string, interval time.Duration) *FileWatcher {
	return &FileWatcher{
		dir:           dir,
		seen:          make(map[string]os.FileInfo),
		interval:      interval,
		events:        make(chan FileEvent, 10),
		stopChan:      make(chan struct{}),
		SendChannel:   sendChannel,
		statusChannel: statusChannel,
	}
}

func (fw *FileWatcher) Events() <-chan FileEvent {
	return fw.events
}

func (fw *FileWatcher) Start() {
	go func() {
		ticker := time.NewTicker(fw.interval)
		defer ticker.Stop()
		for {
			select {
			case <-fw.stopChan:
				return
			case <-ticker.C:
				fw.scan()
			}
		}
	}()
}

func (fw *FileWatcher) Stop() {
	close(fw.stopChan)
}

func (fw *FileWatcher) SendPendingToChannel() error {
	entries, err := os.ReadDir(fw.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fullPath := filepath.Join(fw.dir, entry.Name())

		fileContents, err := os.ReadFile(fullPath)
		if err != nil || len(fileContents) == 0 {
			continue
		}

		file := CreateFile(entry.Name(), fullPath, fileContents)
		if err := file.UnmarshalBinary(fileContents); err != nil {
			log.Error(translation.InfoMessagesCannotDeSerialize, err)
			continue
		}

		fw.SendChannel <- OutboxMessage{
			File:     file,
			Status:   "new",
			TempPath: fullPath,
		}
	}

	return nil
}

func (fw *FileWatcher) scan() {
	currentFiles := make(map[string]os.FileInfo)

	entries, err := os.ReadDir(fw.dir)
	if err != nil {
		log.Error(translation.InfoMessagesCannotWatchDirectory, fw.dir)
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		currentFiles[entry.Name()] = info

		if _, found := fw.seen[entry.Name()]; !found {
			fw.handleNewFile(entry.Name())
		}
	}

	for name := range fw.seen {
		if _, found := currentFiles[name]; !found {
			fw.events <- FileEvent{Type: "delete", Name: name}
		}
	}

	fw.seen = currentFiles
}

func (fw *FileWatcher) handleNewFile(fileName string) {
	fullPath := filepath.Join(fw.dir, fileName)

	fileContents, err := os.ReadFile(fullPath)
	if err != nil || len(fileContents) == 0 {
		return
	}

	file := CreateFile(fileName, fullPath, fileContents)
	if err := file.UnmarshalBinary(fileContents); err != nil {
		log.Error(translation.InfoMessagesCannotDeSerialize, err)
		return
	}

	fw.SendChannel <- OutboxMessage{
		File:     file,
		Status:   "new",
		TempPath: fullPath,
	}
	fw.events <- FileEvent{Type: "create", Name: fileName}
	log.Info(translation.InfoMessagesFileDetected, fileName)
}
