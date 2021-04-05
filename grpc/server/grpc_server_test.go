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
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"github.com/hamed-yousefi/athenz-agent/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

type (
	mockAthenzAgentService struct{}
)

func (m mockAthenzAgentService) CheckAccessWithToken(ctx context.Context, request *v1.AccessCheckRequest) (*v1.AccessCheckResponse, error) {
	return &v1.AccessCheckResponse{AccessCheckStatus: v1.AccessStatus_DENY_DOMAIN_EMPTY}, nil
}

func (m mockAthenzAgentService) GetServiceToken(ctx context.Context, request *v1.ServiceTokenRequest) (*v1.ServiceTokenResponse, error) {
	panic("implement me")
}

func TestRunServer(t *testing.T) {
	log.NewLogrusInitializer().InitialLog(log.Info)
	a := assert.New(t)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx := context.Background()
	go func() {
		err := RunServer(ctx, new(mockAthenzAgentService), port, wg)
		a.NoError(err)
	}()

	<-time.After(2 * time.Second)

	response, err := checkAccessByClientInsecure()
	a.NoError(err)
	a.Equal(v1.AccessStatus_DENY_DOMAIN_EMPTY, response.AccessCheckStatus)
	wg.Done()
	ctx.Done()
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
		err := RunServer(ctx, new(mockAthenzAgentService), randomPort(), wg)
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

func checkAccessByClientSecure() *v1.AccessCheckResponse {
	var conn *grpc.ClientConn

	credential, err := mTLSCredential(config.MtlsProperties{
		CaPath:         ca,
		CrtPath:        clientCrt,
		PrivateKeyPath: clientKey,
	})
	if err != nil {
		common.Fatal(err.Error())
	}

	conn, err = grpc.Dial("127.0.0.1:"+port, grpc.WithTransportCredentials(credential))
	if err != nil {
		logger.Fatalf("unable to connect, error: %s", err.Error())
	}
	defer func() {
		_ = conn.Close()
	}()

	client := ac.NewAthenzAgentClient(conn)

	response, err := client.CheckAccessWithToken(context.Background(),
		&v1.AccessCheckRequest{Token: "token", Access: "access", Resource: "resource"})
	if err != nil {
		common.Fatal(err.Error())
	}

	return response
}

func randomPort() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(55000) + 10000)
}
