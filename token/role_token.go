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
	"strconv"
	"strings"
	"time"
)

const (
	tagVersion               = "v"
	tagDomain                = "d"
	tagRoleNames             = "r"
	tagDomainCompleteRoleSet = "c"
	tagPrincipal             = "p"
	tagIP                    = "i"
	tagGenerationTime        = "t"
	tagExpireTime            = "e"
	tagSalt                  = "a"
	tagSignature             = "s"
	tagKeyId                 = "k"
	tagHostName              = "h"
)

var (
	logger = log.GetLogger(common.GolangFileName())
)

type RoleToken struct {
	Version               string   // the token version e.g. S1, U1
	Domain                string   // domain for which token is valid
	RoleNames             []string // list of comma separated roles
	DomainCompleteRoleSet bool     // the list of roles is complete in domain
	Principal             string   // principal that got the token was generated
	GenerationTime        int64    // time token was generated, nano second
	ExpiryTime            int64    // time token expires, nano second
	KeyId                 string   // identifier - either version or zone name
	Salt                  string   // a random 8 byte salt inner[1] hex encoded
	HostName              string   // host that issued this role token
	IPAddress             string   // ip address that issued this role token
	Signature             string   // signature generated over the roleToken string using Service's private Key and y64 encoded
	SignedToken           string   // roleToken in string format
	UnsignedToken         string   // roleToken with out signature to be validate
	AthenzTokenNoExpiry   bool     // roleToken can be expired for false or can live for ever for true
	AthenzTokenMaxExpiry  int64    // maximum lifetime of the roleToken
}

// validate roleToken by checking field like public key, signature and unsignedToken
// this fields must not be empty. checking generated time and expiry time and then
// verify the roleToken by checking public key and hashing of data and signature.
func (roleToken *RoleToken) Validate(publicKey string, allowedOffset int64, allowNoExpiry bool) (bool, error) {

	// check if data and signature exists
	if roleToken.UnsignedToken == "" || roleToken.Signature == "" {
		return false, common.Errorf("missing data/signature component, data: %s, signature: %s",
			roleToken.UnsignedToken, roleToken.Signature)
	}

	// check if public key exists
	if publicKey == "" {
		return false, common.Errorf("no public key provided, data: %s", roleToken.UnsignedToken)
	}

	now := common.CurrentTimeMillis() / 1000

	// make sure the token does not have a timestamp in the
	// future we'll allow the configured offset between servers.
	if roleToken.GenerationTime != 0 &&
		(roleToken.GenerationTime/int64(time.Second))-allowedOffset > now {
		return false, common.Errorf("token has future generatedTime, generated time: %+v, "+
			"now: %+v, allowed offset: %d", roleToken.GenerationTime, time.Unix(0, now), allowedOffset)
	}

	// make sure we don't have unlimited tokens unless we have
	// explicitly enabled that option for our system. by default
	// they should have an expiration date of less than 30 days.
	if roleToken.ExpiryTime != 0 || !allowNoExpiry {
		expiry := roleToken.ExpiryTime / int64(time.Second)
		if expiry < now {
			return false, common.Errorf("token has expired, expiry time: %+v, now: %+v",
				roleToken.ExpiryTime, time.Unix(0, now))
		}

		if expiry > now+(roleToken.AthenzTokenMaxExpiry*24*60*60)+allowedOffset {
			return false, common.Errorf("token expires too far in the future, expiryTime: %+v"+
				", current time: %+v, max expiry: %d days, allowed offset: %d", roleToken.ExpiryTime,
				time.Unix(0, now), roleToken.AthenzTokenMaxExpiry, allowedOffset)
		}

	}

	err := common.Verify(roleToken.UnsignedToken, roleToken.Signature, publicKey)
	if err != nil {
		logger.Error(err.Error())
		return false, nil
	}
	return true, nil
}

// create new roleToken by a roleToken string that created
// by zpe
func NewRoleToken(signedToken string) (*RoleToken, error) {
	if signedToken == "" {
		return nil, common.Error("input String signedToken must not be empty")
	}

	var err error
	var roleNames string
	roleToken := &RoleToken{}

	i := strings.Index(signedToken, ";s=")
	if i != -1 {
		roleToken.UnsignedToken = signedToken[0:i]
	}

	parts := strings.Split(signedToken, ";")

	for _, part := range parts {
		inner := strings.Split(part, "=")
		if len(inner) != 2 {
			return nil, common.Errorf("malformed token field %s",part)
		}

		switch inner[0] {
		case tagVersion:
			roleToken.Version = inner[1]
		case tagDomain:
			roleToken.Domain = inner[1]
		case tagRoleNames:
			roleNames = inner[1]
		case tagHostName:
			roleToken.HostName = inner[1]
		case tagIP:
			roleToken.IPAddress = inner[1]
		case tagKeyId:
			roleToken.KeyId = inner[1]
		case tagPrincipal:
			roleToken.Principal = inner[1]
		case tagSalt:
			roleToken.Salt = inner[1]
		case tagSignature:
			roleToken.Signature = inner[1]
		case tagDomainCompleteRoleSet:
			if i, err := strconv.Atoi(inner[1]); err != nil {
				return nil, common.Errorf("unable to extract roletoken, error: %s", err.Error())
			} else if i == 1 {
				roleToken.DomainCompleteRoleSet = true
			}
		case tagGenerationTime:
			if roleToken.GenerationTime, err = strconv.ParseInt(inner[1], 10, 64); err != nil {
				return nil, common.Errorf("cannot convert generation timestamp to int64")
			}
		case tagExpireTime:
			if roleToken.ExpiryTime, err = strconv.ParseInt(inner[1], 10, 64); err != nil {
				return nil, common.Errorf("cannot convert expiry timestamp to int64")
			}
		default:
			logger.Info("Unknown ntoken field: " + inner[1])
		}
	}

	roleToken.AthenzTokenNoExpiry = config.ZpeConfig.Properties.AthenzTokenNoExpiry
	// convert days to milliseconds
	roleToken.AthenzTokenMaxExpiry = config.ZpeConfig.Properties.AthenzTokenMaxExpiry

	// the required attributes for the token are domain
	// and roles. The signature will be verified during
	// the authenticate phase but now we'll make sure
	// that domain and roles are present
	if roleToken.Domain == "" {
		return nil, common.Error("signedToken does not contain required domain component")
	}

	if roleNames == "" {
		return nil, common.Error("signedToken does not contain required roles component")
	}

	roleToken.RoleNames = strings.Split(roleNames, ",")
	roleToken.SignedToken = signedToken

	return roleToken, nil
}

func asTime(s, name string) (time.Time, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, common.Errorf("invalid field inner[1] '%s' for field '%s'", s, name)
	}
	return time.Unix(n, 0), nil
}
