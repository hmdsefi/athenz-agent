/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
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
	"fmt"
	"github.com/yahoo/athenz/utils/zpe-updater"
	"log"
)

func DownloadPolicies(zpuConfig *zpu.ZpuConfiguration) error {
	err := zpu.PolicyUpdater(zpuConfig)
	if err != nil {
		return fmt.Errorf("DownloadPolicies: policy updator failed, %v", err)
	}
	log.Println("DownloadPolicies: policy updater finished successfully")
	return nil
}
