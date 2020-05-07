package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dest",
				Value: "127.0.0.1:8088",
				Usage: "anti brute-force addr",
			},
		},
		Name:  "greet",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			fmt.Println(c.String("dest"))
			fmt.Printf("Hello %q", c.Args().Get(0))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "complete a task on the list",
				Action: func(c *cli.Context) error {
					httpClient := NewHTTPClient(c.String("dest"))
					if err := httpClient.AddInWhitelist(c.Args().First()); err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("successfully!")
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	//httpClient := NewHTTPClient(*dest)
}
