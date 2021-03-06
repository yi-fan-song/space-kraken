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
	"io"
	"io/ioutil"
	"time"

	"github.com/fatih/color"
)

// Logger is the interface of the logger
type Logger interface {
	// Log writes the message to the output
	//
	// messages are formatted using fmt.Fprintln
	Log(a ...interface{}) (n int, err error)

	// Log writes the message to the output
	//
	// messages are formatted using fmt.Fprintf
	Logf(format string, a ...interface{}) (n int, err error)

	// Error writes the message to the output with error colors
	//
	// messages are formatted using fmt.Fprintln
	Error(a ...interface{}) (n int, err error)
}

type logger struct {
	prependTime bool

	writer     io.Writer
	logColor   *color.Color
	errorColor *color.Color
}

// New creates a new Logger.
//
// output can be nil, in which case the logs will be discarded
func New(output io.Writer) Logger {
	l := logger{}

	if output != nil {
		l.writer = output
	} else {
		l.writer = ioutil.Discard
	}

	l.logColor = color.New(color.FgCyan)
	l.errorColor = color.New(color.FgRed)
	l.prependTime = true

	return &l
}

func (l logger) Log(a ...interface{}) (n int, err error) {
	if l.prependTime {
		n, err = prependTime(l.writer, l.logColor)
		if err != nil {
			return
		}
	}

	n2, err := l.logColor.Fprintln(l.writer, a...)
	n += n2

	return
}

func (l logger) Logf(format string, a ...interface{}) (n int, err error) {
	if l.prependTime {
		n, err = prependTime(l.writer, l.logColor)
		if err != nil {
			return
		}
	}

	n2, err := l.logColor.Fprintf(l.writer, format, a...)
	n += n2

	return
}

func (l logger) Error(a ...interface{}) (n int, err error) {
	if l.prependTime {
		n, err = prependTime(l.writer, l.errorColor)
		if err != nil {
			return
		}
	}

	n2, err := l.errorColor.Fprintln(l.writer, a...)
	n += n2

	return
}

func prependTime(w io.Writer, color *color.Color) (int, error) {
	loc, err := time.LoadLocation("EST")
	if err != nil {
		return 0, err
	}
	t := time.Now().In(loc)

	return color.Fprintf(w, "[%s] ", t.Format("Mon Jan _2 2006 15:04:05 MST"))
}
