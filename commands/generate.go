package commands

import (
	"github.com/urfave/cli"
	"fmt"
)

var GenerateCommand = cli.Command{
	Name: "generate",
	Usage: "Generate docker-compose.yml file",
	Aliases: []string{"g"},
	Action: generate,
}

func generate(c *cli.Context)  {
	composei := `
        +-------+         ____                                     _
        | || || |        / ___|___  _ __ ___  _ __   ___  ___  ___(_)
    +---+---+---+---+   | |   / _ \| '_ ` + "`" + ` _ \| '_ \ / _ \/ __|/ _ \ |
    | || || | || || |   | |__| (_) | | | | | | |_) | (_) \__ \  __/ |
    +-------+-------+    \____\___/|_| |_| |_| .__/ \___/|___/\___|_|
                                             |_|

                          By Zakariae Filali - 0.0.1
                      https://github.com/kariae/composei


`
	fmt.Println(composei)
}
