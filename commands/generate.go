package commands

import (
	"github.com/urfave/cli"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/kariae/composei/reader"
	"github.com/kariae/composei/compose"
	"github.com/kariae/composei/logger"
	"bufio"
	"os"
	"github.com/kariae/composei/env"
	"github.com/kariae/composei/version"
)

var GenerateCommand = cli.Command{
	Name: "generate",
	Usage: "Generate docker compose file",
	Aliases: []string{"g"},
	Action: generate,
	Flags: []cli.Flag {
		cli.StringFlag{
			Name:  "compose-file, c",
			Usage: "Specify an alternate Compose file",
		},
		cli.StringFlag{
			Name:  "env-file, e",
			Usage: "Specify an alternate environment file",
		},
		cli.BoolFlag{
			Name: "no-env",
			Usage: "Do not generate environment file",
		},
	},
}

// generate is the main command of Composei that generates docker-compose files
func generate(c *cli.Context)  {
	var err error
	var dockerCompose = compose.New()
	r := bufio.NewReader(os.Stdin)
	PrintComposeiAsciiArt()

	// Set compose and env file names
	if c.String("compose-file") != "" {
		dockerCompose.Filename = c.String("compose-file")
	}
	if c.String("env-file") != "" {
		dockerCompose.EnvFilename = c.String("env-file")
	}

	// Loading attributes
	servicesAttributes := compose.InitServicesAttributes()
	networksAttributes := compose.InitNetworksAttributes()
	volumesAttributes  := compose.InitVolumesAttributes()

	// Check if compose file already exists
	if dockerCompose.FileExists() {
		fmt.Printf("A `%s` file already exists in the current directory\n", dockerCompose.Filename)

		replaceExisting := reader.ReadLine(r, "Would you like to replace or edit it", []string{"r", "e"}, false, "")
		switch replaceExisting {
		case "r":
			//Replacing old compose file
			dockerCompose.CreateTopLevel("version", string(dockerCompose.Version))
		case "e":
			// Loading old source
			err = dockerCompose.LoadFile()
			if err != nil {
				logger.ERROR(err.Error())
			}
		}
	} else {
		dockerCompose.CreateTopLevel("version", string(dockerCompose.Version))
	}

	// Insert attributes
	generateData(r, dockerCompose, servicesAttributes, networksAttributes, volumesAttributes)

	// Saving docker compose file && generate environment file
	err = dockerCompose.Save(r, c.Bool("no-env"))
	if err != nil {
		logger.ERROR(err.Error())
	}

	// Display generation message
	fmt.Println(logger.Green(`Congratulation your docker compose configuration file is created.`))
}

// generateData starts the process of generating top levels attributes
// the attributes data is given by user inputs
func generateData(r reader.InputReader, dockerCompose *compose.DockerCompose, servicesAttrs []compose.Attribute, networksAttrs []compose.Attribute, volumesAttrs []compose.Attribute) {
	//var addService string
	//var addNetwork string
	//var addVolume string

	logger.INFO("Enter '-h' anytime to get a short description of the given attribute\n")

	// Services
	for {
		addService := reader.ReadLine(r, "Add new service", []string{reader.YesChoice, reader.NoChoice}, false, "")
		if addService == reader.NoChoice {
			break
		} else {
			isValid := false
			serviceName, serviceAttributes := generateTopLevelAttributes(r, "service", servicesAttrs)
			// Check that at least the service has image or build attribute
			for _, attr := range serviceAttributes {
				if attr.Key == "build" || attr.Key == "image" {
					isValid = true
					break
				}
			}

		    //	Service should have at least image or build context specified
			if isValid {
				dockerCompose.AddService(r, yaml.MapItem{Key:serviceName, Value:serviceAttributes})
			} else {
				logger.ERROR(fmt.Sprintf("Service %s has neither an image nor a build context specified. At least one must be provided.", serviceName))
			}
		}
	}

	// Networks
	for {
		addNetwork := reader.ReadLine(r, "Add new network", []string{reader.YesChoice, reader.NoChoice}, false, "")
		if addNetwork == reader.NoChoice {
			break
		} else {
			var networkAttributes interface{}
			networkName, attrs := generateTopLevelAttributes(r, "network", networksAttrs)
			if len(attrs) == 0 {
				networkAttributes = nil
			} else {
				networkAttributes = attrs
			}
			dockerCompose.AddNetwork(r, yaml.MapItem{Key:networkName, Value:networkAttributes})
		}
	}

	// Volumes
	for {
		addVolume := reader.ReadLine(r, "Add new volume", []string{reader.YesChoice, reader.NoChoice}, false, "")
		if addVolume == reader.NoChoice {
			break
		} else {
			var volumeAttributes interface{}
			volumeName, attrs := generateTopLevelAttributes(r, "volume", volumesAttrs)
			if len(attrs) == 0 {
				volumeAttributes = nil
			} else {
				volumeAttributes = attrs
			}
			dockerCompose.AddVolume(r, yaml.MapItem{Key: volumeName, Value: volumeAttributes})
		}
	}
}

// generateTopLevelAttributes generates attributes for given top level
func generateTopLevelAttributes(r reader.InputReader, topLevel string, topLevelAttrs []compose.Attribute) (string, yaml.MapSlice) {
	var attributesData yaml.MapSlice
	possibleEntries := map[string][]string{}
	topLevelEntryName := reader.ReadLine(r, fmt.Sprintf("Enter %s name", topLevel), []string{}, false, "")
	assetsGetter := env.AssetsGetter{
		GetterFunc: env.Asset,
	}

	for _, attribute := range topLevelAttrs {
		attributeValues := getAttributeValues(r, attribute, possibleEntries)
		if attributeValues != nil {
			attributesData = append(attributesData, yaml.MapItem{Key: attribute.Name, Value: attributeValues})
		}

		if topLevel == "service" && attribute.Name == "image" && attributeValues != nil {
			if possibleEnvVars := env.GetPossibleEnvVars(assetsGetter, attributeValues.(string)); len(possibleEnvVars) > 0 {
				possibleEntries["environment"] = possibleEnvVars
			}
		}

		// TODO: If volume/network added for a service, generate them automatically as possible entries for 'em.
	}

	return topLevelEntryName, attributesData
}

// getAttributeValues get given attribute's value from user input
func getAttributeValues(r reader.InputReader, attribute compose.Attribute, possibleEntries map[string][]string) interface{} {
	Loop:
		for {
			if attribute.IsList {
				var entry string
				var value []string

				attributeName := attribute.Name

				if possibleEntries[attribute.Name] != nil {
					fmt.Printf("%s:\n", attributeName)
					for _, possibleEntry := range possibleEntries[attribute.Name] {
						entry = reader.ReadLine(r, fmt.Sprintf("  - %s", possibleEntry), []string{}, true, attribute.GetDescription())
						if entry != "" {
							value = append(value, fmt.Sprintf("%s=%s", possibleEntry, entry))
						}
					}

					attributeName = fmt.Sprintf("%s (other)", attributeName)
				}

				for ok := true; ok; ok = entry != "" {
					entry = reader.ReadLine(r, fmt.Sprintf("%s", attributeName), []string{}, true, attribute.GetDescription())
					if entry == "-h" {
						// HELP Message
						logger.INFO(attribute.DisplayHelp())

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
				value := reader.ReadLine(r, fmt.Sprintf("%s", attribute.Name), []string{}, true, attribute.GetDescription())
				if value == "-h" {
					// HELP Message
					logger.INFO(attribute.DisplayHelp())
					continue Loop
				} else if value != "" {
					return value
				}
			}
			break
		}
	return nil
}

// PrintComposeiAsciiArt shows splash Ascii art
func PrintComposeiAsciiArt() {
	composei := `
        +-------+         ____                                     _
        | || || |        / ___|___  _ __ ___  _ __   ___  ___  ___(_)
    +---+---+---+---+   | |   / _ \| '_ ` + "`" + ` _ \| '_ \ / _ \/ __|/ _ \ |
    | || || | || || |   | |__| (_) | | | | | | |_) | (_) \__ \  __/ |
    +-------+-------+    \____\___/|_| |_| |_| .__/ \___/|___/\___|_|
                                             |_|

                          By Zakariae Filalis - %s
                      https://github.com/kariae/composei


`
	fmt.Println(logger.Green(fmt.Sprintf(composei, version.Version.String())))
}
