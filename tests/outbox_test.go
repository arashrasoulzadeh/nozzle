// Package tests /*
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package tests

import (
	"Nozzle/src/app"
	"Nozzle/src/io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOutbox(t *testing.T) {
	path := "./temp"
	os.RemoveAll(path)
	io.CreateDirsIfNotExists(path)

	n, e := app.StartDaemon("temp")
	assert.NoError(t, e)
	n.Write("/tmp/2.txt", []byte("test"))

	time.Sleep(5 * time.Second)

}
