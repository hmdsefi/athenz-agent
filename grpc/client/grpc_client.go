/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 3/4/19
 * Time: 1:38 PM
 *
 * Description:
 * this is simple gRPC client for testing our server `CheckAccessWithToken`
 * method.
 *
 */

package client

import (
	"errors"
	"fmt"
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	ac "github.com/hamed-yousefi/athenz-agent/.gen/proto/api/command/v1"
	msg "github.com/hamed-yousefi/athenz-agent/.gen/proto/api/message/v1"
)

var (
	logger = log.GetLogger(common.GolangFileName())
)

// CheckAccessWithClient connects to a athenz-agent server to check an access to a
// resource for a token.
//
// CheckAccessWithClient is useful to test tokens and other kind of tests.
func CheckAccessWithClient(token, access, resource, host, serverPort string) (int32, error) {

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(host+serverPort, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("CheckAccessWithClient: unable to connect, error: %s", err.Error())
	}
	defer func() {
		_ = conn.Close()
	}()

	client := ac.NewAthenzAgentClient(conn)

	response, err := client.CheckAccessWithToken(context.Background(),
		&msg.AccessCheckRequest{Token: token, Access: access, Resource: resource})
	if err != nil {
		return -1, errors.New(fmt.Sprintf("%s> error when calling `CheckAccessWithToken`, error: %s",
			common.FuncName(), err.Error()))
	}

	return int32(response.AccessCheckStatus.Enum().Number()), err
}
