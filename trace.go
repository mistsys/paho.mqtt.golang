/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"log"
	"os"
)

type (
	// Logger interface allows implementations to provide to this package any
	// object that implements the methods defined in it.
	Logger interface {
		Println(v ...interface{})
		Printf(format string, v ...interface{})
	}

	// NOOPLogger implements the logger that does not perform any operation
	// by default. This allows us to efficiently discard the unwanted messages.
	NOOPLogger struct{}
)

func (NOOPLogger) Println(v ...interface{})               {}
func (NOOPLogger) Printf(format string, v ...interface{}) {}

// Internal levels of library output that are initialised to not print
// anything but can be overridden by programmer
var (
	//ERROR    Logger = NOOPLogger{}
	//CRITICAL Logger = NOOPLogger{}
	//WARN     Logger = NOOPLogger{}
	//DEBUG    Logger = NOOPLogger{}
	ERROR    Logger = log.New(os.Stdout, "ERROR: ", 0)
	CRITICAL Logger = log.New(os.Stdout, "CRITICAL: ", 0)
	WARN     Logger = log.New(os.Stdout, "WARN: ", 0)
	DEBUG    Logger = NOOPLogger{}
)
