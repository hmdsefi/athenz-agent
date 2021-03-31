/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 3/5/19
 * Time: 10:03 AM
 *
 * Description:
 * Set some of global properties in agent.conf
 *
 */

package config

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	AgentConfig = new(AgentConfiguration)
)

type (

	AgentConfiguration struct {
		loader     Loader
		Properties *agentProperties
	}

	agentProperties struct {
		RootDir          string // project root directory
		ConfigDir        string // config directory from RootDir
		ZpeConfigFile    string
		ZpuConfigFile    string
		AthenzConfigFile string
	}
)

// LoadAgentConfig reads config file from a specific address and
// loads it into a AgentConfiguration object
func LoadAgentConfig(agentConfig *AgentConfiguration, filePath string) error {

	// load config properties into agentProperties
	agentConfig.Properties = new(agentProperties)
	agentConfig.loader = NewConfigLoader()
	if err := agentConfig.loader.LoadConfig(agentConfig.Properties, filePath); err != nil {
		return errors.New(fmt.Sprintf("LoadAgentConfig: unable to load config from %s : %s", filePath, err.Error()))
	}

	// use default configuration for config loader
	agentConfig.loader.WithDefaultConfig()

	return nil
}
