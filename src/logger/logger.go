/**
 * Copyright (C) 2021 Yi Fan Song <yfsong00@gmail.com>
 *
 * This file is part of space-kraken.
 *
 * space-kraken is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * space-kraken is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with space-kraken.  If not, see <https://www.gnu.org/licenses/>.
 **/

package logger

import (
	"fmt"
	"io"
	"time"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

// Logger is the interface of the logger
type Logger interface {
	// Printf writes the message to the output
	//
	// messages are formatted using fmt.Fprintf
	Printf(format string, a ...interface{})

	// Info writes the message to the output
	Info(...interface{})

	// Infof writes teh message to output with format
	Infof(format string, a ...interface{})

	// Error writes the message to the output with error colors
	//
	// messages are formatted using fmt.Fprintln
	Error(...interface{})
}

type Config struct {
	InfoColor  string
	ErrorColor string
	UseColor   bool
}

type logger struct {
	outWriter io.Writer
	errWriter io.Writer

	infoFormat string
	errFormat  string
}

// New creates a new Logger.
//
// writers can be nil, in which case it will use stdout and stderr
func New(outWriter io.Writer, errWriter io.Writer, config Config) Logger {
	var (
		infoFormat = "[%s] %s\n"
		errFormat  = "[%s] ERROR: %s\n"
	)

	if config.UseColor {
		infoFormat = "[%s]" + config.InfoColor + "%s" + Reset + "\n"
		errFormat = "[%s]" + config.ErrorColor + "ERROR: %s" + Reset + "\n"
	}

	return &logger{
		outWriter:  outWriter,
		errWriter:  errWriter,
		infoFormat: infoFormat,
		errFormat:  errFormat,
	}
}

func (l logger) Printf(format string, a ...interface{}) {
	fmt.Fprintf(l.outWriter, l.infoFormat, getTime(), fmt.Sprintf(format, a...))
}

func (l logger) Info(a ...interface{}) {
	fmt.Fprintf(l.outWriter, l.infoFormat, getTime(), fmt.Sprint(a...))
}

func (l logger) Infof(format string, a ...interface{}) {
	fmt.Fprintf(l.outWriter, l.infoFormat, getTime(), fmt.Sprintf(format, a...))
}

func (l logger) Error(a ...interface{}) {
	fmt.Fprintf(l.errWriter, l.errFormat, getTime(), fmt.Sprint(a...))
}

func getTime() string {
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		fmt.Printf("Failed to load location: %s\n", err.Error())
		return ""
	}
	t := time.Now().In(loc)

	return t.Format("Mon Jan _2 2006 15:04:05 MST")
}
