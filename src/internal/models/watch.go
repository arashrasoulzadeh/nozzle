// Package models
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package models

import (
	"fmt"
	"github.com/arashrasoulzadeh/nozzle/log"
	publicModels "github.com/arashrasoulzadeh/nozzle/src/app/models"
	"github.com/arashrasoulzadeh/nozzle/src/translation"
	"os"
	"time"
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
		for {
			select {
			case <-fw.stopChan:
				return
			default:
				fw.scan()
				time.Sleep(fw.interval)
			}
		}
	}()
}

func (fw *FileWatcher) Stop() {
	close(fw.stopChan)
}

func (fw *FileWatcher) SendPendingToChannel(path string, statusChannel chan publicModels.StatusChannelEnum) error {
	entries, err := os.ReadDir(fw.dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		fmt.Println(entry.Name(), entry.IsDir())
		if !entry.IsDir() {
			fileContents, err := os.ReadFile(fw.dir + "/" + entry.Name())
			if err != nil {
				continue
			}

			file := CreateFile(entry.Name(), fw.dir+"/"+entry.Name(), fileContents)

			// Unmarshal binary data back into file
			if err := file.UnmarshalBinary(fileContents); err != nil {
				log.Error(translation.InfoMessagesCannotDeSerialize, err)
				continue
			}
			om := OutboxMessage{
				File:     file,
				Status:   "new",
				TempPath: fw.dir + "/" + entry.Name(),
			}

			fw.SendChannel <- om
		}
	}
	return nil
}

func (fw *FileWatcher) scan() {
	currentFiles := make(map[string]os.FileInfo)

	entries, err := os.ReadDir(fw.dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		log.Error(translation.InfoMessagesCannotWatchDirectory, fw.dir)

		return
	}

	// Detect new and existing files
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		currentFiles[entry.Name()] = info

		if _, found := fw.seen[entry.Name()]; !found {
			fw.events <- FileEvent{Type: "create", Name: entry.Name()}

			// Read file contents
			fileContents, err := os.ReadFile(fw.dir + "/" + entry.Name())
			if err != nil {
				continue
			}

			// Create file object with actual contents
			file := CreateFile(entry.Name(), fw.dir+"/"+entry.Name(), fileContents)

			// Unmarshal binary data back into file
			if err := file.UnmarshalBinary(fileContents); err != nil {
				log.Error(translation.InfoMessagesCannotDeSerialize, err)
				continue
			}

			fw.SendChannel <- OutboxMessage{
				File:     file,
				Status:   "new",
				TempPath: fw.dir + "/" + entry.Name(),
			}
			log.Info(translation.InfoMessagesFileDetected, entry)
		}
	}

	// Detect deleted files
	for name := range fw.seen {
		if _, found := currentFiles[name]; !found {
			fw.events <- FileEvent{Type: "delete", Name: name}
		}
	}

	fw.seen = currentFiles
}
