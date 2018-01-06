package commands

import (
	"github.com/urfave/cli"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/kariae/composei/libs"
	"github.com/kariae/composei/attributes"
)

var GenerateCommand = cli.Command{
	Name: "generate",
	Usage: "Generate docker-compose.yml file",
	Aliases: []string{"g"},
	Action: generate,
}

var dockerCompose = libs.New()

func generate(c *cli.Context)  {
	var err error
	libs.PrintComposeiAsciiArt()

	// Loading attributes
	servicesAttributes := attributes.InitServicesAttributes()
	networksAttributes := attributes.InitNetworksAttributes()
	volumesAttributes  := attributes.InitVolumesAttributes()

	// Check if `docker-compose.yml` already exists
	if libs.FileExists() {
		fmt.Println("A `docker-compose.yml` file already exists in the current directory")

		replaceExisting := libs.ReadLine("Would you like to replace or edit it", []string{"r", "e"}, false, "")
		switch replaceExisting {
		case "r":
			//Replacing old docker-compose.yml file
			dockerCompose.CreateTopLevel("version", string(libs.Version))
		case "e":
			// Loading old source
			dockerCompose, err = libs.LoadFile()
			if err != nil {
				libs.ERROR(err.Error())
			}
		}
	} else {
		dockerCompose.CreateTopLevel("version", string(libs.Version))
	}

	// Insert attributes
	generateData(servicesAttributes, networksAttributes, volumesAttributes)

	// Saving docker-compose.yml file
	err = dockerCompose.Save()
	if err != nil {
		libs.ERROR(err.Error())
	}
}

func generateData(servicesAttrs []attributes.Attribute, networksAttrs []attributes.Attribute, volumesAttrs []attributes.Attribute) {
	var addService string
	var addNetwork string
	var addVolume string

	libs.INFO("Enter '-h' anytime to get a short description of the given attribute\n")

	// Services
	for {
		addService = libs.ReadLine("Add new service", []string{libs.YesChoice, libs.NoChoice}, false, "")
		if addService == libs.NoChoice {
			break
		} else {
			isValid := false
			serviceName, serviceAttributes := generateTopLevelAttributes("service", servicesAttrs)
			// Check that at least the service has image or build attribute
			for _, attr := range serviceAttributes {
				if attr.Key == "build" || attr.Key == "image" {
					isValid = true
				}
			}

		    //	Service should have at least image or build context specified
			if isValid {
				dockerCompose.AddService(yaml.MapItem{Key:serviceName, Value:serviceAttributes})
			} else {
				libs.ERROR(fmt.Sprintf("Service %s has neither an image nor a build context specified. At least one must be provided.", serviceName))
			}
		}
	}

	// Networks
	for {
		addNetwork = libs.ReadLine("Add new network", []string{libs.YesChoice, libs.NoChoice}, false, "")
		if addNetwork == libs.NoChoice {
			break
		} else {
			var networkAttributes interface{}
			networkName, attrs := generateTopLevelAttributes("network", networksAttrs)
			if len(attrs) == 0 {
				networkAttributes = nil
			} else {
				networkAttributes = attrs
			}
			dockerCompose.AddNetwork(yaml.MapItem{Key:networkName, Value:networkAttributes})
		}
	}

	// Volumes
	for {
		addVolume = libs.ReadLine("Add new volume", []string{libs.YesChoice, libs.NoChoice}, false, "")
		if addVolume == libs.NoChoice {
			break
		} else {
			var volumeAttributes interface{}
			volumeName, attrs := generateTopLevelAttributes("volume", volumesAttrs)
			if len(attrs) == 0 {
				volumeAttributes = nil
			} else {
				volumeAttributes = attrs
			}
			dockerCompose.AddVolume(yaml.MapItem{Key: volumeName, Value: volumeAttributes})
		}
	}
}

func generateTopLevelAttributes(topLevel string, topLevelAttrs []attributes.Attribute) (string, yaml.MapSlice) {
	var attributesData yaml.MapSlice
	possibleEntries := map[string][]string{}
	topLevelEntryName := libs.ReadLine(fmt.Sprintf("Enter %s name", topLevel), []string{}, false, "")

	for _, attribute := range topLevelAttrs {
		attributeValues := getAttributeValues(attribute, possibleEntries)
		if attributeValues != nil {
			attributesData = append(attributesData, yaml.MapItem{Key: attribute.Name, Value: attributeValues})
		}

		if topLevel == "service" && attribute.Name == "image" && attributeValues != nil {
			if possibleEnvVars := libs.GetPossibleEnvVars(attributeValues.(string)); len(possibleEnvVars) > 0 {
				possibleEntries["environment"] = possibleEnvVars
			}
		}
	}

	return topLevelEntryName, attributesData
}
func getAttributeValues(attribute attributes.Attribute, possibleEntries map[string][]string) interface{} {
	Loop:
		for {
			if attribute.IsList {
				var entry string
				var value []string

				attributeName := attribute.Name

				if possibleEntries[attribute.Name] != nil {
					fmt.Println(fmt.Sprintf("%s:", attributeName))
					for _, possibleEntry := range possibleEntries[attribute.Name] {
						entry = libs.ReadLine(fmt.Sprintf("  - %s", possibleEntry), []string{}, true, attribute.GetDescription())
						if entry != "" {
							value = append(value, fmt.Sprintf("%s=%s", possibleEntry, entry))
						}
					}

					attributeName = fmt.Sprintf("%s (other)", attributeName)
				}

				for ok := true; ok; ok = entry != "" {
					entry = libs.ReadLine(fmt.Sprintf("%s", attributeName), []string{}, true, attribute.GetDescription())
					if entry == "-h" {
						// HELP Message
						libs.INFO(attribute.DisplayHelp())

					} else if entry != "" {
						value = append(value, entry)

					} else {
						break
					}
				}
				if len(value) > 0 {
					return value
				}
			} else {
				value := libs.ReadLine(fmt.Sprintf("%s", attribute.Name), []string{}, true, attribute.GetDescription())
				if value == "-h" {
					// HELP Message
					libs.INFO(attribute.DisplayHelp())
					continue Loop
				} else if value != "" {
					return value
				}
			}
			break
		}
	return nil
}
