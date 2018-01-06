package libs

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"fmt"
	"strings"
)

const FileName = "docker-compose.yml"
const Version  = "3"

type DockerCompose struct {
	file   	string
	values 	yaml.MapSlice
	envVars	[]string
}

func New() *DockerCompose {
	self := &DockerCompose{}
	self.values = yaml.MapSlice{}
	return self
}

func FileExists() bool {
	wd := getWorkingDir()
	if _, err := os.Stat(filepath.FromSlash(fmt.Sprintf("%s/%s", wd, FileName))); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func LoadFile() (*DockerCompose, error) {
	var err error

	dockerCompose := New()

	source, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}

	dockerCompose.file = FileName
	err = yaml.Unmarshal([]byte(source), &dockerCompose.values)
	if err != nil {
		return nil, err
	}

	return dockerCompose, nil
}

func (dockerCompose *DockerCompose) AddService(service yaml.MapItem) {
	dockerCompose.addItemToTopLevel("services", service)
}

func (dockerCompose *DockerCompose) AddNetwork(network yaml.MapItem) {
	dockerCompose.addItemToTopLevel("networks", network)
}

func (dockerCompose *DockerCompose) AddVolume(volume yaml.MapItem) {
	dockerCompose.addItemToTopLevel("volumes", volume)
}

func (dockerCompose *DockerCompose) addItemToTopLevel(topLevelName string, item yaml.MapItem) {
	topLevelIndex := dockerCompose.CreateTopLevel(topLevelName, yaml.MapSlice{})

	// Check if item already exists
	addItem := true
	topLevelItems := dockerCompose.values[topLevelIndex].Value.(yaml.MapSlice)
	if topLevelName == "services" {
		for i, existingItem := range topLevelItems {
			if existingItem.Key == item.Key {
				// Plural -> singular with the trivial way possible xD
				switch ReadLine(fmt.Sprintf("A `%s` already exists with the same name, would you like to replace it", strings.TrimRight(topLevelName, "s")), []string{YesChoice, NoChoice}, false, "") {
				case YesChoice:
					topLevelItems = append(topLevelItems[:i], topLevelItems[i+1:]...)
				case NoChoice:
					addItem = false
				}
			}
		}
	}

	if addItem {
		dockerCompose.values[topLevelIndex].Value = append(topLevelItems, yaml.MapSlice{item}...)
	}
}

func (dockerCompose *DockerCompose) CreateTopLevel(topLevelName string, value interface{}) (index int) {
	exists, index, content := dockerCompose.topLevelExists(topLevelName)
	if !exists {
		content = yaml.MapItem{Key:topLevelName, Value:value}
		dockerCompose.values = append(dockerCompose.values, content)
		index = len(dockerCompose.values) - 1
	}

	return
}

func (dockerCompose *DockerCompose) topLevelExists(key string) (exists bool, index int, content yaml.MapItem) {
	exists = false
	index = 0
	for i,v := range dockerCompose.values {
		if v.Key == key {
			exists = true
			index = i
			content = dockerCompose.values[i]
			break
		}
	}

	return
}

func (dockerCompose *DockerCompose) Save() error {
	var err error

	d, err := yaml.Marshal(&dockerCompose.values)

	//fmt.Println(string(d))
	//os.Exit(1)
	if err == nil {
		err = ioutil.WriteFile(FileName, d, 0644)
	}

	return err
}

func getWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		ERROR(err.Error())
	}

	return wd
}


