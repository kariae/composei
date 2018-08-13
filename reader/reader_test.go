package reader

import (
	"testing"
)

func TestReadLine(t *testing.T) {
	testCases := []struct{
		reader			InputReader
		possibleInputs	[]string
		acceptEmpty		bool
		helpMessage		string
		expected		string
	}{
		{
			reader: &InputReaderMock{Content: "test\n"},
			expected: "test",
		},
		{
			reader: &InputReaderMock{Content: "f\ny\n"},
			expected: "y",
			possibleInputs: []string {"y", "N"},
		},
		{
			reader: &InputReaderMock{Content: "\nnot empty\n"},
			expected: "not empty",
			acceptEmpty: false,
		},
	}

	for _, c := range testCases {
		actualInput := ReadLine(c.reader, "dummy message", c.possibleInputs, c.acceptEmpty, c.helpMessage)
		if actualInput != c.expected {
			t.Fatalf("Expected %s but got %s", c.expected, actualInput)
		}

	}
}
