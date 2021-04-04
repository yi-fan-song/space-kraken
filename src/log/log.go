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

// Package log provides a logger client.
package log

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

// DefaultConfig holds default values
var DefaultConfig = Config{
	InfoColor:  Cyan,
	ErrorColor: Red,
	UseColor:   true,
}

// Config holds config info for the logger.
type Config struct {
	InfoColor  string
	ErrorColor string
	UseColor   bool
}

// Logger is the struct holding information on log location and format.
type Logger struct {
	outWriter io.Writer
	errWriter io.Writer

	printFormat string
	infoFormat  string
	errFormat   string
}

// New creates a new Logger.
func New(outWriter io.Writer, errWriter io.Writer, config Config) Logger {
	var (
		printFormat = "[%s] %s\n"
		infoFormat  = "[%s] %s\n"
		errFormat   = "[%s] ERROR: %s\n"
	)

	if config.UseColor {
		infoFormat = "[%s]" + config.InfoColor + "%s" + Reset + "\n"
		errFormat = "[%s]" + config.ErrorColor + "ERROR: %s" + Reset + "\n"
	}

	return Logger{
		outWriter:   outWriter,
		errWriter:   errWriter,
		printFormat: printFormat,
		infoFormat:  infoFormat,
		errFormat:   errFormat,
	}
}

// Info logs info level to output.
func (l Logger) Info(a ...interface{}) {
	fmt.Fprintf(l.outWriter, l.infoFormat, getTime(), fmt.Sprint(a...))
}

// Infof logs info level to output with formatting.
func (l Logger) Infof(format string, a ...interface{}) {
	fmt.Fprintf(l.outWriter, l.infoFormat, getTime(), fmt.Sprintf(format, a...))
}

// Error logs error level to error output.
func (l Logger) Error(a ...interface{}) {
	fmt.Fprintf(l.errWriter, l.errFormat, getTime(), fmt.Sprint(a...))
}

func getTime() string {
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		return ""
	}
	t := time.Now().In(loc)

	return t.Format("Mon Jan _2 2006 15:04:05 MST")
}
