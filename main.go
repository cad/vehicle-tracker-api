//go:generate swagger generate spec -o data/swagger.json
//go:generate statik -src data/
package main

import (
	"github.com/cad/vehicle-tracker-api/server"
	"fmt"
	"os"
	"flag"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args)>0 {
		switch args[0] {
		case "serve":
			server.ExecuteServer()
		default:
			fmt.Println("Invalid command")
			os.Exit(1)
		}
	} else {
		fmt.Println("Invalid command, please retry with the following format $ main serve")
		os.Exit(1)
	}

}
