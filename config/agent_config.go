/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
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
	"github.com/hamed-yousefi/athenz-agent/common"
	convertor "github.com/xhit/go-str2duration/v2"
	"math/rand"
	"strconv"
	"time"
)

var (
	// AgentConfig is a global variable of AgentConfiguration type. It holds agent's
	// configuration in runtime.
	AgentConfig = newAgentConfiguration()
)

type (
	// AgentConfiguration holds agent's properties. It uses a Loader to load
	// configuration into Properties field.
	AgentConfiguration struct {
		loader Loader
		// Properties holds agent's properties
		Properties *agentProperties
	}

	agentProperties struct {
		Server ServerProperties
		Config Properties
		Log    logProperties
	}

	// Properties is a struct that stores configuration file's paths.
	Properties struct {
		ZpeConfigFile    string `mapstructure:"zpe_config_file"`
		ZpuConfigFile    string `mapstructure:"zpu_config_file"`
		AthenzConfigFile string `mapstructure:"athenz_config_file"`
	}

	// ServerProperties is a struct that represents grpc server information.
	ServerProperties struct {
		Name string
		Port string
		MtlsProperties
	}

	// MtlsProperties is a struct that stores mutual TLS configurations.
	MtlsProperties struct {
		CaPath         string `mapstructure:"ca_path"`
		CrtPath        string `mapstructure:"crt_path"`
		PrivateKeyPath string `mapstructure:"key_path"`
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

// newZpeConfiguration creates a new instance of ZpeConfiguration with
// an empty properties to prevent nil pointer exception.
func newAgentConfiguration() *AgentConfiguration {
	return &AgentConfiguration{Properties: new(agentProperties)}
}

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
		agentConfig.Properties.Server.Port = strconv.Itoa(rand.Intn(55000) + 10000)
	}

	// use default configuration for config loader
	agentConfig.loader.WithDefaultConfig()

	return nil
}

// GetLevel returns log level.
func (p logProperties) GetLevel() string {
	return p.Level
}

// GetPath returns the path that log files must be stored there.
func (p logProperties) GetPath() string {
	return p.Path
}

// GetMaxAge returns the max age of a log file before it gets purged from the file system.
func (p logProperties) GetMaxAge() time.Duration {
	maxAge, err := convertor.ParseDuration(p.MaxAge)
	if err != nil {
		common.Fatalf("invalid input, MaxAge: %s", p.MaxAge)
	}

	return maxAge
}

// GetRotationTime return the time between rotation.
func (p logProperties) GetRotationTime() time.Duration {
	rotationTime, err := convertor.ParseDuration(p.RotationTime)
	if err != nil {
		common.Fatalf("invalid input, RotationTime: %s", p.RotationTime)
	}

	return rotationTime
}

// GetMaxSize returns the log file size between rotation.
func (p logProperties) GetMaxSize() int64 {
	byteCount, err := units.ParseBase2Bytes(p.MaxSize)
	if err != nil {
		common.Fatalf("invalid input, MaxSize: %s", p.MaxSize)
	}

	return int64(byteCount)
}

// GetFilenamePattern returns filename pattern.
func (p logProperties) GetFilenamePattern() string {
	return p.FilenamePattern
}

// IsEmpty checks if MtlsProperties has value or not. If not returns true else
// returns false.
func (p MtlsProperties) IsEmpty() bool {
	if p.CaPath == "" && p.PrivateKeyPath == "" && p.CrtPath == "" {
		return true
	}
	return false
}
