// Package tests /*
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package tests

import (
	"encoding/base64"
	"fmt"
	"github.com/arashrasoulzadeh/nozzle/src/app"
	publicModels "github.com/arashrasoulzadeh/nozzle/src/app/models"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type File struct {
	uuid    string
	path    string
	payload []byte
	md5     string
}

func (f *File) MarshalBinary() ([]byte, error) {
	return []byte(fmt.Sprintf("%s|%s|%s|%s", f.uuid, f.path, base64.StdEncoding.EncodeToString(f.payload), f.md5)), nil
}

func TestOutbox(t *testing.T) {
	path := "./temp"
	os.RemoveAll(path)
	err := os.MkdirAll(path, os.ModePerm)
	assert.NoError(t, err)

	f := File{uuid: "Test", path: "/tmp", payload: []byte("Test"), md5: "test"}

	for i := 0; i < 3; i++ {
		f.path = "/tmp/" + strconv.Itoa(i) + ".txt"
		filePath := filepath.Join(path, fmt.Sprintf("test_%d.noz", i))

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}

		b, err := f.MarshalBinary()
		assert.NoError(t, err)

		_, err = file.Write(b)

		file.Close()
	}

	n, e := app.StartDaemon("temp")
	assert.NoError(t, e)
	n.Write("/tmp/arash.txt", []byte("test"))
	n.Write("/tmp/arash2.txt", []byte("test"))
	n.Write("/tmp/arash3.txt", []byte("test"))
	n.Write("/tmp/arash4.txt", []byte("test"))

	body, err := n.Read("/tmp/arash.txt")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), body)

loop:
	for {
		select {
		case status := <-n.StatusChannel:
			fmt.Println("received status", status)

			if status == publicModels.StatusChannelExit {
				break loop
			}
		}
	}
}
