/*
Modified by Alessio Greggi © 2025. Based on work by Furkan Pehlivan <furkanpehlivan34@gmail.com>.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level string, format string) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(levelFromString(level)).With().Timestamp().Logger()
	if format == "json" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else if format == "text" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})
	}

	return &logger
}

func levelFromString(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "off":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}
