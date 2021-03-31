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
	zpeconst "gitlab.com/trialblaze/athenz-agent"
	"gitlab.com/trialblaze/athenz-agent/config"
	"gitlab.com/trialblaze/athenz-agent/matcher"
	"gitlab.com/trialblaze/athenz-agent/token"
	"gitlab.com/trialblaze/athenz-agent/common/util"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var fileStatusMap = make(map[string]*zpeFileStatus)
var lastTokenCleanup = util.CurrentTimeMillis()
var PolicyDirectory string

// key is the domain name, value is a map keyed by role name with list of assertions
var DomainStandardRoleAllowMap = make(map[string]*RoleMap)

// wild card role map, keys and values same as domRoleMap above
var DomainWildcardRoleAllowMap = make(map[string]*RoleMap)

// key is the domain name, value is a map keyed by role name with list of assertions
var DomainStandardRoleDenyMap = make(map[string]*RoleMap)

// wild card role map, keys and values same as domRoleMap above
var DomainWildcardRoleDenyMap = make(map[string]*RoleMap)

// cache of active Role Tokens
var RoleTokenCacheMap = make(map[string]*token.RoleToken)

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
			patter, _ := regexp.Compile(value)
			return matcher.ZpeMatchRegex{Pattern: patter}
		}
	}
}

// Process the given policy file list and determine if any of the
// policy domain files have been updated. New ones will be loaded
// into the policy domain map.
func LoadDB(files []os.FileInfo) {
	if files == nil {
		log.Println("loadDb: no policy files to load")
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
			log.Println(err)
		}
	}
}

// Loads and parses the given file. It will create the domain assertion
// list per role and put it into the domain policy maps(domRoleMap, domWildcardRoleMap).
func loadFile(file os.FileInfo) error {

	path := PolicyDirectory + "/" + file.Name()
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("loadFile: unable to load file info: %v, Error: %v", path, err)
	}

	readFile, err := os.OpenFile(path, os.O_RDONLY, 0444)
	defer readFile.Close()
	if err != nil {
		return fmt.Errorf("loadFile: Cannot open file: %v , Error: %v", path, err)
	}

	var domainSignedPolicyData *zts.DomainSignedPolicyData
	err = json.NewDecoder(readFile).Decode(&domainSignedPolicyData)
	if err != nil {
		return fmt.Errorf("loadFile: Unable to decode policy file: %v, Error: %v", path, err)
	}

	if domainSignedPolicyData == nil {
		//	mark this file as an invalid file
		fileStatus := fileStatusMap[fileInfo.Name()]
		if fileStatus != nil {
			fileStatus.isValidPolFile = false
		}
		return fmt.Errorf("loadFile: Unable to decode policy file: %v, Error: %v", path, err)
	}

	// first let's verify the ZTS signature for our policy file
	signedPolicyData := domainSignedPolicyData.SignedPolicyData

	verified := false
	input, err := zpuUtil.ToCanonicalString(signedPolicyData)
	if err != nil {
		return fmt.Errorf("loadFile: Unable to convert to string, Error: %v", err)
	}

	pubKey := config.KeyStore.GetZtsPublicKey(domainSignedPolicyData.KeyId)
	ztsKey, err := new(zmssvctoken.YBase64).DecodeString(pubKey)
	if err != nil {
		return fmt.Errorf("loadFile: verification of data with zts key having id:\"%v\" failed, Error :%v",
			domainSignedPolicyData.KeyId, err)
	}
	err = util.Verify(input, domainSignedPolicyData.Signature, string(ztsKey))
	if err == nil {
		verified = true
	}

	var policyData *zts.PolicyData
	if verified {
		policyData = signedPolicyData.PolicyData

		inputPolicy, err := zpuUtil.ToCanonicalString(policyData)
		if err != nil {
			return fmt.Errorf("loadFile: Unable to convert to string, Error: %v", err)
		}

		zmsKey, err := new(zmssvctoken.YBase64).DecodeString(config.KeyStore.GetZmsPublicKey(signedPolicyData.ZmsKeyId))
		if err != nil {
			return fmt.Errorf("loadFile: verification of data with zms key having id:\"%v\" failed, Error :%v",
				signedPolicyData.ZmsKeyId, err)
		}

		err = util.Verify(inputPolicy, signedPolicyData.ZmsSignature, string(zmsKey))
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
		return fmt.Errorf("loadFile: Policy file is invalid: %v", path)
	}

	domainName := string(policyData.Domain)

	roleStandardAllowMap := make(map[string][]map[string]interface{})
	roleWildcardAllowMap := make(map[string][]map[string]interface{})
	roleStandardDenyMap := make(map[string][]map[string]interface{})
	roleWildcardDenyMap := make(map[string][]map[string]interface{})
	for _, policy := range policyData.Policies {
		for _, assertion := range policy.Assertions {
			strAssert := make(map[string]interface{})
			strAssert[zpeconst.ZpeFieldPolicyName] = policy.Name
			strAssert[zpeconst.ZpeActionMatchStruct] = getMatchObject(assertion.Action)

			rsrc := util.StripDomainPrefix(assertion.Resource, domainName, assertion.Resource)
			strAssert[zpeconst.ZpeFieldResource] = rsrc
			strAssert[zpeconst.ZpeResourceMatchStruct] = getMatchObject(rsrc)

			pRoleName := util.StripDomainPrefix(assertion.Role, domainName, assertion.Role)
			reg := regexp.MustCompile("^role.")
			pRoleName = reg.ReplaceAllString(pRoleName, "$1")
			strAssert[zpeconst.ZpeFieldRole] = pRoleName
			matchStruct := getMatchObject(pRoleName)
			strAssert[zpeconst.ZpeRoleMatchStruct] = matchStruct

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
	now := util.CurrentTimeMillis()
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
