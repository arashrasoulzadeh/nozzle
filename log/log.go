// Package log /*
// Copyright © 2025 Arash Rasoulzadeh <arashrasoulzadeh@gmail.com>
package log

import (
	"Nozzle/src/translation"
	"fmt"
)

func Info(t translation.InfoMessages, fields ...any) {
	//fmt.Println("INFO: "+t, fields)
}
func Error(t translation.InfoMessages, fields ...any) {
	fmt.Println("ERROR: "+t, fields)
}
