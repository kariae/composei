package compose

import (
	"testing"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"github.com/kariae/composei/reader"
)

func TestFileExistsWhenThereIsNoFile(t *testing.T)  {
	dockerCompose := New()
	// remove file if it exists
	removeDockerComposeFile(dockerCompose.Filename)

	if dockerCompose.FileExists() {
		t.Fatalf("Expected file `%s` to not exists.", dockerCompose.Filename)
	}
}

func TestFileExistsWhenThereIsAFile(t *testing.T)  {
	dockerCompose := New()

	// create tmp docker-compose file
	file, err:= createDockerComposeFile(dockerCompose.Filename)
	defer os.Remove(file)
	if err != nil {
		t.Error(err)
	}

	if !dockerCompose.FileExists() {
		t.Fatalf("Expected file `%s` to exists.", dockerCompose.Filename)
	}
}

func TestLoadFile(t *testing.T) {
	dockerCompose := New()
	file, err := createDockerComposeFile(dockerCompose.Filename)
	defer os.Remove(file)
	if err != nil {
		t.Error(err)
	}
	dockerCompose.LoadFile()

	// Test loading version
	versionTitle, versionNumber := "version", "3"
	if dockerCompose.Values[0].Key != versionTitle {
		t.Fatalf("Expected %s but got %s", versionTitle, dockerCompose.Values[0].Key)
	}
	if dockerCompose.Values[0].Value != versionNumber {
		t.Fatalf("Expected %s but got %s", versionNumber, dockerCompose.Values[0].Value)
	}

	// Test loading top level attribute (services)
	servicesTitle, composeiService := "services", yaml.MapSlice{
		yaml.MapItem{
			Key: "composei",
			Value: yaml.MapSlice{
				yaml.MapItem{
					Key: "image",
					Value: "golang:alpine",
				},
			},
		},
	}
	if dockerCompose.Values[1].Key != servicesTitle {
		t.Fatalf("Expected %s but got %s", servicesTitle, dockerCompose.Values[0].Key)
	}

	actualValue, _ := yaml.Marshal(&dockerCompose.Values[1].Value)
	expectedValue, _ := yaml.Marshal(composeiService)

	if string(actualValue) != string(expectedValue) {
		t.Fatalf("Expected %s but got %s", composeiService, dockerCompose.Values[1].Value)
	}
}

func TestTopLevelExists(t *testing.T) {
	dockerCompose := New()

	file, err := createDockerComposeFile(dockerCompose.Filename)
	defer os.Remove(file)
	if err != nil {
		t.Error(err)
	}
	dockerCompose.LoadFile()

	exists, _, _ := dockerCompose.topLevelExists("services")
	if !exists {
		t.Fatalf("Expected `services` top level to be existing.")
	}
}

func TestTopLevelExistsWhenIdDoesNot(t *testing.T) {
	dockerCompose := New()

	file, err := createDockerComposeFile(dockerCompose.Filename)
	defer os.Remove(file)
	if err != nil {
		t.Error(err)
	}
	dockerCompose.LoadFile()

	exists, _, _ := dockerCompose.topLevelExists("networks")
	if exists {
		t.Fatalf("Expected `networks` top level to be existing.")
	}
}

func TestCreateTopLevelWhenNew(t *testing.T) {
	dockerCompose := New()

	topLevelName, topLevelValue := "version", "3"

	dockerCompose.CreateTopLevel(topLevelName, topLevelValue)

	if dockerCompose.Values[0].Key != topLevelName {
		t.Fatalf("Expected %s but got %s", topLevelName, dockerCompose.Values[0].Key)
	}

	if dockerCompose.Values[0].Value != topLevelValue {
		t.Fatalf("Expected %s but got %s", topLevelValue, dockerCompose.Values[0].Value)
	}
}

func TestCreateTopLevelWhenExisting(t *testing.T) {
	dockerCompose := New()

	topLevelName, topLevelValue := "version", "3"

	dockerCompose.CreateTopLevel(topLevelName, topLevelValue)
	// Try to create the top level a second time
	dockerCompose.CreateTopLevel(topLevelName, "2")

	if dockerCompose.Values[0].Key != topLevelName {
		t.Fatalf("Expected %s but got %s", topLevelName, dockerCompose.Values[0].Key)
	}

	if dockerCompose.Values[0].Value != topLevelValue {
		t.Fatalf("Expected %s but got %s", topLevelValue, dockerCompose.Values[0].Value)
	}
}

func TestAddServiceToTopLevel(t *testing.T) {
	dockerCompose := New()

	composeiService := yaml.MapItem{
		Key: "composei",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key: "image",
				Value: "golang:alpine",
			},
		},
	}
	
	r := reader.InputReaderMock{}

	dockerCompose.addItemToTopLevel(&r, "services", composeiService)

	actualValue, _ := yaml.Marshal(&dockerCompose.Values[0].Value)
	expectedValue, _ := yaml.Marshal(yaml.MapSlice{composeiService})

	if string(actualValue) != string(expectedValue) {
		t.Fatalf("Expected %s but got %s", composeiService, dockerCompose.Values[1].Value)
	}
}

func TestAddServiceWithExistingNameReplaceIt(t *testing.T) {
	dockerCompose := New()

	composeiService := yaml.MapItem{
		Key: "composei",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key: "image",
				Value: "golang:alpine",
			},
		},
	}

	composeiServiceRC := yaml.MapItem{
		Key: "composei",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key: "image",
				Value: "golang:rc-alpine",
			},
		},
	}

	r := reader.InputReaderMock{Content:"y\n"}
	dockerCompose.addItemToTopLevel(&r, "services", composeiService)

	// Add a service with the same name
	dockerCompose.addItemToTopLevel(&r, "services", composeiServiceRC)


	actualValue, _ := yaml.Marshal(&dockerCompose.Values[0].Value)
	expectedValue, _ := yaml.Marshal(yaml.MapSlice{composeiServiceRC})


	if string(actualValue) != string(expectedValue) {
		t.Fatalf("Expected %s but got %s", expectedValue, actualValue)
	}
}

func TestAddServiceWithExistingNameWithoutReplacingIt(t *testing.T) {
	dockerCompose := New()

	composeiService := yaml.MapItem{
		Key: "composei",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key: "image",
				Value: "golang:alpine",
			},
		},
	}

	composeiServiceRC := yaml.MapItem{
		Key: "composei",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key: "image",
				Value: "golang:rc-alpine",
			},
		},
	}

	r := reader.InputReaderMock{Content:"N\n"}
	dockerCompose.addItemToTopLevel(&r, "services", composeiService)

	// Add a service with the same name
	dockerCompose.addItemToTopLevel(&r, "services", composeiServiceRC)

	actualValue, _ := yaml.Marshal(&dockerCompose.Values[0].Value)
	expectedValue, _ := yaml.Marshal(yaml.MapSlice{composeiService})


	if string(actualValue) != string(expectedValue) {
		t.Fatalf("Expected %s but got %s", expectedValue, actualValue)
	}
}

func TestSave(t *testing.T) {
	dockerCompose := New()
	r := reader.InputReaderMock{}

	dockerCompose.Save(&r, true)
	defer os.Remove(dockerCompose.Filename)

	if _, err := os.Stat(filepath.FromSlash(dockerCompose.Filename)); os.IsNotExist(err) {
		t.Fatalf("%s expected to be created.", dockerCompose.Filename)
	}
}

func TestSaveWithEnvFile(t *testing.T) {
	dockerCompose := New()
	r := reader.InputReaderMock{Content:"alpine"}

	dockerCompose.Values = yaml.MapSlice{
		yaml.MapItem{
			Key: "services",
			Value: yaml.MapSlice{
				yaml.MapItem{
					Key: "composei",
					Value: yaml.MapSlice{
						yaml.MapItem{
							Key: "image",
							Value: "golang:${DIST}",
						},
					},
				},
			},
		},
	}

	dockerCompose.Save(&r, false)
	defer os.Remove(dockerCompose.Filename)
	defer os.Remove(dockerCompose.EnvFilename)

	if _, err := os.Stat(filepath.FromSlash(dockerCompose.Filename)); os.IsNotExist(err) {
		t.Fatalf("%s expected to be created.", dockerCompose.Filename)
	}

	if _, err := os.Stat(filepath.FromSlash(dockerCompose.EnvFilename)); os.IsNotExist(err) {
		t.Fatalf("%s expected to be created.", dockerCompose.EnvFilename)
	}
}

func TestGenerateEnvFile(t *testing.T) {
	dockerCompose := New()
	dockerComposeContent := "version: \"3\"\nservices:\n  composei:\n    image: golang:${DIST}\n    container_name: ${PREFIX}.app"

	// Generate .env file
	r := reader.InputReaderMock{Content: "alpine\ncomposei"}
	dockerCompose.generateEnvFile(&r, dockerComposeContent)
	defer os.Remove(dockerCompose.EnvFilename)

	source, err := ioutil.ReadFile(dockerCompose.EnvFilename)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(source), "DIST=alpine") {
		t.Fatalf("Expected %s in .env file", "DIST=alpine")
	}

	if !strings.Contains(string(source), "PREFIX=composei") {
		t.Fatalf("Expected %s in .env file", "PREFIX=composei")
	}
}

func removeDockerComposeFile(filename string) {
	wd := getWorkingDir()
	os.Remove(filepath.FromSlash(filepath.Join(wd, filename)))
}

func createDockerComposeFile(filename string) (string, error) {
	wd := getWorkingDir()
	configFilePath := filepath.FromSlash(filepath.Join(wd, filename))

	var file, err = os.Create(configFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString("version: \"3\"\nservices:\n  composei:\n    image: golang:alpine")
	if err != nil {
		return "", err
	}

	return configFilePath, nil
}
