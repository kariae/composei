package reader

import "strings"

type InputReaderMock struct{
	Content	string
	index	int8
}

func (r *InputReaderMock) ReadString(delimiter byte) (string, error) {
	content := strings.Split(r.Content, string(delimiter))
	input := content[r.index]
	r.index += 1
	return input, nil
}
