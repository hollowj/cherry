// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package repl2

import (
	l "log"
	"os"

	"github.com/abiosoft/ishell/v2"
)

var logger *CliLogger

// Log has log methods
type Log interface {
	Print(...interface{})
	Println(...interface{})
	Printf(string, ...interface{})
}

// CliLogger log by *ishell.Shell
type CliLogger struct {
	shell Log
}

// NewCliLogger creates a clilog object pointer
func NewCliLogger(shell *ishell.Shell) *CliLogger {
	if shell == nil {
		return &CliLogger{shell: l.New(os.Stdout, "", 0)}
	}
	return &CliLogger{shell: shell}
}

// Printf comment
func (l CliLogger) Printf(format string, v ...interface{}) {
	l.shell.Printf(format, v...)
}

// Print comment
func (l CliLogger) Print(v ...interface{}) {
	l.shell.Print(v...)
}

// Println comment
func (l CliLogger) Println(v ...interface{}) {
	err, ok := v[0].(error)
	if ok {
		l.shell.Println(err)
		s, ok := l.shell.(*ishell.Shell)
		if ok {
			l.shell.Print(">>> ")
			s.ShowPrompt(true)
		}
	} else {
		l.shell.Println(v...)
	}
}

// Fatalf comment
func (l CliLogger) Fatalf(format string, v ...interface{}) {
	l.shell.Printf(format, v...)
}

// Fatal comment
func (l CliLogger) Fatal(v ...interface{}) {
	l.shell.Print(v...)
}

// Fatalln comment
func (l CliLogger) Fatalln(v ...interface{}) {
	l.shell.Println(v...)
}

// Panicf comment
func (l CliLogger) Panicf(format string, v ...interface{}) {
	l.shell.Printf(format, v...)
}

// Panic comment
func (l CliLogger) Panic(v ...interface{}) {
	l.shell.Print(v...)
}

// Panicln comment
func (l CliLogger) Panicln(v ...interface{}) {
	l.shell.Println(v...)
}
