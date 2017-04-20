//go:generate swagger generate spec -o data/swagger.json
//go:generate statik -src data/
package main

import (
	"os"
	"github.com/urfave/cli"
	"github.com/cad/vehicle-tracker-api/server"
	"github.com/cad/vehicle-tracker-api/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "vehicle-tracker-api"
	app.Usage = "Vehicle Tracker API Server"
	app.Version = config.VERSION
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "",
			Action: func(c *cli.Context) {
				println("action:", "run")
				configPath := c.String("config-path")

				// Run the App here
				server.ExecuteServer(configPath)
			},

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config-path",
					Value: "config.json",
					Usage: "Path to the config file",
				},
			},
		},
	}

	app.Run(os.Args)
}
