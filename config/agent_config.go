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
	"github.com/alecthomas/units"
	convertor "github.com/xhit/go-str2duration/v2"
	"gitlab.com/trialblaze/athenz-agent/common"
	"math/rand"
	"time"
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
		Server ServerProperties
		Config Properties
		Log    logProperties
	}

	Properties struct {
		ZpeConfigFile    string `mapstructure:"zpe_config_file"`
		ZpuConfigFile    string `mapstructure:"zpu_config_file"`
		AthenzConfigFile string `mapstructure:"athenz_config_file"`
	}

	ServerProperties struct {
		Name string
		Port string
	}

	// logProperties represents log config
	logProperties struct {
		Level           string
		Path            string
		MaxAge          string `mapstructure:"max_age"`
		MaxSize         string `mapstructure:"max_size"`
		FilenamePattern string `mapstructure:"filename_pattern"`
		RotationTime    string `mapstructure:"rotation_time"`
	}
)

// LoadGlobalAgentConfig loads config file from input path into the global
// variable AgentConfig.
func LoadGlobalAgentConfig(filePath string) error {
	return LoadAgentConfig(AgentConfig, filePath)
}

// LoadAgentConfig reads config file from a specific address and loads it into
//a AgentConfiguration object
func LoadAgentConfig(agentConfig *AgentConfiguration, filePath string) error {

	// load config properties into agentProperties
	agentConfig.Properties = new(agentProperties)
	agentConfig.loader = NewConfigLoader()
	if err := agentConfig.loader.LoadConfig(agentConfig.Properties, filePath); err != nil {
		return common.Errorf("unable to load config from %s : %s", filePath, err.Error())
	}

	if agentConfig.Properties.Server.Port == "" {
		rand.Seed(time.Now().UnixNano())
		agentConfig.Properties.Server.Port = string(rune(rand.Intn(55000) + 10000))
	}

	// use default configuration for config loader
	agentConfig.loader.WithDefaultConfig()

	return nil
}

func (p logProperties) GetLevel() string {
	return p.Level
}

func (p logProperties) GetPath() string {
	return p.Path
}

func (p logProperties) GetMaxAge() time.Duration {
	maxAge, err := convertor.ParseDuration(p.MaxAge)
	if err != nil {
		common.Fatalf("invalid input, MaxAge: %s", p.MaxAge)
	}

	return maxAge
}

func (p logProperties) GetRotationTime() time.Duration {
	rotationTime, err := convertor.ParseDuration(p.RotationTime)
	if err != nil {
		common.Fatalf("invalid input, RotationTime: %s", p.RotationTime)
	}

	return rotationTime
}

func (p logProperties) GetMaxSize() int64 {
	byteCount, err := units.ParseBase2Bytes(p.MaxSize)
	if err != nil {
		common.Fatalf("invalid input, MaxSize: %s", p.MaxSize)
	}

	return int64(byteCount)
}

func (p logProperties) GetFilenamePattern() string {
	return p.FilenamePattern
}
