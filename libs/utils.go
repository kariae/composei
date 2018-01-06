package libs

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"github.com/fatih/color"
)

const YesChoice  = "y"
const NoChoice   = "N"
const HelpChoice = "-h"

var reader = bufio.NewReader(os.Stdin)

func ReadLine(message string, possibleInputs []string, acceptEmpty bool, helpMessage string) (value string) {
	mandatory := ""
	possibleInputsString := ""
	messageToDisplay := message

	if !acceptEmpty && len(possibleInputs) == 0 {
		mandatory = Red("*")
		messageToDisplay = fmt.Sprintf("%s (%s)", message, mandatory)
	}

	if len(possibleInputs) > 0 {
		possibleInputsString = Yellow(strings.Join(possibleInputs, "/"))
		messageToDisplay = fmt.Sprintf("%s [%s]", message, possibleInputsString)
	}

	fmt.Print(fmt.Sprintf("%s: ", messageToDisplay))

	input, _ := reader.ReadString('\n')
	input = strings.Trim(input," \r\n")

	if len(possibleInputs) == 0 { // Accept any input
		value = input
	} else {
		for _, possibleInput := range append(possibleInputs, HelpChoice) {
			if input == possibleInput {
				value = input
			}
		}
	}

	if (input == HelpChoice || value == "") && !acceptEmpty {
		if input == HelpChoice && helpMessage != "" {
			fmt.Println(helpMessage)
		}

		value = ReadLine(message, possibleInputs, acceptEmpty, helpMessage)
	}

	return
}

func PrintComposeiAsciiArt() {
	composei := `
        +-------+         ____                                     _
        | || || |        / ___|___  _ __ ___  _ __   ___  ___  ___(_)
    +---+---+---+---+   | |   / _ \| '_ ` + "`" + ` _ \| '_ \ / _ \/ __|/ _ \ |
    | || || | || || |   | |__| (_) | | | | | | |_) | (_) \__ \  __/ |
    +-------+-------+    \____\___/|_| |_| |_| .__/ \___/|___/\___|_|
                                             |_|

                          By Zakariae Filali - 0.0.1
                      https://github.com/kariae/composei


`
	color.Green(composei)
}
