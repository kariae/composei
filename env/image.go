package env

import (
	"strings"
	"fmt"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"github.com/kariae/composei/logger"
)

var envFilePath = filepath.Join("env", "variables")

type getAssets func (filePath string) ([]byte, error)
type AssetsGetter struct {
	GetterFunc	getAssets
}

func GetPossibleEnvVars(assetsGetter AssetsGetter, image string) []string {
	var envVars []string
	var err error

	s := strings.Split(image, ":")
	imageName := s[0]

	envVars, err = assetsGetter.loadEnvVarsFile(fmt.Sprintf("%s.yaml", imageName))

	if err != nil {
		logger.ERROR(err.Error())
	}

	return envVars
}

func (assetsGetter *AssetsGetter) loadEnvVarsFile(fileName string) ([]string, error) {
	var data map[string]interface{}
	var vars []string
	var err error

	source, err := assetsGetter.GetterFunc(filepath.Join(envFilePath, fileName))
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
