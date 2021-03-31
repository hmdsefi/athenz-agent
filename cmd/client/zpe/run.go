/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Date: 3/29/21
 * Time: 5:40 PM
 *
 * Description:
 *
 */

package zpe

import (
	"fmt"
	"gitlab.com/trialblaze/athenz-agent/grpc/client"
	"log"
)

func Run() {
	if token == "" || access == "" || resource == "" || port == "" {
		log.Fatal("main: token, access, resource or port is empty, please set value for all of them")
	}

	val, err := client.CheckAccessWithClient(token, access, resource, host, port)
	if err != nil {
		log.Fatalf("CheckAccessWithClient: error when calling `CheckAccessWithToken`, error: %s", err.Error())
	}

	fmt.Printf("Response from server: %d", val)
}
