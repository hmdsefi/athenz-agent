/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/14/19
 * Time: 4:57 PM
 *
 * Description:
 * This is a gRPC server. Start gRPC server by passing
 * a PermissionServer object and a port number to RunServer
 * method. Use PermissionService object that implements
 * PermissionServer interface and get port from zpe config.
 *
 */

package server

import (
	"context"
	"gitlab.com/trialblaze/athenz-agent/common"
	"gitlab.com/trialblaze/athenz-agent/common/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"

	ac "gitlab.com/trialblaze/grpc-go/pkg/api/common/command/v1"
)

var (
	logger = log.GetLogger(common.GolangFileName())
)

func RunServer(ctx context.Context, ps ac.PermissionServer, port string, waitGrp *sync.WaitGroup) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	ac.RegisterPermissionServer(server, ps)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// wait here until one of this channels receive something
		select {
		case <-c:
			// sig is a ^C, handle it
			logger.Info("shutting down 'athenz-agent' gRPC server...")
			server.GracefulStop()
			waitGrp.Done()
		case <-ctx.Done():
			logger.Info("shutting down 'athenz-agent' gRPC server...")
			server.GracefulStop()
			waitGrp.Done()
		}
	}()

	// start gRPC server
	logger.Info("'athenz-agent' gRPC server listening on port: "+ port)
	return server.Serve(listen)
}
