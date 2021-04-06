/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/24/19
 * Time: 11:00 AM
 *
 * Description:
 *
 */

package monitor

import (
	"fmt"
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"github.com/hamed-yousefi/athenz-agent/config"
	"github.com/hamed-yousefi/athenz-agent/downloader"
	"time"
)

var (
	zpuLogger = log.GetLogger(common.GolangFileName())
)

type (
	// zpuMonitor is an implementation of monitor. It monitors ZPU policy
	// downloader.
	zpuMonitor struct{}
)

// NewZpuMonitor creates new instance Monitor type from zpuMonitor.
func NewZpuMonitor() Monitor {
	return zpuMonitor{}
}

// Start starts a process and monitor it. Most of the time this function
// runs in a separate goroutine, because of that it accept a channel as
// input argument.
func (z zpuMonitor) Start(downloadChan chan<- string) {
	for {
		zpuLogger.Info("Start downloading policy files...")
		err := downloader.NewPolicyDownloader(config.ZpuConfig.Properties).DownloadPolicies()
		if err != nil {
			zpuLogger.Error(err.Error())
			downloadChan <- fmt.Sprintf("Policy updator failed, %s", err.Error())
		}
		<-time.After(time.Duration(config.ZpeConfig.Properties.ZpuDownloadInterval) * time.Second)
	}
}
