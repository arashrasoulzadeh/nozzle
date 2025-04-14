// Package models /*
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package models

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type File struct {
	uuid    string
	path    string
	payload []byte
	md5     string
}

func CreateFile(uuid string, path string, payload []byte) File {
	hash := md5.Sum(payload)

	// Convert to hex string
	md5String := hex.EncodeToString(hash[:])

	return File{
		uuid:    uuid,
		path:    path,
		payload: []byte(payload),
		md5:     md5String,
	}
}

func (f *File) MarshalBinary() ([]byte, error) {
	// Example implementation: convert the File struct to a byte slice
	// You can customize this based on the actual fields and requirements
	return []byte(fmt.Sprintf("%s|%s|%s|%s", f.uuid, f.path, base64.StdEncoding.EncodeToString(f.payload), f.md5)), nil
}
func (f *File) UnmarshalBinary(data []byte) error {
	// Split the binary data into fields
	fields := bytes.Split(data, []byte("|"))
	if len(fields) != 4 {
		return fmt.Errorf("invalid binary format")
	}
	// Validate that we have all required fields from the binary data

	// Assign fields back to struct
	f.uuid = string(fields[0])
	f.path = string(fields[1])
	f.payload = fields[2]
	f.md5 = string(fields[3])

	return nil
}
