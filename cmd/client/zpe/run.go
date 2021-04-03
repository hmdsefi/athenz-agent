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
	"gitlab.com/trialblaze/athenz-agent/common"
	"gitlab.com/trialblaze/athenz-agent/common/log"
	"gitlab.com/trialblaze/athenz-agent/grpc/client"
)

func Run() {
	log.NewLogrusInitializer().InitialLog(log.Info)
	logger := log.GetLogger(common.GolangFileName())

	if token == "" || access == "" || resource == "" || port == "" {
		logger.Fatal("token, access, resource or port is empty, please set value for all of them")
	}

	val, err := client.CheckAccessWithClient(token, access, resource, host, port)
	if err != nil {
		logger.Fatalf("error when calling agent's client, error: %s", err.Error())
	}

	logger.Info(fmt.Sprintf("resource: %s, access: %s, access_status: %d", resource, access, val))
	fmt.Printf("Response from server: %d\n", val)
}
