package env

import (
	"testing"
	"reflect"
	"errors"
	"gopkg.in/yaml.v2"
)

func TestGetPossibleEnvVars(t *testing.T) {
	mysqlEnvVarsData := []string{
		"MYSQL_ROOT_PASSWORD",
		"MYSQL_DATABASE",
	}
	mysqlEnvVars, _ := yaml.Marshal(map[string]interface{}{
		"vars": mysqlEnvVarsData,
	})

	var phpEnvVars []string

	testCases := []struct{
		image		string
		envVars		[]byte
		err			error
		expected	[]string
	}{
		{
			image: "mysql:latest",
			envVars: mysqlEnvVars,
			expected: mysqlEnvVarsData,
		},
		{
			image: "mysql",
			envVars: mysqlEnvVars,
			expected: mysqlEnvVarsData,
		},
		{
			// Case without environment variables
			image: "php",
			envVars: []byte{},
			err: errors.New("asset env/variables/php.yaml not found"),
			expected: phpEnvVars,
		},
	}

	for _, c := range testCases {
		// AssetsGetter mock
		assetsGetter := AssetsGetter{
			GetterFunc: func(filePath string) ([]byte, error) {
				return c.envVars, c.err
			},
		}

		possibleEnvVars := GetPossibleEnvVars(assetsGetter, c.image)

		if ! reflect.DeepEqual(possibleEnvVars, c.expected) {
			t.Fatalf("For %s image, Expected %s but got %s", c.image, c.expected, possibleEnvVars)
		}
	}
}
