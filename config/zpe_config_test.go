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
 *
 */

package config

import (
	"fmt"
	"gitlab.com/trialblaze/athenz-agent/common/util"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testConfigDirPrefix  = "config"
	testAthenzConfigFile = "athenz.json"
	testZpeConfigFile    = "zpe.conf"
)

func CreateFile(fileName, content string) error {
	if util.Exists(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			return fmt.Errorf("unable to remove file: %v, Error: %v", fileName, err)
		}
	}

	err := ioutil.WriteFile(fileName, []byte(content), 0755)
	if err != nil {
		return fmt.Errorf("unable to write file: %v, Error: %v", fileName, err)
	}

	return nil
}

func RemoveAll(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Println(err)
	}
}

func TestReadAthenzConfig(t *testing.T) {

	a := assert.New(t)
	athenzConfig := new(AthenzConfiguration)

	dir, err := ioutil.TempDir("./", testConfigDirPrefix)
	a.Nil(err)
	defer RemoveAll(dir)

	configPath := dir + "/" + testAthenzConfigFile

	// check config file with missing keys
	err = CreateFile(configPath, `{"zmsUrl":"zms_url","zmsPublicKeys":[{"id":"0","key":"zmsKey"}]}`)
	a.Nil(err)
	err = LoadAthenzConfig(athenzConfig, configPath)
	a.Nil(err)
	a.Equal("zms_url", athenzConfig.Properties.ZmsUrl)
	a.Empty(athenzConfig.Properties.ZtsUrl)
	a.Equal(1, len(athenzConfig.Properties.ZmsPublicKeys))
	a.Empty(athenzConfig.Properties.ZtsPublicKeys)
	a.Equal(0, len(athenzConfig.Properties.ZtsPublicKeys))
	a.Equal("0", athenzConfig.Properties.ZmsPublicKeys[0].Id)
	a.Equal("zmsKey", athenzConfig.Properties.ZmsPublicKeys[0].Key)

	//check if file content is incorrect
	athenzConfig = new(AthenzConfiguration)
	err = CreateFile(configPath, `"zmsUrl":"zms_url","zmsPublicKeys":[{"id":"0","key":"zmsKey"}]}`)
	a.Nil(err)
	err = LoadAthenzConfig(athenzConfig, configPath)
	a.NotNil(err)
	a.Empty(athenzConfig)

	//check if file content is correct
	athenzConfig = new(AthenzConfiguration)
	err = CreateFile(configPath, `{"zmsUrl":"zms_url","ztsUrl":"zts_url","ztsPublicKeys":[{"id":"0","key":"key0"}],"zmsPublicKeys":[{"id":"1","key":"key1"}]}`)
	a.Nil(err)
	err = LoadAthenzConfig(athenzConfig, configPath)
	a.Nil(err)
	a.Equal(athenzConfig.Properties.ZmsUrl, "zms_url")
	a.Equal(athenzConfig.Properties.ZtsUrl, "zts_url")
	a.Equal(1, len(athenzConfig.Properties.ZmsPublicKeys))
	a.Equal(1, len(athenzConfig.Properties.ZtsPublicKeys))
	a.Equal("0", athenzConfig.Properties.ZtsPublicKeys[0].Id)
	a.Equal("key0", athenzConfig.Properties.ZtsPublicKeys[0].Key)
	a.Equal("1", athenzConfig.Properties.ZmsPublicKeys[0].Id)
	a.Equal("key1", athenzConfig.Properties.ZmsPublicKeys[0].Key)
}

func TestReadZpeConfig(t *testing.T) {
	a := assert.New(t)

	zpeConfig := new(ZpeConfiguration)

	dir, err := ioutil.TempDir("./", testConfigDirPrefix)
	a.Nil(err)
	defer RemoveAll(dir)

	configPath := dir + "/" + testZpeConfigFile

	// check if file missing some key values
	err = CreateFile(configPath, `{"policy_files_dir": "./resource/policy","athenz_config_dir":""}`)
	a.Nil(err)
	err = LoadZpeConfig(zpeConfig, configPath)
	a.Nil(err)
	a.Empty(zpeConfig.Properties.CleanupTokenInterval)
	a.Empty(zpeConfig.Properties.AthenzConfigDir)
	a.Equal(zpeConfig.Properties.PolicyFilesDir, "./resource/policy")

	// check if file content is incorrect
	zpeConfig = new(ZpeConfiguration)
	err = CreateFile(configPath, `"policy_files_dir": "./resource/policy","cleanup_token_interval":600,"athenz_config_dir":"./resource"}`)
	a.Nil(err)
	err = LoadZpeConfig(zpeConfig, configPath)
	a.NotNil(err)
	a.Empty(zpeConfig)

	// check if file content is correct
	zpeConfig = new(ZpeConfiguration)
	err = CreateFile(configPath, `{"policy_files_dir": "./resource/policy","cleanup_token_interval":600,"athenz_config_dir":"./resource","athenz_token_no_expiry":true,"athenz_token_max_expiry":30,"allowed_offset":300}`)
	a.Nil(err)
	err = LoadZpeConfig(zpeConfig, configPath)
	a.Nil(err)
	a.Equal(zpeConfig.Properties.CleanupTokenInterval, int64(600))
	a.Equal(zpeConfig.Properties.PolicyFilesDir, "./resource/policy")
	a.Equal(zpeConfig.Properties.AthenzConfigDir, "./resource")
	a.True(zpeConfig.Properties.AthenzTokenNoExpiry)
	a.Equal(zpeConfig.Properties.AthenzTokenMaxExpiry, int64(30))
	a.Equal(zpeConfig.Properties.AllowedOffset, int64(300))
}
