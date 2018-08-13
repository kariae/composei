package reader

import (
	"fmt"
	"strings"
	"github.com/kariae/composei/logger"
)

const YesChoice  = "y"
const NoChoice   = "N"
const HelpChoice = "-h"

type InputReader interface {
	ReadString(delimiter byte) (string, error)
}

func ReadLine(reader InputReader, message string, possibleInputs []string, acceptEmpty bool, helpMessage string) (value string) {
	mandatory := ""
	possibleInputsString := ""
	messageToDisplay := message

	if !acceptEmpty && len(possibleInputs) == 0 {
		mandatory = logger.Red("*")
		messageToDisplay = fmt.Sprintf("%s (%s)", message, mandatory)
	}

	if len(possibleInputs) > 0 {
		possibleInputsString = logger.Yellow(strings.Join(possibleInputs, "/"))
		messageToDisplay = fmt.Sprintf("%s [%s]", message, possibleInputsString)
	}

	fmt.Printf("%s: ", messageToDisplay)

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

		value = ReadLine(reader, message, possibleInputs, acceptEmpty, helpMessage)
	}

	return
}
