/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Date: 3/29/21
 * Time: 5:39 PM
 *
 * Description:
 *
 */

package main

import (
	"gitlab.com/trialblaze/athenz-agent/cmd/client/zpe"
	"log"
	"os"
)

// main entry point for the athenz-agent client
func main() {
	app := zpe.BuildCLI()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
