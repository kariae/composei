package libs

import (
	"strings"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/kariae/composei/env_vars"
	"path/filepath"
)

const envFilePath = "env_vars"

func GetPossibleEnvVars(image string) []string {
	var envVars []string
	var err error

	s := strings.Split(image, ":")
	imageName := s[0]

	envVars, err = loadEnvVarsFile(fmt.Sprintf("%s.yaml", imageName))
	if err != nil {
		ERROR(err.Error())
	}

	return envVars
}

func loadEnvVarsFile(fileName string) ([]string, error) {
	var data map[string]interface{}
	var vars []string
	var err error

	source, err := env_vars.Asset(filepath.Join(envFilePath, fileName))
	if err != nil {
		return vars, nil
	}

	err = yaml.Unmarshal(source, &data)
	if err != nil {
		return nil, err
	}

	for _, variable := range data["vars"].([]interface{}) {
		vars = append(vars, variable.(string))
	}

	return vars, nil
}
