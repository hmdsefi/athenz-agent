/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
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
	"errors"
	"fmt"
)

var (
	KeyStore = &AthenzConfiguration{}
	ZpeConfig  = &ZpeConfiguration{}
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
		PolicyFilesDir       string
		CleanupTokenInterval int64 // in seconds format
		AthenzConfigDir      string
		AthenzTokenNoExpiry  bool
		AthenzTokenMaxExpiry int64 // in days format
		AllowedOffset        int64 // in seconds format
		GRPCServerPort       string
		CertFilePath         string // TLS cert, this will be used for communicating with ZTS server
		KeyFilePath          string // Key for the TLS cert, this will be used for communicating with ZTS server
		DomainName           string // domain name that belong to the service that use this agent
		ServiceName          string // the name that server register in athenz with
		RoleNames            string // list of service comma separated role names
		TokenExpirationMin   int32  // in minutes format, It will be used for getting roleToken from ZTS server
		TokenExpirationMax   int32  // in minutes format, It will be used for getting roleToken from ZTS server
		KeyVersion           string // The key-version should be the same string that was used to register the key with Athenz
		NTokenExpiration     int64  // the duration for which the token is valid, in minutes format
		ZpuDownloadInterval  int64  // in seconds format
	}

	PublicKeys struct {
		Id  string
		Key string
	}

	athenzProperties struct {
		ZtsUrl        string
		ZmsUrl        string
		ZtsPublicKeys []PublicKeys
		ZmsPublicKeys []PublicKeys
	}
)

// LoadAthenzConfig reads config file from a specific address and
// loads it into a AthenzConfiguration object
func LoadAthenzConfig(athenzConfig *AthenzConfiguration, filePath string) error {

	// load config properties into athenzProperties
	athenzConfig.Properties = new(athenzProperties)
	athenzConfig.loader = NewConfigLoader()
	if err := athenzConfig.loader.LoadConfig(athenzConfig.Properties, filePath); err != nil {
		return errors.New(fmt.Sprintf("LoadAthenzConfig: unable to load config from %s : %s", filePath, err.Error()))
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

// LoadZpeConfig reads config file from a specific address and
// loads it into a ZpeConfiguration object
func LoadZpeConfig(zpeConfig *ZpeConfiguration, filePath string) error {

	// load config properties into zpeProperties
	zpeConfig.Properties = new(zpeProperties)
	zpeConfig.loader = NewConfigLoader()
	if err := zpeConfig.loader.LoadConfig(zpeConfig.Properties, filePath); err != nil {
		return errors.New(fmt.Sprintf("LoadZpeConfig: unable to load config from %s : %s", filePath, err.Error()))
	}

	// use default configuration for config loader
	zpeConfig.loader.WithDefaultConfig()

	return nil
}
