/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/12/19
 * Time: 11:04 AM
 *
 * Description:
 * In here we describe roleToken. Also you can create a
 * role token by a signed token string. also you can validate
 * roleToken.
 *
 */

package token

import (
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/hamed-yousefi/athenz-agent/common/log"
	"github.com/hamed-yousefi/athenz-agent/config"
	"github.com/stretchr/testify/assert"
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	svcVersion      = "S1"
	svcDomain       = "sports"
	role1           = "storage.tenant.weather.updater"
	role2           = "fantasy.tenant.sports.admin"
	role3           = "fantasy.tenant.sports.reader"
	role4           = "fantasy.tenant.sports.writer"
	role5           = "fantasy.tenant.sports.scanner"
	host            = "somehost.somecompany.com"
	salt            = "saltstring"
	testPrivateKey0 = "testdata/private_key0.key"
	testPublicKey0  = "testdata/public_key0.key"
	configPath      = "testdata/zpe.toml"
)

func setup() {
	log.NewLogrusInitializer().InitialLog(log.Info)

	if err := config.LoadGlobalZpeConfig(configPath); err != nil {
		common.Fatalf("unable to load config, %s: ", err)
	}
}

func TestNewRoleToken(t *testing.T) {
	setup()

	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=trialblaze;r=role1,role2;s=signature")
	a.NoError(err)
	a.NotNil(roleToken)
	a.Equal(roleToken.Version, "S1")
	a.Equal(roleToken.Domain, "trialblaze")
	a.Len(roleToken.RoleNames, 2)
	a.Equal(roleToken.RoleNames[0], "role1")
	a.Equal(roleToken.RoleNames[1], "role2")
	a.Equal(roleToken.Signature, "signature")
}

func TestNewRoleTokenEmpty(t *testing.T) {
	setup()

	a := assert.New(t)

	roleToken, err := NewRoleToken("")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal("token.NewRoleToken-> input String signedToken must not be empty", err.Error())
}

func TestNewRoleTokenWithoutDomain(t *testing.T) {
	setup()
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;r=role1,role2;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal("token.NewRoleToken-> signedToken does not contain required domain component", err.Error())
}

func TestNewRoleTokenEmptyDomain(t *testing.T) {
	setup()
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=;r=role1,role2;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal("token.NewRoleToken-> signedToken does not contain required domain component", err.Error())
}

func TestNewRoleTokenWithoutRole(t *testing.T) {
	setup()
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=trialblaze;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal("token.NewRoleToken-> signedToken does not contain required roles component", err.Error())
}

func TestNewRoleTokenEmptyRole(t *testing.T) {
	setup()
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=trialblaze;r=;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal("token.NewRoleToken-> signedToken does not contain required roles component", err.Error())
}

func TestNewRoleTokenInvalidVersion(t *testing.T) {
	setup()
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1=S2;d=trialblaze;r=role1,role2;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal("token.NewRoleToken-> malformed token field v=S1=S2", err.Error())
}

func TestTimeConversion(t *testing.T) {
	setup()
	a := assert.New(t)

	unix := "1549981415"
	wrongValue := "j2je2k3e23"
	_, err := asTime(unix, tagGenerationTime)
	a.NoError(err)

	_, err = asTime(wrongValue, tagGenerationTime)
	a.NotNil(err)
}

func TestValidateNilSignature(t *testing.T) {
	setup()
	a := assert.New(t)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt + ";h=" + host + ";r=" + role1
	var pubKey string

	roleToken, err := NewRoleToken(signedToken)
	a.NoError(err)
	a.NotNil(roleToken)

	isValid, err := roleToken.Validate(pubKey, 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.Equal(err.Error(), "token.(*RoleToken).Validate-> missing data/signature component, data: , signature: ")
}

func TestValidateNilPublicKey(t *testing.T) {
	setup()
	a := assert.New(t)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt + ";h=" + host + ";r=" + role1 + ";s=somesignature"
	var pubKey string

	roleToken, err := NewRoleToken(signedToken)
	a.NoError(err)
	a.NotNil(roleToken)

	isValid, err := roleToken.Validate(pubKey, 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "token.(*RoleToken).Validate-> no public key provided"))
}

func TestValidateFutureTimeStamp(t *testing.T) {
	setup()
	a := assert.New(t)
	generatedToken := strconv.FormatInt((common.CurrentTimeMillis()/1000+4600)*int64(time.Second), 10)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";t=" + generatedToken + ";s=" + "somesignature"

	roleToken, err := NewRoleToken(signedToken)
	a.NoError(err)
	a.NotNil(roleToken)

	isValid, err := roleToken.Validate("someInvalidPubKey", 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "token.(*RoleToken).Validate-> token has future generatedTime, generated time:"))

}

func TestValidateNoExpiry(t *testing.T) {
	setup()
	a := assert.New(t)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";s=" + "somesignature"

	roleToken, err := NewRoleToken(signedToken)
	a.NoError(err)
	a.NotNil(roleToken)
	roleToken.AthenzTokenNoExpiry = false

	isValid, err := roleToken.Validate("someInvalidPubKey", 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "token.(*RoleToken).Validate-> token has expired"))
}

func TestValidateTooFarExpiryTimestamp(t *testing.T) {
	setup()
	invalidPublicKey := "-----BEGIN PUBLIC KEY-----\nsomeInvalidPubKey\n-----END PUBLIC KEY-----"
	a := assert.New(t)
	expiration := strconv.FormatInt((common.CurrentTimeMillis()/1000+(30*24*60*60)+10)*int64(time.Second), 10)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";e=" + expiration + ";s=" + "somesignature"

	roleToken, err := NewRoleToken(signedToken)
	a.NoError(err)
	a.NotNil(roleToken)

	roleToken.AthenzTokenMaxExpiry = 30
	isValid, err := roleToken.Validate(invalidPublicKey, 5, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "token.(*RoleToken).Validate-> token expires too far in the future"))

	isValid, _ = roleToken.Validate(invalidPublicKey, 20, false)
	a.False(isValid)
}

func TestValidate(t *testing.T) {
	setup()
	a := assert.New(t)
	generatedToken := strconv.FormatInt((common.CurrentTimeMillis()/1000-30)*int64(time.Second), 10)
	expiration := strconv.FormatInt((common.CurrentTimeMillis()/1000+30)*int64(time.Second), 10)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";t=" + generatedToken + ";e=" + expiration

	data, err := ioutil.ReadFile(testPrivateKey0)
	a.NoError(err)
	a.NotNil(data)

	signer, err := zmssvctoken.NewSigner(data)
	a.NoError(err)
	a.NotNil(signer)

	signature, err := signer.Sign(signedToken)
	a.NoError(err)
	a.NotNil(signature)

	signedToken = signedToken + ";s=" + signature
	roleToken, err := NewRoleToken(signedToken)
	a.NoError(err)
	a.NotNil(roleToken)

	pubKey, err := ioutil.ReadFile(testPublicKey0)
	a.NoError(err)
	a.NotNil(pubKey)

	isValid, err := roleToken.Validate(string(pubKey), 3600, false)
	a.NoError(err)
	a.True(isValid)
}
