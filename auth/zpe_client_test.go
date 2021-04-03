/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
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
	"github.com/ardielle/ardielle-go/rdl"
	"github.com/stretchr/testify/assert"
	"github.com/yahoo/athenz/clients/go/zts"
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	zpuUtil "github.com/yahoo/athenz/utils/zpe-updater/util"
	"gitlab.com/trialblaze/athenz-agent/cache"
	"gitlab.com/trialblaze/athenz-agent/common"
	"gitlab.com/trialblaze/athenz-agent/common/log"
	"gitlab.com/trialblaze/athenz-agent/config"
	"gitlab.com/trialblaze/grpc-go/pkg/api/common/message/v1"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	tmpDir          = "testdata"
	policyDirPrefix = "policy"
	policyFile      = "testdata/angler.pol"
	zmsPrivateKey0  = "testdata/zms_private_k0.pem"
	ztsPrivateKey0  = "testdata/zts_private_k0.pem"
	athenzConfigPath    = "testdata/athenz.json"
	zpeConfigPath       = "testdata/zpe.toml"
)

var testTempFolder string


func signPolicy(domainSignedPolicyData *zts.DomainSignedPolicyData,
	zmsPrivateKey string, ztsPrivateKey string) error{
	zmsData, err := ioutil.ReadFile(zmsPrivateKey)
	if err != nil {
		return common.Errorf("unable to open zms private key file, error: %s", err.Error())
	}

	signer, err := zmssvctoken.NewSigner(zmsData)
	if err != nil {
		return err
	}
	policyData, err := zpuUtil.ToCanonicalString(domainSignedPolicyData.SignedPolicyData.PolicyData)
	if err != nil {
		return err
	}
	signature, err := signer.Sign(policyData)
	if err != nil {
		return err
	}

	domainSignedPolicyData.SignedPolicyData.ZmsSignature = signature
	domainSignedPolicyData.SignedPolicyData.ZmsKeyId = "0"

	ztsData, err := ioutil.ReadFile(ztsPrivateKey)
	if err != nil {
		return common.Errorf("unable to open zts private key file, error: %s", err.Error())
	}

	signer, err = zmssvctoken.NewSigner(ztsData)
	if err != nil {
		return err
	}

	policyData, err = zpuUtil.ToCanonicalString(domainSignedPolicyData.SignedPolicyData)
	if err != nil {
		return err
	}

	signature, err = signer.Sign(policyData)
	if err != nil {
		return err
	}

	domainSignedPolicyData.Signature = signature
	domainSignedPolicyData.KeyId = "0"

	return nil
}

func preparePolicyFiles(expiry time.Time) error {

	log.NewLogrusInitializer().InitialLog(log.Info)

	if err := config.LoadGlobalZpeConfig(zpeConfigPath); err != nil {
		common.Fatalf("unable to load config, %s: ", err)
	}

	if err := config.LoadGlobalAthenzConfig(athenzConfigPath); err != nil {
		common.Fatalf("unable to load config, %s: ", err)
	}

	readFile, err := os.OpenFile(policyFile, os.O_RDONLY, 0444)
	if err != nil {
		return common.Errorf("cannot open file: %s , error: %s", policyFile, err.Error())
	}
	defer func() {
		cErr := readFile.Close()
		if cErr != nil {
			common.Fatal(err.Error())
		}
	}()

	var domainSignedPolicyData *zts.DomainSignedPolicyData

	if err := json.NewDecoder(readFile).Decode(&domainSignedPolicyData); err != nil {
		return common.Errorf("unable to decode policy file: %s, error: %s", policyFile, err.Error())
	}

	if expiry.UnixNano() > 0 {
		expiry = expiry.Add(48 * time.Hour)
		domainSignedPolicyData.SignedPolicyData.Expires = rdl.Timestamp{Time: expiry}
	}

	if err := signPolicy(domainSignedPolicyData,zmsPrivateKey0, ztsPrivateKey0); err !=nil {
		common.Fatal(err.Error())
	}

	testTempFolder, err = ioutil.TempDir(tmpDir, policyDirPrefix)
	if err != nil {
		return common.Errorf("unable to create policy directory: %s", err.Error())
	}

	data, _ := json.Marshal(domainSignedPolicyData)
	err = common.CreateFile(testTempFolder+"/angler.pol", string(data))
	if err != nil {
		return common.Error("unable to create policy file")
	}
	return nil
}

func createRoleToken(role, domain string) string {
	generatedToken := strconv.FormatInt((common.CurrentTimeMillis()/1000-30)*int64(time.Second), 10)
	expiration := strconv.FormatInt((common.CurrentTimeMillis()/1000+300)*int64(time.Second), 10)
	signedToken := "v=S1;d=" + domain + ";h=localhost" + ";r=" + role +
		";t=" + generatedToken + ";e=" + expiration + ";k=0"

	data, _ := ioutil.ReadFile(ztsPrivateKey0)

	signer, _ := zmssvctoken.NewSigner(data)
	signature, _ := signer.Sign(signedToken)

	signedToken = signedToken + ";s=" + signature

	return signedToken
}

func TestPermissionService_CheckAccessWithTokenPolicyFileExpired(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Time{})
	a.NoError(err)
	defer func() {
		err := common.RemoveAll(testTempFolder)
		if err != nil {
			common.Fatal(err.Error())
		}
	}()

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "read", Resource: "angler:stuff",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(9))
}

func TestPermissionService_CheckAccessWithTokenAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "read", Resource: "angler:stuff",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenDeny(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "throw", Resource: "angler:stuff",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(int32(1), status.AccessCheckStatus)

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenStartWith(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("public", "angler")

	request := &v1.AccessCheckRequest{Access: "fish", Resource: "angler:stockedpondBigBassLake",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenWildcardDeny(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("managerkernco", "angler")

	request := &v1.AccessCheckRequest{Access: "manage", Resource: "angler:pondsVenturaCounty",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(1))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenWildcardAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("managerkernco", "angler")

	request := &v1.AccessCheckRequest{Access: "manage", Resource: "angler:pondsKernCounty",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenMatchAllAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("matchall", "angler")

	request := &v1.AccessCheckRequest{Access: "all", Resource: "angler:anything",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenMatchRegexAllow(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("matchregex", "angler")

	request := &v1.AccessCheckRequest{Access: "regex", Resource: "angler:nhllllllkings",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow1(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:coretech",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(int32(0), status.AccessCheckStatus)

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow2(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:corecommit",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(status.AccessCheckStatus, int32(0))

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow3(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:borec",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(int32(0), status.AccessCheckStatus)

	_ = os.RemoveAll(testTempFolder)
}

func TestPermissionService_CheckAccessWithTokenFullRegexAllow4(t *testing.T) {
	a := assert.New(t)
	err := preparePolicyFiles(time.Now())
	a.NoError(err)

	files, _ := ioutil.ReadDir(testTempFolder)
	cache.PolicyDirectory = testTempFolder
	cache.LoadDB(files)

	signedToken := createRoleToken("full_regex", "angler")

	request := &v1.AccessCheckRequest{Access: "full_regex", Resource: "angler:cored",
		Token: signedToken}

	tst := PermissionService{}
	ctx := context.Background()
	status, err := tst.CheckAccessWithToken(ctx, request)
	a.NoError(err)
	a.Equal(int32(0), status.AccessCheckStatus)

	_ = os.RemoveAll(testTempFolder)
}
