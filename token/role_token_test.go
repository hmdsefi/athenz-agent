/**
 * Copyright (c) 2019 TRIALBLAZE PTY. LTD. All rights reserved.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/12/19
 * Time: 4:47 PM
 *
 * Description:
 *
 */

package token

import (
	"github.com/stretchr/testify/assert"
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	"gitlab.com/trialblaze/athenz-agent/common/util"
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
	testPrivateKey0 = "../resource/test_data/private_key0.key"
	testPublicKey0  = "../resource/test_data/public_key0.key"
)

func TestNewRoleToken(t *testing.T) {
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=trialblaze;r=role1,role2;s=signature")
	a.Nil(err)
	a.NotNil(roleToken)
	a.Equal(roleToken.Version, "S1")
	a.Equal(roleToken.Domain, "trialblaze")
	a.Len(roleToken.RoleNames, 2)
	a.Equal(roleToken.RoleNames[0], "role1")
	a.Equal(roleToken.RoleNames[1], "role2")
	a.Equal(roleToken.Signature, "signature")
}

func TestNewRoleTokenEmpty(t *testing.T) {
	a := assert.New(t)

	roleToken, err := NewRoleToken("")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal(err.Error(), "NewRoleToken: input String signedToken must not be empty")
}

func TestNewRoleTokenWithoutDomain(t *testing.T) {
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;r=role1,role2;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal(err.Error(), "NewRoleToken: signedToken does not contain required domain component")
}

func TestNewRoleTokenEmptyDomain(t *testing.T) {
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=;r=role1,role2;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal(err.Error(), "NewRoleToken: signedToken does not contain required domain component")
}

func TestNewRoleTokenWithoutRole(t *testing.T) {
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=trialblaze;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal(err.Error(), "NewRoleToken: signedToken does not contain required roles component")
}

func TestNewRoleTokenEmptyRole(t *testing.T) {
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1;d=trialblaze;r=;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal(err.Error(), "NewRoleToken: signedToken does not contain required roles component")
}

func TestNewRoleTokenInvalidVersion(t *testing.T) {
	a := assert.New(t)

	roleToken, err := NewRoleToken("v=S1=S2;d=trialblaze;r=role1,role2;s=signature")
	a.NotNil(err)
	a.Nil(roleToken)
	a.Equal(err.Error(), "NewRoleToken: malformed token field v=S1=S2")
}

func TestTimeConversion(t *testing.T) {
	a := assert.New(t)

	unix := "1549981415"
	wrongValue := "j2je2k3e23"
	_, err := asTime(unix, tagGenerationTime)
	a.Nil(err)

	_, err = asTime(wrongValue, tagGenerationTime)
	a.NotNil(err)
}

func TestValidateNilSignature(t *testing.T) {

	a := assert.New(t)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt + ";h=" + host + ";r=" + role1
	var pubKey string

	roleToken, err := NewRoleToken(signedToken)
	a.Nil(err)
	a.NotNil(roleToken)

	isValid, err := roleToken.Validate(pubKey, 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.Equal(err.Error(), "RoleToken:Validate: missing data/signature component, data: , signature: ")
}

func TestValidateNilPublicKey(t *testing.T) {
	a := assert.New(t)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt + ";h=" + host + ";r=" + role1 + ";s=somesignature"
	var pubKey string

	roleToken, err := NewRoleToken(signedToken)
	a.Nil(err)
	a.NotNil(roleToken)

	isValid, err := roleToken.Validate(pubKey, 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "RoleToken:Validate: no public key provided"))
}

func TestValidateFutureTimeStamp(t *testing.T) {

	a := assert.New(t)
	generatedToken := strconv.FormatInt((util.CurrentTimeMillis()/1000+4600)*int64(time.Second), 10)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";t=" + generatedToken + ";s=" + "somesignature"

	roleToken, err := NewRoleToken(signedToken)
	a.Nil(err)
	a.NotNil(roleToken)

	isValid, err := roleToken.Validate("someInvalidPubKey", 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "RoleToken:Validate: token has future generatedTime, generated time:"))

}

func TestValidateNoExpiry(t *testing.T) {

	a := assert.New(t)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";s=" + "somesignature"

	roleToken, err := NewRoleToken(signedToken)
	a.Nil(err)
	a.NotNil(roleToken)
	roleToken.AthenzTokenNoExpiry = false

	isValid, err := roleToken.Validate("someInvalidPubKey", 3600, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "RoleToken:Validate: token has expired"))
}

func TestValidateTooFarExpiryTimestamp(t *testing.T) {
	a := assert.New(t)
	expiration := strconv.FormatInt((util.CurrentTimeMillis()/1000+(30*24*60*60)+10)*int64(time.Second), 10)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";e=" + expiration + ";s=" + "somesignature"

	roleToken, err := NewRoleToken(signedToken)
	a.Nil(err)
	a.NotNil(roleToken)

	roleToken.AthenzTokenMaxExpiry = 30
	isValid, err := roleToken.Validate("someInvalidPubKey", 5, false)
	a.NotNil(err)
	a.False(isValid)
	a.True(strings.HasPrefix(err.Error(), "RoleToken:Validate: token expires too far int the future"))

	isValid, err = roleToken.Validate("someInvalidPubKey", 20, false)
	a.Nil(err)
	a.False(isValid)
}

func TestValidate(t *testing.T) {

	a := assert.New(t)
	generatedToken := strconv.FormatInt((util.CurrentTimeMillis()/1000-30)*int64(time.Second), 10)
	expiration := strconv.FormatInt((util.CurrentTimeMillis()/1000+30)*int64(time.Second), 10)
	signedToken := "v=" + svcVersion + ";d=" + svcDomain + ";a=" + salt +
		";h=" + host + ";r=" + role1 + ";t=" + generatedToken + ";e=" + expiration

	data, err := ioutil.ReadFile(testPrivateKey0)
	a.Nil(err)
	a.NotNil(data)

	signer, err := zmssvctoken.NewSigner(data)
	a.Nil(err)
	a.NotNil(signer)

	signature, err := signer.Sign(signedToken)
	a.Nil(err)
	a.NotNil(signature)

	signedToken = signedToken + ";s=" + signature
	roleToken, err := NewRoleToken(signedToken)
	a.Nil(err)
	a.NotNil(roleToken)

	pubKey, err := ioutil.ReadFile(testPublicKey0)
	a.Nil(err)
	a.NotNil(pubKey)

	isValid, err := roleToken.Validate(string(pubKey), 3600, false)
	a.Nil(err)
	a.True(isValid)
}
