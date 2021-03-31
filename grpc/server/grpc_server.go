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
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	ac "gitlab.com/trialblaze/grpc-go/pkg/api/common/command/v1"
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
			log.Println("shutting down 'athenz-agent' gRPC server...")
			server.GracefulStop()
			waitGrp.Done()
		case <-ctx.Done():
			log.Println("shutting down 'athenz-agent' gRPC server...")
			server.GracefulStop()
			waitGrp.Done()
		}
	}()

	// start gRPC server
	log.Printf("'athenz-agent' gRPC server listenting on port: %s", port)
	return server.Serve(listen)
}
