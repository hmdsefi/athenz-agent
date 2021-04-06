/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Date: 3/28/21
 * Time: 6:19 PM
 *
 * Description:
 *
 */

package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
)

const (
	// EnvKeyZpeConfigPath the environment variable key for ZPE config path
	EnvKeyZpeConfigPath = "ZPE_CONFIG_PATH"
	// EnvKeyAthenzConfigPath the environment variable key for Athenz server config path
	EnvKeyAthenzConfigPath = "ATHENZ_CONFIG_PATH"
	// EnvKeyZpuConfigPath the environment variable key for ZPU config path
	EnvKeyZpuConfigPath = "ZPU_CONFIG_PATH"
	// EnvKeyAgentConfigPath the environment variable key for agent config path
	EnvKeyAgentConfigPath = "AGENT_CONFIG_PATH"

	// DefaultZpeConfigPath the default value ZPE config
	DefaultZpeConfigPath = "config/zpe.conf"
	// DefaultAthenzConfigPath the default value for Athenz config
	DefaultAthenzConfigPath = "config/athenz.conf"
	// DefaultZpuConfigPath the default value for ZPU config
	DefaultZpuConfigPath = "config/zpu.conf"
	// DefaultAgentConfigPath the default value for agent config
	DefaultAgentConfigPath = "config/agent.json"
)

type (

	// Loader is the interface that wraps the basic config loader.
	Loader interface {
		// LoadConfig loads config file from filePath ro config object
		LoadConfig(config interface{}, filePath string) error
		// WithDefaultConfig performs default configuration to config loader
		WithDefaultConfig()
	}

	// viperLoader is the implementation of Loader interface for viper.
	//
	// It holds a viper struct for each config file, regarding the viper
	// field it's possible to have custom setting per config files.
	viperLoader struct {
		v *viper.Viper
	}
)

// NewConfigLoader creates new instance of Leader type
func NewConfigLoader() Loader {
	return &viperLoader{
		v: viper.New(),
	}
}

// LoadConfig loads config file from filePath ro config object
func (c *viperLoader) LoadConfig(config interface{}, filePath string) error {

	c.v.SetConfigFile(filePath)

	// load config to memory
	if err := c.v.ReadInConfig(); err != nil {
		return err
	}

	// unmarshall config to the input interface
	if err := c.v.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

// WithDefaultConfig performs default configuration to config loader
//
// WithDefaultConfig enables config file reload on new changes at runtime
// and environment variable detection.
func (c *viperLoader) WithDefaultConfig() {

	// load new config changes at runtime and notifies it
	c.v.WatchConfig()
	c.v.OnConfigChange(notify)

	// enable environment variable detection
	c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	c.v.AllowEmptyEnv(true)
	c.v.AutomaticEnv()
}

// notify notifies the new event
func notify(e fsnotify.Event) {
	// TODO Use log.Info instead
	fmt.Printf("Config file changed: %s\n", e.Name)
}
