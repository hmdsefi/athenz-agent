/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Date: 3/28/21
 * Time: 6:16 PM
 *
 * Description:
 *
 */

package zpe

import (
	"context"
	"fmt"
	"gitlab.com/trialblaze/athenz-agent/auth"
	"gitlab.com/trialblaze/athenz-agent/config"
	"gitlab.com/trialblaze/athenz-agent/grpc/server"
	"gitlab.com/trialblaze/athenz-agent/monitor"
	"gitlab.com/trialblaze/athenz-agent/common/util"
	"log"
	"sync"
)

func Run() {

	var cacheError string
	var waitGrp sync.WaitGroup

	loadConfigs()

	// make new directory for metric file, if it doesn't exist
	if err := util.CreateAllDirectories(config.ZpuConfig.Properties.MetricsDir); err != nil {
		log.Fatal("Main> cannot create metrics directory, error: ", err.Error())
	}

	// make new directory for policy files, if it doesn't exist
	if err := util.CreateAllDirectories(config.ZpuConfig.Properties.PolicyFileDir); err != nil {
		log.Fatal("Main> cannot create policy directory, error: ", err.Error())
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

	permissionService := &auth.PermissionService{}

	// start policy downloader
	go monitor.StartDownloader(downloaderChan)
	// start caching policy files into memory
	go monitor.StartCache(cacheChan)

	// start gRPC server in a goroutine
	waitGrp.Add(1)
	go func() {
		if err := server.RunServer(ctx, permissionService, config.ZpeConfig.Properties.GRPCServerPort, &waitGrp); err != nil {
			serverStatusChan <- fmt.Sprintf("Main>startGRPCServer: gRPC server failed to start, error: %s", err.Error())
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
		// do nothing
	}

	// print error reason and exit with code 1
	log.Fatal(cacheError)

}

// loadConfigs loads all configurations
func loadConfigs() {
	if err := config.LoadAgentConfig(config.AgentConfig, agentConfPath); err != nil {
		log.Fatal("Main> unable to read agent config file, error: ", err.Error())
	}

	if err := config.LoadZpeConfig(config.ZpeConfig, zpeConfigPath); err != nil {
		log.Fatal("Main.loadConfigs> unable to open zpe config file, error: ", err.Error())
	}

	if err := config.LoadAthenzConfig(config.KeyStore, athenzConfigPath); err != nil {
		log.Fatal("Main.loadConfigs> unable to open athenz config file, error: ", err.Error())
	}


	if err := config.LoadZpuConfig(config.ZpuConfig, athenzConfigPath, zpuConfigPath); err != nil {
		log.Fatal("Main.loadConfigs> unable to open zpu config file, error: ", err.Error())
	}
}
