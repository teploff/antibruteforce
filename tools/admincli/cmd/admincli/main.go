package main

import (
	"fmt"
	"log"
	"os"

	"github.com/teploff/antibruteforce/tools/admincli/transport/http"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dest",
				Value: "127.0.0.1:8088",
				Usage: "anti brute-force http-service addr",
			},
		},
		Name:  "admin panel",
		Usage: "CLI for administrator use cases to anti brute-force service via http client.",
		Commands: []*cli.Command{
			{
				Name:    "reset_bucket_by_login",
				Aliases: []string{"rbl"},
				Usage:   "reset leaky bucket in the rate limiter for given login",
				Action: func(c *cli.Context) error {
					httpClient := http.NewClient(c.String("dest"))
					if err := httpClient.ResetBucketByLogin(c.Args().First()); err != nil {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println("Successfully!")
					}
					return nil
				},
			},
			{
				Name:    "reset_bucket_by_password",
				Aliases: []string{"rbp"},
				Usage:   "reset leaky bucket in the rate limiter for given password",
				Action: func(c *cli.Context) error {
					httpClient := http.NewClient(c.String("dest"))
					if err := httpClient.ResetBucketByPassword(c.Args().First()); err != nil {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println("Successfully!")
					}
					return nil
				},
			},
			{
				Name:    "reset_bucket_by_ip",
				Aliases: []string{"rbi"},
				Usage:   "reset leaky bucket in the rate limiter for given ip",
				Action: func(c *cli.Context) error {
					httpClient := http.NewClient(c.String("dest"))
					if err := httpClient.ResetBucketByIP(c.Args().First()); err != nil {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println("Successfully!")
					}
					return nil
				},
			},
			{
				Name:    "add_in_whitelist",
				Aliases: []string{"aw"},
				Usage:   "add subnet(ip + mask) in the whitelist. For example: 192.168.130.0/24",
				Action: func(c *cli.Context) error {
					httpClient := http.NewClient(c.String("dest"))
					if err := httpClient.AddInWhitelist(c.Args().First()); err != nil {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println("Successfully!")
					}
					return nil
				},
			},
			{
				Name:    "remove_from_whitelist",
				Aliases: []string{"rw"},
				Usage:   "remove subnet(ip + mask) from the whitelist. For example: 192.168.130.0/24",
				Action: func(c *cli.Context) error {
					httpClient := http.NewClient(c.String("dest"))
					if err := httpClient.RemoveFromWhitelist(c.Args().First()); err != nil {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println("Successfully!")
					}
					return nil
				},
			},
			{
				Name:    "add_in_blacklist",
				Aliases: []string{"ab"},
				Usage:   "add subnet(ip + mask) in the blacklist. For example: 192.168.130.0/24",
				Action: func(c *cli.Context) error {
					httpClient := http.NewClient(c.String("dest"))
					if err := httpClient.AddInBlacklist(c.Args().First()); err != nil {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println("Successfully!")
					}
					return nil
				},
			},
			{
				Name:    "remove_from_blacklist",
				Aliases: []string{"rb"},
				Usage:   "remove subnet(ip + mask) from the blacklist. For example: 192.168.130.0/24",
				Action: func(c *cli.Context) error {
					httpClient := http.NewClient(c.String("dest"))
					if err := httpClient.RemoveFromBlacklist(c.Args().First()); err != nil {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println("Successfully!")
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
