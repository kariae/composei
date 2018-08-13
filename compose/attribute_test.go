package compose

import (
	"testing"
	"fmt"
)

func TestInitServicesAttributes(t *testing.T)  {
	attributes := InitServicesAttributes()
	if len(attributes) < 1 {
		t.Fatalf("Expected services attribute length to be > 0")
	}
}

func TestInitNetworksAttributes(t *testing.T)  {
	attributes := InitNetworksAttributes()
	if len(attributes) < 1 {
		t.Fatalf("Expected networks attributes length to be > 0")
	}
}

func TestInitVolumesAttributes(t *testing.T)  {
	attributes := InitVolumesAttributes()
	if len(attributes) < 1 {
		t.Fatalf("Expected volumes attributes length to be > 0")
	}
}

func TestDisplayHelpWhenDescriptionExists(t *testing.T) {
	description := "The image to start the container from. Can either be a repository/tag or a partial image ID"

	attribute := Attribute{
		InputDescription: description,
	}

	helpMessage := attribute.DisplayHelp()
	expectedResult := fmt.Sprintf("%s", attribute.InputDescription)

	if helpMessage != expectedResult {
		t.Fatalf("Expected `%s` but got `%s`", expectedResult, helpMessage)
	}
}

func TestDisplayHelpWhenExampleExists(t *testing.T) {
	description := "The image to start the container from. Can either be a repository/tag or a partial image ID"

	attribute := Attribute{
		InputDescription: description,
		Example:          "nginx:alpine",
	}

	helpMessage := attribute.DisplayHelp()
	expectedResult := fmt.Sprintf("%s (Example: %s)", attribute.InputDescription, attribute.Example)

	if helpMessage != expectedResult {
		t.Fatalf("Expected `%s` but got `%s`", description, helpMessage)
	}
}

func TestDisplayHelpWhenDescriptionDoesNotExists(t *testing.T) {
	attribute := Attribute{}

	helpMessage := attribute.DisplayHelp()
	expectedResult := ""

	if helpMessage != expectedResult {
		t.Fatalf("Expected `` but got `%s`", helpMessage)
	}
}
