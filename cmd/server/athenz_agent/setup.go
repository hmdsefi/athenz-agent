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

package athenz_agent

import (
	"github.com/urfave/cli"
	"github.com/hamed-yousefi/athenz-agent/config"
)

var (
	athenzConfigPath string
	zpuConfigPath    string
	zpeConfigPath    string
	agentConfPath    string
)

// BuildCLI is the main entry point for the cadence server
func BuildCLI() *cli.App {

	app := cli.NewApp()
	app.Name = "AthenzAgent"
	app.Usage = "AthenzAgent server"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "zpe-config, e",
			Value: config.DefaultZpeConfigPath,
			Usage:  "ZPE utility configuration path",
			EnvVar: config.EnvKeyZpeConfigPath,
			Destination: &zpeConfigPath,
		},
		cli.StringFlag{
			Name:   "athenz-config, a",
			Value:  config.DefaultAthenzConfigPath,
			Usage:  "Athenz configuration file path for ZMS/ZTS urls and public keys",
			EnvVar: config.EnvKeyAthenzConfigPath,
			Destination: &athenzConfigPath,
		},
		cli.StringFlag{
			Name:   "zpu-config, u",
			Value:  config.DefaultZpuConfigPath,
			Usage:  "ZPU utility configuration path",
			EnvVar: config.EnvKeyZpuConfigPath,
			Destination: &zpuConfigPath,
		},
		cli.StringFlag{
			Name:   "agent-config, c",
			Value: config.DefaultAgentConfigPath,
			Usage:  "Agent configuration file path",
			EnvVar: config.EnvKeyAgentConfigPath,
			Destination: &agentConfPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{""},
			Usage:   "start agent server",
			Action: func(c *cli.Context) {
				run()
			},
		},
	}

	return app
}
