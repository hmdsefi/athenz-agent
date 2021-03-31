/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/18/19
 * Time: 8:32 AM
 *
 * Description:
 *
 */

package auth

import (
	"encoding/json"
	"fmt"
	"github.com/ardielle/ardielle-go/rdl"
	"github.com/stretchr/testify/assert"
	"github.com/yahoo/athenz/clients/go/zts"
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	zpuUtil "github.com/yahoo/athenz/utils/zpe-updater/util"
	"gitlab.com/trialblaze/athenz-agent/cache"
	"gitlab.com/trialblaze/athenz-agent/config"
	"gitlab.com/trialblaze/athenz-agent/common/util"
	"gitlab.com/trialblaze/grpc-go/pkg/api/common/message/v1"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	testPolicyDir       = "../resource"
	testPolicyDirPrefix = "policy"
	testPolicyFile      = "../resource/test_data/angler.pol"
	testZmsPrivateKey0  = "../resource/test_data/zms_private_k0.pem"
	testZtsPrivateKey0  = "../resource/test_data/zts_private_k0.pem"
	testAthenzConfig    = "../resource/test_data/athenz.conf"
	testZpeConfig       = "../resource/test_data/zpe.conf"
)

var testTempFolder string

func preparePolicyFiles(expiry time.Time) error {
	readFile, err := os.OpenFile(testPolicyFile, os.O_RDONLY, 0444)
	defer readFile.Close()
	if err != nil {
		return fmt.Errorf("cannot open file: %v , Error: %v", testPolicyFile, err)
	}

	var domainSignedPolicyData *zts.DomainSignedPolicyData
	err = json.NewDecoder(readFile).Decode(&domainSignedPolicyData)
	if err != nil {
		return fmt.Errorf("unable to decode policy file: %v, Error: %v", testPolicyFile, err)
	}

	if expiry.UnixNano() > 0 {
		expiry = expiry.Add(48 * time.Hour)
		domainSignedPolicyData.SignedPolicyData.Expires = rdl.Timestamp{expiry}
	}

	zmsData, err := ioutil.ReadFile(testZmsPrivateKey0)
	if err != nil {
		return fmt.Errorf("cannot open zms private key file")
	}

	signer, _ := zmssvctoken.NewSigner(zmsData)
	policyData, _ := zpuUtil.ToCanonicalString(domainSignedPolicyData.SignedPolicyData.PolicyData)
	signature, _ := signer.Sign(policyData)
	domainSignedPolicyData.SignedPolicyData.ZmsSignature = signature
	domainSignedPolicyData.SignedPolicyData.ZmsKeyId = "0"

	ztsData, err := ioutil.ReadFile(testZtsPrivateKey0)
	if err != nil {
		return fmt.Errorf("cannot open zts private key file")
	}

	signer, _ = zmssvctoken.NewSigner(ztsData)
	policyData, _ = zpuUtil.ToCanonicalString(domainSignedPolicyData.SignedPolicyData)
	signature, _ = signer.Sign(policyData)
	domainSignedPolicyData.Signature = signature
	domainSignedPolicyData.KeyId = "0"

	testTempFolder, err = ioutil.TempDir(testPolicyDir, testPolicyDirPrefix)
	if err != nil {
		return fmt.Errorf("unable to create policy directory")
	}

	data, _ := json.Marshal(domainSignedPolicyData)
	err = util.CreateFile(testTempFolder+"/angler.pol", string(data))
	if err != nil {
		return fmt.Errorf("unable to create policy file")
	}
	return nil
}

func createRoleToken(role, domain string) string {
	generatedToken := strconv.FormatInt((util.CurrentTimeMillis()/1000-30)*int64(time.Second), 10)
	expiration := strconv.FormatInt((util.CurrentTimeMillis()/1000+300)*int64(time.Second), 10)
	signedToken := "v=S1;d=" + domain + ";h=localhost" + ";r=" + role +
		";t=" + generatedToken + ";e=" + expiration + ";k=0"

	data, _ := ioutil.ReadFile(testZtsPrivateKey0)

	signer, _ := zmssvctoken.NewSigner(data)
	signature, _ := signer.Sign(signedToken)

	signedToken = signedToken + ";s=" + signature

	return signedToken
}

func TestPermissionService_CheckAccessWithTokenPolicyFileExpired(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Time{})
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "read", Resource: "angler:stuff",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(9))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "read", Resource: "angler:stuff",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenDeny(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "throw", Resource: "angler:stuff",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(1))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenStartWith(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "fish", Resource: "angler:stockedpondBigBassLake",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenWildcardDeny(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("managerkernco", "angler")

	request := &v1.AccessCheckRequest{Access: "manage", Resource: "angler:pondsVenturaCounty",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(1))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenWildcardAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("managerkernco", "angler")

	request := &v1.AccessCheckRequest{Access: "manage", Resource: "angler:pondsKernCounty",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenMatchAllAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("matchall", "angler")

	request := &v1.AccessCheckRequest{Access: "all", Resource: "angler:anything",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenMatchRegexAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("matchregex", "angler")

	request := &v1.AccessCheckRequest{Access: "regex", Resource: "angler:nhllllllkings",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow1(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:oretech",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow2(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:orecommit",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow3(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:orec",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow4(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.Nil(err)

	config.KeyStore, _ = config.LoadAthenzConfig(testAthenzConfig)
	config.ZConfig, _ = config.LoadZpeConfig(testZpeConfig)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:ored",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.Nil(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}
