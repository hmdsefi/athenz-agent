/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
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
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"

	ac "gitlab.com/trialblaze/grpc-go/pkg/api/common/command/v1"
	msg "gitlab.com/trialblaze/grpc-go/pkg/api/common/message/v1"
)

func CheckAccessWithClient(token, access, resource, host, serverPort string) (int32, error) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(host+serverPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("CheckAccessWithClient: unable to connect, error: %s", err.Error())
	}
	defer func() {
		_ = conn.Close()
	}()

	c := ac.NewPermissionClient(conn)

	response, err := c.CheckAccessWithToken(context.Background(),
		&msg.AccessCheckRequest{Token: token, Access: access, Resource: resource})
	if err != nil {
		return -1, errors.New(fmt.Sprintf("CheckAccessWithClient: error when calling `CheckAccessWithToken`, error: %s", err.Error()))
	}

	return response.AccessCheckStatus, err
}
