/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Date: 3/28/21
 * Time: 6:16 PM
 *
 * Description:
 *
 */

package athenzagent

import (
	"context"
	"fmt"
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"github.com/hamed-yousefi/athenz-agent/config"
	"github.com/hamed-yousefi/athenz-agent/grpc/api"
	"github.com/hamed-yousefi/athenz-agent/grpc/server"
	"github.com/hamed-yousefi/athenz-agent/monitor"
	"sync"
)

func run() {

	var cacheError string
	var waitGrp sync.WaitGroup

	loadConfigs()
	logInit := log.NewLogrusInitializer()
	logInit.InitialLog(log.GetLevel(config.AgentConfig.Properties.Log.GetLevel())).
		SetupRotation(config.AgentConfig.Properties.Log)

	logger := log.GetLogger(common.GolangFileName())

	// make new directory for metric file, if it doesn't exist
	if err := common.CreateAllDirectories(config.ZpuConfig.Properties.MetricsDir); err != nil {
		logger.Fatalf("cannot create metrics directory, error: %s", err.Error())
	}

	// make new directory for policy files, if it doesn't exist
	if err := common.CreateAllDirectories(config.ZpuConfig.Properties.PolicyFileDir); err != nil {
		logger.Fatalf("cannot create policy directory, error: %s" + err.Error())
	}

	// ZPU channel, it's pipeline for sending error
	cacheChan := make(chan string)
	// gRPC server channel, it's pipeline for sending error
	serverStatusChan := make(chan string)
	// ZPE channel, it's pipeline for sending error
	downloaderChan := make(chan string)
	serverIsShutdown := make(chan string)
	// error message channel, this channel listen to all error channels
	done := make(chan string)

	// gRPC server context
	ctx := context.Background()
	// create context with its cancellation method
	ctx, cancel := context.WithCancel(ctx)

	permissionService := &api.PermissionService{}

	// start policy downloader
	go monitor.NewZpuMonitor().Start(downloaderChan)
	// start caching policy files into memory
	go monitor.NewCacheMonitor().Start(cacheChan)

	// start gRPC server in a goroutine
	waitGrp.Add(1)
	go func() {
		if err := server.RunServer(ctx, permissionService, config.AgentConfig.Properties.Server.Port, &waitGrp); err != nil {
			serverStatusChan <- fmt.Sprintf("%s> gRPC server failed to start, error: %s", common.FuncName(), err.Error())
		}
	}()

	// os.Signal notifier goroutine
	go func() {
		waitGrp.Wait()
		serverIsShutdown <- "Shutdown by OS signal"
	}()

	// this goroutine caches ZPE, ZPU, and gRPC server goroutines
	// errors, so if a goroutine prone an error this function cache
	// that error.
	go func() {
		select {
		case msg := <-cacheChan:
			done <- msg
			return
		case msg := <-serverStatusChan:
			done <- msg
			return
		case msg := <-downloaderChan:
			done <- msg
			return
		}
	}()

	select {
	// wait for cache goroutine
	case cacheError = <-done:
		// something bad happened, shutdown gRPC server
		cancel()
		waitGrp.Wait()
	// wait for shutting down server by OS signals
	case cacheError = <-serverIsShutdown:
		// cancel the server context to prevent context leak
		cancel()
	}

	// print error reason and exit with code 1
	logger.Fatal(cacheError)

}

// loadConfigs loads all configurations
func loadConfigs() {
	// use function name in logs and errors
	funcName := common.FuncName()

	if err := config.LoadGlobalAgentConfig(agentConfPath); err != nil {
		common.Fatalf("%s> unable to read agent config file, error: %s", funcName, err.Error())
	}

	if err := config.LoadGlobalZpeConfig(zpeConfigPath); err != nil {
		common.Fatalf("%s> unable to open zpe config file, error: %s", funcName, err.Error())
	}

	if err := config.LoadGlobalAthenzConfig(athenzConfigPath); err != nil {
		common.Fatalf("%s> unable to open athenz config file, error: %s", funcName, err.Error())
	}

	if err := config.LoadGlobalZpuConfig(athenzConfigPath, zpuConfigPath); err != nil {
		common.Fatalf("%s> unable to open zpu config file, error: %s", funcName, err.Error())
	}
}
