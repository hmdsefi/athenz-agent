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
	"github.com/hamed-yousefi/athenz-agent/config"
	"github.com/hamed-yousefi/athenz-agent/downloader"
	"time"
)

func StartDownloader(downloadChan chan<- string) {
	for {
		err := downloader.DownloadPolicies(config.ZpuConfig.Properties)
		if err != nil {
			downloadChan <- fmt.Sprintf("Policy updator failed, %s", err.Error())
		}
		<-time.After(time.Duration(config.ZpeConfig.Properties.ZpuDownloadInterval) * time.Second)
	}
}
