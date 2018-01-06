package attributes

import (
	"fmt"
)

const servicesCategory = 1
const networksCategory = 2
const volumesCategory  = 3

type Attribute struct {
	Name             string
	InputDescription string
	Example          string
	Required         bool
	IsList           bool
	Category         int
	PossibleValues   []string
}

func (attribute Attribute) DisplayHelp() string {
	example := ""
	if attribute.Example != "" {
		example = fmt.Sprintf("(Example: %s)", attribute.Example)
	}

	return fmt.Sprintf("%s %s.", attribute.InputDescription, example)
}

func InitServicesAttributes() []Attribute {
	return []Attribute{
		{
			Name:             "image",
			InputDescription: "The image to start the container from. Can either be a repository/tag or a partial image ID",
			Example:          "nginx:alpine",
			Category:         servicesCategory,
		},
		{
			Name:             "build",
			InputDescription: "String containing a path to the build context",
			Example:          "./app",
			Category:         servicesCategory,
		},
		{
			Name:             "container_name",
			InputDescription: "Custom container Name, rather than a generated default Name",
			Category:         servicesCategory,
		},
		{
			Name:             "ports",
			InputDescription: "Expose ports",
			Example:          "8000:8000",
			IsList:           true,
			Category:         servicesCategory,
		},
		{
			Name:             "volumes",
			InputDescription: "Host paths or named volumes, specified as sub-options to a service",
			IsList:           true,
			Category:         servicesCategory,
		},
		{
			Name:             "depends_on",
			InputDescription: "Express dependency between services",
			Category:         servicesCategory,
		},
		{
			Name:             "entrypoint",
			InputDescription: "Override the default entrypoint",
			Example:          "/app/entrypoint.sh",
			Category:         servicesCategory,
		},
		{
			Name:             "restart",
			InputDescription: "Restart service options",
			PossibleValues:   []string{"no", "always", "on-failure", "unless-stopped"},
			Category:         servicesCategory,
		},
		{
			Name:             "environment",
			InputDescription: "Environment variables",
			Example:          "RACK_ENV=development",
			IsList:           true,
			Category:         servicesCategory,
		},
		{
			Name:             "env_file",
			InputDescription: "Environment variables from a file (Single Value)",
			Example:          "secrets.env",
			IsList:           true,
			Category:         servicesCategory,
		},
		{
			Name:             "networks",
			InputDescription: "Networks to join, referencing entries under the top-level networks key",
			IsList:           true,
			Category:         servicesCategory,
		},
		{
			Name:             "external_links",
			InputDescription: "Link to containers started outside this `docker-compose.yml` or even outside of Compose",
			Example:          "project_db_1:postgresql",
			IsList:           true,
			Category:         servicesCategory,
		},
		{
			Name:             "labels",
			InputDescription: "Metadata to containers using Docker labels",
			IsList:           true,
			Category:         networksCategory,
		},
	}
}

func InitNetworksAttributes() []Attribute {
	return []Attribute{
		{
			Name:             "driver",
			InputDescription: "Driver to be used for this network",
			PossibleValues:   []string{"bridge", "overlay"},
			Category:         networksCategory,
		},
		{
			Name:             "driver_opts",
			InputDescription: "Specify a list of options as key-Value pairs to pass to the driver for this network",
			IsList:           true,
			Category:         networksCategory,
		},
		{
			Name:             "labels",
			InputDescription: "Metadata to containers using Docker labels",
			IsList:           true,
			Category:         networksCategory,
		},
		{
			Name:             "external",
			InputDescription: "If set to true, specifies that this network has been created outside of Compose",
			PossibleValues:   []string{"true", "false"},
			Category:         networksCategory,
		},
	}
}

func (attribute *Attribute) GetDescription() string {
	example := ""
	if attribute.Example != "" {
		example = fmt.Sprintf("(Example: %s)", attribute.Example)
	}

	return fmt.Sprintf("%s %s.", attribute.InputDescription, example)
}

func InitVolumesAttributes() []Attribute {
	return []Attribute{
		{
			Name:             "driver",
			InputDescription: "Driver to be used for this volume",
			Category:         volumesCategory,
		},
		{
			Name:             "driver_opts",
			InputDescription: "Specify a list of options as key-Value pairs to pass to the driver for this volume",
			IsList:           true,
			Category:         volumesCategory,
		},
		{
			Name:             "labels",
			InputDescription: "Metadata to containers using Docker labels",
			IsList:           true,
			Category:         volumesCategory,
		},
		{
			Name:             "external",
			InputDescription: "If set to true, specifies that this volume has been created outside of Compose",
			PossibleValues:   []string{"true", "false"},
			Category:         volumesCategory,
		},
	}
}
