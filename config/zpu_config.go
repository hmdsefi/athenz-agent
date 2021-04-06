/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
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
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/yahoo/athenz/utils/zpe-updater"
)

var (
	// ZpuConfig is a global variable of ZpuConfiguration type. It holds ZPU configurations.
	ZpuConfig = new(ZpuConfiguration)
)

type (
	// ZpuConfiguration holds ZPU's properties. It uses a Loader to load
	// configuration into Properties field.
	ZpuConfiguration struct {
		// Yahoo zpu configuration
		Properties *zpu.ZpuConfiguration
	}
)

// LoadGlobalZpuConfig loads config file from input path into the global
// variable ZpuConfig.
func LoadGlobalZpuConfig(athenzConfPath, zpuConfPath string) error {
	return LoadZpuConfig(ZpuConfig, athenzConfPath, zpuConfPath)
}

// LoadZpuConfig reads config file from a specific address and  loads it into
// a ZpuConfiguration object.
func LoadZpuConfig(zpuConfig *ZpuConfiguration, athenzConfPath, zpuConfPath string) error {

	var err error
	// load ZPU configuration using Yahoo zpe-updater
	zpuConfig.Properties, err = zpu.NewZpuConfiguration(".", athenzConfPath, zpuConfPath)
	if err != nil {
		return common.Errorf("unable to get zpu configuration, error: %s", err.Error())
	}

	return nil
}
