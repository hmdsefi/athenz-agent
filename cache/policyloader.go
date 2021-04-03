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
 * In here we try to cache policy files into memory to check
 * access much more faster and prevent to opening and closing
 * files. So, we will Cache them in some maps to have quick
 * access to domains and their roles. Also, we will cache
 * roleTokens that want to give a access to a resource, we
 * will cleanup cached roleTokens in some iteration to prevent
 * caching expired roleTokens.
 *
 */

package cache

import (
	"encoding/json"
	"fmt"
	"github.com/yahoo/athenz/clients/go/zts"
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	zpuUtil "github.com/yahoo/athenz/utils/zpe-updater/util"
	"gitlab.com/trialblaze/athenz-agent/common"
	"gitlab.com/trialblaze/athenz-agent/common/log"
	"gitlab.com/trialblaze/athenz-agent/config"
	"gitlab.com/trialblaze/athenz-agent/matcher"
	"gitlab.com/trialblaze/athenz-agent/token"

	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	logger           = log.GetLogger(common.GolangFileName())
	fileStatusMap    = make(map[string]*zpeFileStatus)
	lastTokenCleanup = common.CurrentTimeMillis()
	PolicyDirectory  string

	// key is the domain name, value is a map keyed by role name with list of assertions
	DomainStandardRoleAllowMap = make(map[string]*RoleMap)

	// wild card role map, keys and values same as domRoleMap above
	DomainWildcardRoleAllowMap = make(map[string]*RoleMap)

	// key is the domain name, value is a map keyed by role name with list of assertions
	DomainStandardRoleDenyMap = make(map[string]*RoleMap)

	// wild card role map, keys and values same as domRoleMap above
	DomainWildcardRoleDenyMap = make(map[string]*RoleMap)

	// cache of active Role Tokens
	RoleTokenCacheMap = make(map[string]*token.RoleToken)
)

type zpeFileStatus struct {
	fileName         string
	domainName       string
	lastModifiedDate time.Time
	isValidPolFile   bool
}

type RoleMap struct {
	Expiry      int64
	RoleDataMap map[string][]map[string]interface{}
}

func getMatchObject(value string) matcher.ZpeMatch {
	if value == "*" {
		return matcher.ZpeMatchAll{}
	} else {
		anyCharMatch := strings.Index(value, "*")
		singleCharMatch := strings.Index(value, "?")

		if anyCharMatch == -1 && singleCharMatch == -1 {
			return matcher.ZpeMatchEqual{MatchValue: value}
		} else if anyCharMatch == len(value)-1 && singleCharMatch == -1 {
			return matcher.ZpeMatchStartsWith{Prefix: value[:anyCharMatch]}
		} else {
			regexMatcher, err := matcher.NewZpeMatchRegex(value)
			if err != nil {
				logger.Error(fmt.Sprintf("unable to create pattern for '%s', error: %s", value, err.Error()))
			}
			return regexMatcher
		}
	}
}

// Process the given policy file list and determine if any of the
// policy domain files have been updated. New ones will be loaded
// into the policy domain map.
func LoadDB(files []os.FileInfo) {
	if files == nil {
		logger.Info("loadDb: no policy files to load")
		return
	}
	for _, policyFile := range files {
		fileStatus := fileStatusMap[policyFile.Name()]
		if fileStatus != nil {

			//	check if file does not exist
			if _, err := os.Stat(PolicyDirectory + "/" + fileStatus.fileName); os.IsNotExist(err) {
				delete(fileStatusMap, policyFile.Name())
				if !fileStatus.isValidPolFile || fileStatus.domainName == "" {
					continue
				}

				DomainStandardRoleAllowMap[fileStatus.domainName] = new(RoleMap)
				DomainWildcardRoleAllowMap[fileStatus.domainName] = new(RoleMap)
				DomainStandardRoleDenyMap[fileStatus.domainName] = new(RoleMap)
				DomainWildcardRoleDenyMap[fileStatus.domainName] = new(RoleMap)
				continue
			}

			// check if file was modified since last time it was loaded
			if policyFile.ModTime().UnixNano() <= fileStatus.lastModifiedDate.UnixNano() {
				if fileStatus.isValidPolFile {
					continue
				}
			}

		} else {
			fileStatusMap[policyFile.Name()] = &zpeFileStatus{fileName: policyFile.Name(),
				lastModifiedDate: policyFile.ModTime()}
		}
		err := loadFile(policyFile)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

// Loads and parses the given file. It will create the domain assertion
// list per role and put it into the domain policy maps(domRoleMap, domWildcardRoleMap).
func loadFile(file os.FileInfo) error {

	path := PolicyDirectory + "/" + file.Name()
	fileInfo, err := os.Stat(path)
	if err != nil {
		return common.Errorf("unable to load file info: %s, error: %s", path, err)
	}

	readFile, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return common.Errorf("unable to open file: %s , error: %s", path, err)
	}
	defer func() {
		err := readFile.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	var domainSignedPolicyData *zts.DomainSignedPolicyData
	err = json.NewDecoder(readFile).Decode(&domainSignedPolicyData)
	if err != nil {
		return common.Errorf("unable to decode policy file: %s, error: %s", path, err.Error())
	}

	if domainSignedPolicyData == nil {
		//	mark this file as an invalid file
		fileStatus := fileStatusMap[fileInfo.Name()]
		if fileStatus != nil {
			fileStatus.isValidPolFile = false
		}
		return common.Errorf("unable to decode policy file: %s", path)
	}

	// first let's verify the ZTS signature for our policy file
	signedPolicyData := domainSignedPolicyData.SignedPolicyData

	verified := false
	input, err := zpuUtil.ToCanonicalString(signedPolicyData)
	if err != nil {
		return common.Errorf("unable to convert to string, error: %s", err.Error())
	}

	pubKey := config.KeyStore.GetZtsPublicKey(domainSignedPolicyData.KeyId)
	ztsKey, err := new(zmssvctoken.YBase64).DecodeString(pubKey)
	if err != nil {
		return common.Errorf("verification of data with zts key having id: '%s' failed, error: %s",
			domainSignedPolicyData.KeyId, err.Error())
	}
	err = common.Verify(input, domainSignedPolicyData.Signature, string(ztsKey))
	if err == nil {
		verified = true
	} else {
		logger.Error("invalid policy, error: " + err.Error())
	}

	var policyData *zts.PolicyData
	if verified {
		policyData = signedPolicyData.PolicyData

		inputPolicy, err := zpuUtil.ToCanonicalString(policyData)
		if err != nil {
			return common.Errorf("unable to convert to string, error: %s", err)
		}

		zmsKey, err := new(zmssvctoken.YBase64).DecodeString(config.KeyStore.GetZmsPublicKey(signedPolicyData.ZmsKeyId))
		if err != nil {
			return common.Errorf("verification of data with zms key having id:'%s' failed, error: %s",
				signedPolicyData.ZmsKeyId, err)
		}

		err = common.Verify(inputPolicy, signedPolicyData.ZmsSignature, string(zmsKey))
		if err != nil {
			verified = false
		}
	}

	if !verified || policyData == nil {
		//	mark this file as an invalid file
		fileStatus := fileStatusMap[fileInfo.Name()]
		if fileStatus != nil {
			fileStatus.isValidPolFile = false
		}
		return common.Errorf("policy file is invalid: %s", path)
	}

	domainName := string(policyData.Domain)

	roleStandardAllowMap := make(map[string][]map[string]interface{})
	roleWildcardAllowMap := make(map[string][]map[string]interface{})
	roleStandardDenyMap := make(map[string][]map[string]interface{})
	roleWildcardDenyMap := make(map[string][]map[string]interface{})
	for _, policy := range policyData.Policies {
		for _, assertion := range policy.Assertions {
			strAssert := make(map[string]interface{})
			strAssert[common.ZpeFieldPolicyName] = policy.Name
			strAssert[common.ZpeActionMatchStruct] = getMatchObject(assertion.Action)

			rsrc := common.StripDomainPrefix(assertion.Resource, domainName, assertion.Resource)
			strAssert[common.ZpeFieldResource] = rsrc
			strAssert[common.ZpeResourceMatchStruct] = getMatchObject(rsrc)

			pRoleName := common.StripDomainPrefix(assertion.Role, domainName, assertion.Role)
			reg := regexp.MustCompile("^role.")
			pRoleName = reg.ReplaceAllString(pRoleName, "$1")
			strAssert[common.ZpeFieldRole] = pRoleName
			matchStruct := getMatchObject(pRoleName)
			strAssert[common.ZpeRoleMatchStruct] = matchStruct

			if assertion.Effect != nil && assertion.Effect.String() == "DENY" {
				if reflect.TypeOf(matchStruct).Name() == "ZpeMatchEqual" {
					computeIfAbsent(pRoleName, roleStandardDenyMap, strAssert)
				} else {
					computeIfAbsent(pRoleName, roleWildcardDenyMap, strAssert)
				}
			} else {
				if reflect.TypeOf(matchStruct).Name() == "ZpeMatchEqual" {
					computeIfAbsent(pRoleName, roleStandardAllowMap, strAssert)
				} else {
					computeIfAbsent(pRoleName, roleWildcardAllowMap, strAssert)
				}
			}

		}
	}

	fileStatus := fileStatusMap[fileInfo.Name()]
	if fileStatus != nil {
		fileStatus.isValidPolFile = true
		fileStatus.domainName = domainName
	}

	expires := signedPolicyData.Expires.UnixNano()

	DomainStandardRoleAllowMap[domainName] = &RoleMap{Expiry: expires, RoleDataMap: roleStandardAllowMap}
	DomainWildcardRoleAllowMap[domainName] = &RoleMap{Expiry: expires, RoleDataMap: roleWildcardAllowMap}
	DomainStandardRoleDenyMap[domainName] = &RoleMap{Expiry: expires, RoleDataMap: roleStandardDenyMap}
	DomainWildcardRoleDenyMap[domainName] = &RoleMap{Expiry: expires, RoleDataMap: roleWildcardDenyMap}

	return nil
}

// this method will check if there is a slice for the
// key then append new item to that slice, else create
// new slice and append new item to that
func computeIfAbsent(key string, roleMap map[string][]map[string]interface{}, mapSlice map[string]interface{}) {
	if assertSlice, ok := roleMap[key]; ok {
		assertSlice = append(assertSlice, mapSlice)
		roleMap[key] = assertSlice
	} else {
		newMapSlice := make([]map[string]interface{}, 0)
		roleMap[key] = append(newMapSlice, mapSlice)
	}
}

// lets cleanup the every thing we cache and
// be prepared for caching policies
func CleanupRoleTokenCache() {
	//is it time to cleanup
	now := common.CurrentTimeMillis()
	if now < int64(time.Duration(config.ZpeConfig.Properties.CleanupTokenInterval)*time.Microsecond)+lastTokenCleanup {
		return
	}

	expired := make([]string, 0)
	for key, roleToken := range RoleTokenCacheMap {
		if roleToken == nil {
			continue
		}
		expiry := roleToken.ExpiryTime / int64(time.Millisecond)
		if expiry < now {
			expired = append(expired, key)
		}
	}

	// now we will remove expired roleTokens
	for _, key := range expired {
		delete(RoleTokenCacheMap, key)
	}
	// update last cleanup time
	lastTokenCleanup = now
}
