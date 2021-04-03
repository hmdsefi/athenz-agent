/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/1/21
 * Time: 9:54 PM
 *
 * Description:
 *
 */

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	filePath = "testdata/agent.toml"
)

func TestLoadAgentConfig(t *testing.T) {
	a := assert.New(t)
	config := new(AgentConfiguration)
	err := LoadAgentConfig(config,filePath)
	a.NoError(err)

	a.Equal("9091", config.Properties.Server.Port)
	a.Equal("sidecar-agent", config.Properties.Server.Name)
	a.Equal("testdata/zpu.conf", config.Properties.Config.ZpuConfigFile)
	a.Equal("info", config.Properties.Log.Level)
}

func TestLogProperties_GetMaxAge(t *testing.T) {
	a := assert.New(t)
	config := new(AgentConfiguration)
	err := LoadAgentConfig(config,filePath)
	a.NoError(err)

	a.Equal(time.Duration(24)*time.Hour, config.Properties.Log.GetRotationTime())
	a.Equal(time.Duration(720)*time.Hour, config.Properties.Log.GetMaxAge())
	a.Equal(int64(20971520), config.Properties.Log.GetMaxSize())
}
