/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/24/19
 * Time: 10:30 AM
 *
 * Description:
 *
 */

package downloader

import (
	"github.com/yahoo/athenz/utils/zpe-updater"
	"gitlab.com/trialblaze/athenz-agent/common"
	"gitlab.com/trialblaze/athenz-agent/common/log"
	)

var (
	logger = log.GetLogger(common.GolangFileName())
)

func DownloadPolicies(zpuConfig *zpu.ZpuConfiguration) error {
	err := zpu.PolicyUpdater(zpuConfig)
	if err != nil {
		logger.Error(err.Error())
		return common.Errorf("DownloadPolicies: policy updater failed, %s", err.Error())
	}
	logger.Info("DownloadPolicies: policy updater finished successfully")
	return nil
}
