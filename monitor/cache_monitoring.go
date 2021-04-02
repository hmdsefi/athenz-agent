/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/12/19
 * Time: 8:37 AM
 *
 * Description:
 * This file has a task to refresh our cached policies and
 * roleTokens.
 *
 */

package monitor

import (
	"fmt"
	"gitlab.com/trialblaze/athenz-agent/cache"
	"gitlab.com/trialblaze/athenz-agent/common"
	"gitlab.com/trialblaze/athenz-agent/config"
	"time"
)

func StartCache(cacheChan chan<- string) {
	for {
		cache.CleanupRoleTokenCache()
		files, err := common.LoadFileStatus(config.ZpeConfig.Properties.PolicyFilesDir)
		if err != nil {
			cacheChan <- fmt.Sprintf("unable to read policy directory, error: %s", err.Error())
		}
		cache.LoadDB(files)
		<-time.After(time.Duration(config.ZpeConfig.Properties.CleanupTokenInterval) * time.Second)
	}
}
