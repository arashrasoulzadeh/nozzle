// Package io /*
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package io

import (
	"bufio"
	"fmt"
	"github.com/arashrasoulzadeh/nozzle/log"
	"github.com/arashrasoulzadeh/nozzle/src/translation"
	"io"
	"os"
	"path/filepath"
)

func SaveToFile(dir string, filename string, data []byte) error {
	path := dir

	if filename != "" {
		path = dir + "/" + filename
	}

	log.Info(translation.InfoMessagesSavingFile, path)

	err := CreateDirsIfNotExists(path)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("failed to close file", err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return writer.Flush()

}

func DeleteFile(path string) error {
	log.Info(translation.InfoMessagesDeletingFile, path)
	return os.Remove(path)
}

func LoadFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("failed to close file", err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func CreateDirsIfNotExists(fp string) error {
	dir := filepath.Dir(fp)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directories: %v", err)
	}

	_, err = os.OpenFile(fp, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	return nil
}
