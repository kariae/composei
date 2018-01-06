package libs

import (
	"github.com/fatih/color"
	"fmt"
)

// Level defines all available log levels for log messages.
type Level int

// Log levels.
const (
	Error Level = iota
	Warning
	Info
	Debug
)

var levelNames = []string{
	"ERROR",
	"WARNING",
	"INFO",
	"DEBUG",
}

func ERROR(message string) {
	Log(Error, message)
}

func WARNING(message string) {
	Log(Warning, message)
}

func INFO(message string) {
	Log(Info, message)
}

func DEBUG(message string) {
	Log(Debug, message)
}

func Log(level Level, message string) {
	c := getLevelColor(level).SprintFunc()
	fmt.Println(c(fmt.Sprintf("%s: %s", levelNames[level], message)))
}

func Red(message string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(message)
}

func Yellow(message string) string {
	red := color.New(color.FgYellow).SprintFunc()
	return red(message)
}

func getLevelColor(level Level) *color.Color {
	c := color.New(color.FgBlack)

	switch level {
	case Error:
		c = color.New(color.FgRed)
	case Warning:
		c = color.New(color.FgYellow)
	case Info:
		c = color.New(color.FgBlue)
	case Debug:
		c = color.New(color.FgGreen)
	}

	return c
}
