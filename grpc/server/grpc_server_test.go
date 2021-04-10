/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/5/21
 * Time: 2:22 AM
 *
 * Description:
 *
 */

package server

import (
	"context"
	ac "github.com/hamed-yousefi/athenz-agent/.gen/proto/api/command/v1"
	v1 "github.com/hamed-yousefi/athenz-agent/.gen/proto/api/message/v1"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"github.com/hamed-yousefi/athenz-agent/config"
	"github.com/hamed-yousefi/athenz-agent/grpc/server/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	ca        = "testdata/ca-crt.pem"
	serverCrt = "testdata/server-crt.pem"
	serverKey = "testdata/server-key.pem"
	clientCrt = "testdata/client-crt.pem"
	clientKey = "testdata/client-key.pem"
)

var (
	port = randomPort()
)

func TestRunServerWrongPort(t *testing.T) {
	log.NewLogrusInitializer().InitialLog(log.Info)
	a := assert.New(t)
	wg := new(sync.WaitGroup)
	ctx := context.Background()
	err := RunServer(ctx, new(mock.AthenzAgentService), strconv.Itoa(math.MaxInt32), wg)
	a.Error(err)
	a.Equal("listen tcp: address 2147483647: invalid port", err.Error())
}

func TestRunServerInvalidCredential(t *testing.T) {
	log.NewLogrusInitializer().InitialLog(log.Info)
	config.AgentConfig.Properties.Server.CrtPath = "invalidPath"
	config.AgentConfig.Properties.Server.PrivateKeyPath = "invalidPath"
	config.AgentConfig.Properties.Server.CaPath = "invalidPath"

	a := assert.New(t)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx := context.Background()
	err := RunServer(ctx, new(mock.AthenzAgentService), randomPort(), wg)
	a.Error(err)
	a.Equal("open invalidPath: no such file or directory", err.Error())
}

func TestRunServerInvalidCa(t *testing.T) {
	log.NewLogrusInitializer().InitialLog(log.Info)
	config.AgentConfig.Properties.Server.CrtPath = serverCrt
	config.AgentConfig.Properties.Server.PrivateKeyPath = serverKey
	config.AgentConfig.Properties.Server.CaPath = "invalidPath"

	a := assert.New(t)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx := context.Background()
	err := RunServer(ctx, new(mock.AthenzAgentService), randomPort(), wg)
	a.Error(err)
	a.Equal("open invalidPath: no such file or directory", err.Error())
}

func TestRunServer(t *testing.T) {
	log.NewLogrusInitializer().InitialLog(log.Info)
	config.AgentConfig.Properties.Server.CrtPath = ""
	config.AgentConfig.Properties.Server.PrivateKeyPath = ""
	config.AgentConfig.Properties.Server.CaPath = ""
	a := assert.New(t)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := RunServer(ctx, new(mock.AthenzAgentService), port, wg)
		a.NoError(err)
	}()

	<-time.After(2 * time.Second)

	response, err := checkAccessByClientInsecure()
	a.NoError(err)
	a.Equal(v1.AccessStatus_DENY_DOMAIN_EMPTY, response.AccessCheckStatus)

	// cancel the server to shut it down.
	cancel()
}

func TestRunServerWithTLS(t *testing.T) {
	log.NewLogrusInitializer().InitialLog(log.Info)
	config.AgentConfig.Properties.Server.CrtPath = serverCrt
	config.AgentConfig.Properties.Server.PrivateKeyPath = serverKey
	config.AgentConfig.Properties.Server.CaPath = ca

	a := assert.New(t)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx := context.Background()
	go func() {
		err := RunServer(ctx, new(mock.AthenzAgentService), randomPort(), wg)
		a.NoError(err)
	}()

	<-time.After(2 * time.Second)

	wg.Done()
	ctx.Done()
}

func checkAccessByClientInsecure() (*v1.AccessCheckResponse, error) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial("127.0.0.1:"+port, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("unable to connect, error: %s", err.Error())
	}
	defer func() {
		_ = conn.Close()
	}()

	client := ac.NewAthenzAgentClient(conn)

	return client.CheckAccessWithToken(context.Background(),
		&v1.AccessCheckRequest{Token: "token", Access: "access", Resource: "resource"})

}

func randomPort() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(55000) + 10000)
}
