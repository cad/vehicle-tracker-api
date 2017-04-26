//go:generate swagger generate spec -o data/swagger.json
//go:generate statik -src data/
package main

import (
	"os"
	"log"
	"fmt"
	"github.com/urfave/cli"
	"github.com/cad/vehicle-tracker-api/server"
	"github.com/cad/vehicle-tracker-api/repository"
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
			Usage:   "Start API server.",
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
		{
			Name:    "createsuperuser",
			Aliases: []string{"c"},
			Usage:   "Create a new user.",
			Action: func(c *cli.Context) {
				println("action:", "createsuperuser")
				configPath := c.String("config-path")
				if err := config.LoadConfigFile(configPath); err != nil {
					fmt.Printf("Error: %s loading configuration file: %s\n", configPath, err)
					os.Exit(1)
				}

				email := c.String("email")
				password := c.String("password")
				repository.ConnectDB(config.C.DB.Type , config.C.DB.URL)
				user, err := repository.CreateNewUser(email, password)
				if err != nil {
					log.Fatal("Can not create user:", err.Error())
					return
				}
				log.Println("User:", user.Email, "created successfuly!")
				defer repository.CloseDB()
			},

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config-path",
					Value: "config.json",
					Usage: "Path to the config file",
				},
				cli.StringFlag{
					Name:  "email",
					Value: "user@example.com",
					Usage: "User's email address",
				},
				cli.StringFlag{
					Name:  "password",
					Value: "1234",
					Usage: "User's password",
				},

			},
		},
	}

	app.Run(os.Args)
}
