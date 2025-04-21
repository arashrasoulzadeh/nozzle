// Package models
//
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package models

type StatusChannelEnum int

const (
	StatusChannelExit        StatusChannelEnum = 0
	StatusChannelFileWritten StatusChannelEnum = 1
)

type StatusChannelFrame struct {
	code StatusChannelEnum
}
