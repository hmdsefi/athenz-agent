/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Date: 3/28/21
 * Time: 7:43 PM
 *
 * Description:
 *
 */

package zpe

import (
	"github.com/urfave/cli"
)

var (
	token            string
	access           string
	resource         string
	host             string
	port       string
)

// BuildCLI is the main entry point for the cadence server
func BuildCLI() *cli.App {

	app := cli.NewApp()
	app.Name = "AthenzAgent"
	app.Usage = "AthenzAgent server"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "token, t",
			Value:       "",
			Usage:       "RoleToken for client mode",
			Destination: &token,
		},
		cli.StringFlag{
			Name:        "access, a",
			Value:       "",
			Usage:       "The action you want to do on a resource with a specific roleToken",
			Destination: &access,
		},
		cli.StringFlag{
			Name:        "resource, r",
			Value:       "",
			Usage:       "The resource you want to have access",
			Destination: &resource,
		},
		cli.StringFlag{
			Name:        "host, h",
			Value:       "",
			Usage:       "gRPC server address that client wants to connect",
			Destination: &host,
		},
		cli.StringFlag{
			Name:        "port, p",
			Value:       "",
			Usage:       "gRPC server port number that client wants to connect",
			Destination: &port,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "client",
			Aliases: []string{""},
			Usage:   "check an access to a resource by client api",
			Action: func(c *cli.Context) {
				Run()
			},
		},
	}

	return app
}
