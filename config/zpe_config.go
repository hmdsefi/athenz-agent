/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/12/19
 * Time: 8:37 AM
 *
 * Description:
 * To use `athenz-agent` you need two main config, zpe.conf and
 * athenz.conf. Yo can use this file to read config files and
 * move them into memory. Use LoadAthenzConfig and give it the
 * athenz.conf address and it will return a AthenzConfig object.
 * Use LoadZpeConfig and give it zpe.conf and it will return a
 * ZpeConfig object to you.
 *
 */

package config

import (
	"fmt"
	"github.com/hamed-yousefi/athenz-agent/common"
)

var (
	KeyStore  = newAthenzConfiguration()
	ZpeConfig = newZpeConfiguration()
)

type (
	AthenzConfiguration struct {
		loader     Loader
		Properties *athenzProperties
	}

	ZpeConfiguration struct {
		loader     Loader
		Properties *zpeProperties
	}

	zpeProperties struct {
		PolicyFilesDir       string `mapstructure:"policy_files_dir"`
		// in seconds format
		CleanupTokenInterval int64  `mapstructure:"cleanup_token_interval"`
		AthenzConfigDir      string `mapstructure:"athenz_config_dir"`
		AthenzTokenNoExpiry  bool   `mapstructure:"athenz_token_no_expiry"`
		// in days format
		AthenzTokenMaxExpiry int64  `mapstructure:"athenz_token_max_expiry"`
		// in seconds format
		AllowedOffset        int64  `mapstructure:"allowed_offset"`
		// TLS cert, this will be used for communicating with ZTS server
		CertFilePath         string `mapstructure:"cert_file_path"`
		// Key for the TLS cert, this will be used for communicating with ZTS server
		KeyFilePath          string `mapstructure:"key_file_path"`
		// domain name that belong to the service that use this agent
		DomainName           string `mapstructure:"domain_name"`
		// the name that server register in athenz with
		ServiceName          string `mapstructure:"service_name"`
		// list of service comma separated role names
		RoleNames            string `mapstructure:"role_names"`
		// in minutes format, It will be used for getting roleToken from ZTS server
		TokenExpirationMin   int32  `mapstructure:"token_expiration_min"`
		// in minutes format, It will be used for getting roleToken from ZTS server
		TokenExpirationMax   int32  `mapstructure:"token_expiration_max"`
		// The key-version should be the same string that was used to register the key with Athenz
		KeyVersion           string `mapstructure:"key_version"`
		// the duration for which the token is valid, in minutes format
		NTokenExpiration     int64  `mapstructure:"ntoken_expiration"`
		// in seconds format
		ZpuDownloadInterval  int64 `mapstructure:"zpu_download_interval"`
	}

	PublicKeys struct {
		Id  string
		Key string
	}

	athenzProperties struct {
		ZtsUrl        string
		ZmsUrl        string
		ZtsPublicKeys []PublicKeys `mapstructure:"ztsPublicKeys"`
		ZmsPublicKeys []PublicKeys`mapstructure:"zmsPublicKeys"`
	}
)

// newAthenzConfiguration creates a new instance of AthenzConfiguration with
// an empty properties to prevent nil pointer exception.
func newAthenzConfiguration() *AthenzConfiguration {
	return &AthenzConfiguration{Properties: new(athenzProperties)}
}

// LoadGlobalAthenzConfig loads config file from input path into the global
// variable KeyStore.
func LoadGlobalAthenzConfig(filePath string) error {
	return LoadAthenzConfig(KeyStore, filePath)
}

// LoadAthenzConfig reads config file from a specific address and  loads it into
// a AthenzConfiguration object.
func LoadAthenzConfig(athenzConfig *AthenzConfiguration, filePath string) error {

	// load config properties into athenzProperties
	athenzConfig.Properties = new(athenzProperties)
	athenzConfig.loader = NewConfigLoader()
	if err := athenzConfig.loader.LoadConfig(athenzConfig.Properties, filePath); err != nil {
		return common.Errorf("unable to load config from %s : %s", filePath, err.Error())
	}

	// use default configuration for config loader
	athenzConfig.loader.WithDefaultConfig()

	return nil
}

func (config AthenzConfiguration) GetZtsPublicKey(id string) string {
	for _, ztsPublicKeys := range config.Properties.ZtsPublicKeys {
		if ztsPublicKeys.Id == id {
			return ztsPublicKeys.Key
		}
	}
	return ""
}

func (config AthenzConfiguration) GetZmsPublicKey(id string) string {
	for _, zmsPublicKey := range config.Properties.ZmsPublicKeys {
		if zmsPublicKey.Id == id {
			return zmsPublicKey.Key
		}
	}
	return ""
}

// newZpeConfiguration creates a new instance of ZpeConfiguration with
// an empty properties to prevent nil pointer exception.
func newZpeConfiguration() *ZpeConfiguration {
	return &ZpeConfiguration{Properties: new(zpeProperties)}
}

// LoadGlobalZpeConfig loads config file from input path into the global
// variable ZpeConfig.
func LoadGlobalZpeConfig(filePath string) error {
	return LoadZpeConfig(ZpeConfig, filePath)
}

// LoadZpeConfig reads config file from a specific address and
// loads it into a ZpeConfiguration object
func LoadZpeConfig(zpeConfig *ZpeConfiguration, filePath string) error {

	// load config properties into zpeProperties
	zpeConfig.Properties = new(zpeProperties)
	zpeConfig.loader = NewConfigLoader()
	if err := zpeConfig.loader.LoadConfig(zpeConfig.Properties, filePath); err != nil {
		return common.Errorf(fmt.Sprintf("unable to load config from %s: %s", filePath, err.Error()))
	}

	// use default configuration for config loader
	zpeConfig.loader.WithDefaultConfig()

	return nil
}
