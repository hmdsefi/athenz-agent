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
 * Time: 1:24 PM
 *
 * Description:
 *
 */

package config

import (
	"fmt"
	"github.com/yahoo/athenz/utils/zpe-updater"
)

const (
	LogFile = "agent.out"
	LogDir  = "logs"
)

var (
	ZpuConfig = new(ZpuConfiguration)
)

type (
	ZpuConfiguration struct {
		// Yahoo zpu configuration
		Properties *zpu.ZpuConfiguration
	}
)

func LoadZpuConfig(zpuConfig *ZpuConfiguration, athenzConf, zpuConf string) error {

	//if _, err := os.Stat(LogDir); os.IsNotExist(err) {
	//	err := os.Mkdir(LogDir, 0755)
	//	if err != nil {
	//		fmt.Println("Main:createNecessaryFileFolder: cannot create logs directory")
	//		os.Exit(1)
	//	}
	//}


	var err error
	// load ZPU configuration using Yahoo zpe-updater
	zpuConfig.Properties, err = zpu.NewZpuConfiguration(".", athenzConf, zpuConf)
	if err != nil {
		return fmt.Errorf("LoadZpuConfig: unable to get zpu configuration, Error: %v", err)
	}

	return nil
}
