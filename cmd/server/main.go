/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/19/19
 * Time: 12:45 PM
 *
 * Description:
 * You can run `athenz-agent` project via main method in
 * command line. Before start the project make sure that
 * zpe.conf and athenz.conf files placed in their
 * place.
 *
 */

package main

import (
	"gitlab.com/trialblaze/athenz-agent/cmd/server/zpe"
	"log"
	"os"
)

// main entry point for the athenz-agent server
func main() {
	app := zpe.BuildCLI()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
