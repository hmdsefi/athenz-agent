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
 *
 */

package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"github.com/hamed-yousefi/athenz-agent/config"
	"github.com/hamed-yousefi/athenz-agent/token"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

const (
	configDirPrefix  = "config"
	policyDirPrefix  = "policy"
	polFile          = "test.pol"
)

func setup() {
	log.NewLogrusInitializer().InitialLog(log.Info)
}

func TestGetMatchObject(t *testing.T) {
	a := assert.New(t)

	matchObject := getMatchObject("*")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchAll")

	matchObject = getMatchObject("**")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchRegex")

	matchObject = getMatchObject("?*")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchRegex")

	matchObject = getMatchObject("?")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchRegex")

	matchObject = getMatchObject("test?again*")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchRegex")

	matchObject = getMatchObject("*test")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchRegex")

	matchObject = getMatchObject("test")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchEqual")

	matchObject = getMatchObject("(test|again)")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchEqual")

	matchObject = getMatchObject("test*")
	a.True(reflect.TypeOf(matchObject).Name() == "ZpeMatchStartsWith")
}

func TestLoadDBNull(t *testing.T) {
	setup()
	LoadDB(nil)
}

func TestLoadDB(t *testing.T) {
	setup()
	a := assert.New(t)

	policyDir, err := ioutil.TempDir("./", policyDirPrefix)
	a.NoError(err)
	defer func() {
		err := common.RemoveAll(policyDir)
		if err != nil {
			common.Fatal(err.Error())
		}
	}()
	PolicyDirectory = policyDir

	policyPath := policyDir + "/" + polFile
	err = common.CreateFile(policyPath, `{"signedPolicyData":{"expires":"2017-06-09T06:11:12.125Z","modified" : "2017-06-02T06:11:12.125Z","policyData":{"domain":"sys.auth","policies":[{"assertions":[{"action":"*","effect":"ALLOW","resource":"*","role":"sys.auth:role.admin"},{"action":"*","effect":"DENY","resource":"*","role":"sys.auth:role.non-admin"}],"name":"sys.auth:policy.admin"}]},"zmsKeyId":"0","zmsSignature":"Y2HuXmgL86PL1WnleGFHwPmNEqUdWgDxmmIsDnF5f5oqakacqTtwt9JNqDV9nuJ7LnKl3zsZoDQSAtcHMu4IGA--"},"signature":"XJnQ4t33D4yr7NtUjLaWhXULFr76z.z0p3QV4uCkA5KR9L4liVRmICYwVmnXxvHAlImKlKLv7sbIHNsjBfGfCw--","keyId": "0"}`)
	a.NoError(err)

	// check if zms and zts public keys not exist input must
	// be invalid
	files, _ := common.LoadFileStatus(policyDir)
	LoadDB(files)
	a.Len(DomainWildcardRoleDenyMap, 0)
	a.Len(DomainStandardRoleAllowMap, 0)
	a.Len(DomainWildcardRoleAllowMap, 0)
	a.Len(DomainStandardRoleDenyMap, 0)
	a.False(fileStatusMap[polFile].isValidPolFile)

	// use athenz config file to verify input and signature
	// and then cache the policies in memory
	configDir, err := ioutil.TempDir("./", configDirPrefix)
	a.Nil(err)
	defer common.RemoveAll(configDir)
	configPath := configDir + "/athenz.json"
	err = common.CreateFile(configPath, `{"zmsUrl":"https://dev.zms.athenzcompany.com:4443/","ztsUrl":"https://dev.zts.athenzcompany.com:4443/","ztsPublicKeys":[{"id":"0","key":"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTHpmU09UUUpmRW0xZW00TDNza3lOVlEvYngwTU9UcQphK1J3T0gzWmNNS3lvR3hPSm85QXllUmE2RlhNbXZKSkdZczVQMzRZc3pGcG5qMnVBYmkyNG5FQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"},{"id":"1","key": "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FETGlLY1hjUDlrMWRJcGU4bm1OS3pBaWpGcApuY0VWbEFveS8xcHordE5ETjExcDQ0MTJEREhXejhFSUNiVkE0RE16Wm1ta09URFdlUDBQSWdnNTg0RlF1SGpsCmsyOWU4VjJXT3pqQWZybGlad0dKbm1mdlBhb3FOQkNhZDI3cWFubm1MOVU3cTcvSEdRWmpMeGdoaXhGa0FtczEKaHFlbnlkb2JSVkhheHV3cDB3SURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"}],"zmsPublicKeys":[{"id":"0","key":"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTHpmU09UUUpmRW0xZW00TDNza3lOVlEvYngwTU9UcQphK1J3T0gzWmNNS3lvR3hPSm85QXllUmE2RlhNbXZKSkdZczVQMzRZc3pGcG5qMnVBYmkyNG5FQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"},{"id":"1","key":"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTHpmU09UUUpmRW0xZW00TDNza3lOVlEvYngwTU9UcQphK1J3T0gzWmNNS3lvR3hPSm85QXllUmE2RlhNbXZKSkdZczVQMzRZc3pGcG5qMnVBYmkyNG5FQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"}]}`)
	a.Nil(err)
	if err := config.LoadGlobalAthenzConfig(configPath); err != nil {
		common.Fatalf("unable to load config, %s: ", err)
	}

	files, _ = common.LoadFileStatus(policyDir)
	LoadDB(files)
	a.True(fileStatusMap[polFile].isValidPolFile)

}

func TestCleanupRoleTokenCache(t *testing.T) {
	setup()
	a := assert.New(t)

	dir, err := ioutil.TempDir("./", configDirPrefix)
	a.NoError(err)
	defer func() {
		err := common.RemoveAll(dir)
		if err != nil {
			common.Fatal(err.Error())
		}
	}()

	configPath := dir + "/zpe.json"
	err = common.CreateFile(configPath, `{"policy_files_dir": "./resource/policy","cleanup_token_interval":10,"athenz_config_dir":"./resource"}`)
	a.NoError(err)
	err = config.LoadGlobalZpeConfig(configPath)
	a.NoError(err)

	lastTokenCleanup = common.CurrentTimeMillis()
	oldLTC := lastTokenCleanup
	RoleTokenCacheMap["role1"] = &token.RoleToken{ExpiryTime: time.Now().UnixNano() - (10 * int64(time.Second))}
	RoleTokenCacheMap["role2"] = &token.RoleToken{ExpiryTime: time.Now().UnixNano() - (5 * int64(time.Second))}
	RoleTokenCacheMap["role3"] = &token.RoleToken{ExpiryTime: time.Now().UnixNano() + (5 * int64(time.Second))}
	RoleTokenCacheMap["role4"] = &token.RoleToken{ExpiryTime: time.Now().UnixNano() + (10 * int64(time.Second))}

	// this is not right time to cleanup
	CleanupRoleTokenCache()
	a.True(oldLTC == lastTokenCleanup)
	a.Len(RoleTokenCacheMap, 4)

	lastTokenCleanup = common.CurrentTimeMillis() - int64(time.Duration(15)*time.Microsecond)
	oldLTC = lastTokenCleanup

	// this is right time to cleanup cached roles
	CleanupRoleTokenCache()
	a.True(lastTokenCleanup > oldLTC)
	a.Len(RoleTokenCacheMap, 2)
	_, ok := RoleTokenCacheMap["role1"]
	a.False(ok)
	_, ok = RoleTokenCacheMap["role2"]
	a.False(ok)
	_, ok = RoleTokenCacheMap["role3"]
	a.True(ok)
	_, ok = RoleTokenCacheMap["role4"]
	a.True(ok)
}
