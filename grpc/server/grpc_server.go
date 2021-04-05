/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
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
	"crypto/tls"
	"crypto/x509"
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"github.com/hamed-yousefi/athenz-agent/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"sync"

	ac "github.com/hamed-yousefi/athenz-agent/.gen/proto/api/command/v1"
)

var (
	logger = log.GetLogger(common.GolangFileName())
)

func RunServer(ctx context.Context, ps ac.AthenzAgentServer, port string, waitGrp *sync.WaitGroup) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	credential, err := mTLSCredential(config.AgentConfig.Properties.Server.MtlsProperties)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// register service
	server := grpc.NewServer(grpc.Creds(credential))
	ac.RegisterAthenzAgentServer(server, ps)

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
	logger.Info("'athenz-agent' gRPC server listening on port: " + port)
	return server.Serve(listen)
}

func mTLSCredential(properties config.MtlsProperties) (credentials.TransportCredentials, error) {

	if properties.IsEmpty() {
		return nil, nil
	}
	certificate, err := tls.LoadX509KeyPair(properties.CrtPath, properties.PrivateKeyPath)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	data, err := ioutil.ReadFile(properties.CaPath)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(data) {
		logger.Error("append ca cert failed!")
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(tlsConfig), err
}
