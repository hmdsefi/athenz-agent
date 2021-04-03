/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/14/19
 * Time: 4:17 PM
 *
 * Description:
 * This file contains PermissionService struct that implements
 * gRPC PermissionServer interface. You can use PermissionService
 * to create gRPC server. After that gRPC server starts, you can
 * call CheckAccessWithToken and GetServiceToken with your gRPC
 * client. CheckAccessWithToken will be used for checking an access
 * to a specific resource by a roleToken and GetServiceToken generates
 * roleToken.
 *
 */

package auth

import (
	"crypto/tls"
	"github.com/yahoo/athenz/clients/go/zts"
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	"gitlab.com/trialblaze/athenz-agent/cache"
	"gitlab.com/trialblaze/athenz-agent/common"
	"gitlab.com/trialblaze/athenz-agent/config"
	"gitlab.com/trialblaze/athenz-agent/matcher"
	"gitlab.com/trialblaze/athenz-agent/token"
	"gitlab.com/trialblaze/grpc-go/pkg/api/common/message/v1"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// Constant values that will return by
// CheckAccessWithToken method
const (
	Allow                 = 0
	Deny                  = 1
	DenyRoleTokenExpired  = 2
	DenyRoleTokenInvalid  = 3
	DenyInvalidParameters = 4
	DenyDomainMismatch    = 5
	DenyDomainNotFound    = 6
	DenyNoMatch           = 7
	DenyDomainEmpty       = 8
	DenyDomainExpired     = 9
)

// We will implement gRPC PermissionServer
// interface for this struct to use it in
// gRPC server.
// This interface has two method:
// 		* CheckAccessWithToken
//      * GetServiceToken
type PermissionService struct{}

// This method implements one of PermissionServer
// interface. CheckAccessWithToken accept a struct
// named AccessCheckRequest that contains roleToke,
// access and resource that roleToken wants to use.
// This method will return a AccessCheckResponse
// type that contains an access number between 0
// and 9.
func (permService PermissionService) CheckAccessWithToken(ctx context.Context,
	req *v1.AccessCheckRequest) (*v1.AccessCheckResponse, error) {

	// first try to get RoleToken from
	// cached RoleTokens
	roleToken, ok := cache.RoleTokenCacheMap[req.Token]
	if !ok {
		// this is first time that we trying to create
		// this rToken, so we will cache it after
		// validation step
		rToken, err := token.NewRoleToken(req.Token)
		if err != nil {
			return nil, common.Errorf("unable to create RoleToken, error: %s", err.Error())
		}

		// validate the rToken
		pubKey := config.KeyStore.GetZtsPublicKey(rToken.KeyId)
		ztsKey, err := new(zmssvctoken.YBase64).DecodeString(pubKey)
		isValid, err := rToken.Validate(string(ztsKey), config.ZpeConfig.Properties.AllowedOffset, false)
		if err != nil {
			return nil, common.Errorf("validation failed, error: %s", err.Error())
		}
		roleToken = rToken

		// rToken is not valid, so check expiry first
		if !isValid {
			// check the rToken expiration
			now := common.CurrentTimeMillis()
			if rToken.ExpiryTime != 0 && (rToken.ExpiryTime/int64(time.Millisecond)) < now {
				return &v1.AccessCheckResponse{AccessCheckStatus: DenyRoleTokenExpired}, nil
			}

			return &v1.AccessCheckResponse{AccessCheckStatus: DenyRoleTokenInvalid}, nil
		}

		cache.RoleTokenCacheMap[req.Token] = rToken
	} else {
		// check the cached token expiration
		// if it was expired remove it from
		// cached tokens
		now := common.CurrentTimeMillis()
		if roleToken.ExpiryTime != 0 && (roleToken.ExpiryTime/int64(time.Millisecond)) < now {
			delete(cache.RoleTokenCacheMap, req.Token)
			return &v1.AccessCheckResponse{AccessCheckStatus: DenyRoleTokenExpired}, nil
		}
	}

	return allowAction(req.Access, req.Resource, roleToken.Domain, roleToken.RoleNames)
}

// This method implements one of PermissionServer.
// GetServiceToken accept a struct
// named ServiceTokenRequest that is a empty struct
// used in gRPC request message. This method will
// return ServiceTokenResponse type that contains
// a token string.
// There are three ways to getting roleToken from
// ZTS server:
//		* Using athenz service identity certificate,
//		  it means that we can use our service private
//		  key and cert file
//		* Using ntoken from a file, this ntoken will
//		  be expired in some time periods
//		* Using ntoken as command-line (not recommended
//		  since others running ps might see your ntoken)
// we will use athenz service identity certificate
// in here to get roleToken from ZTS server. we're
// using copper argos which only uses tls and the
// attestation data contains the authentication details
func (permService PermissionService) GetServiceToken(ctx context.Context,
	req *v1.ServiceTokenRequest) (*v1.ServiceTokenResponse, error) {

	tlsConfig, err := getTLSConfigFromFiles(config.ZpeConfig.Properties.KeyFilePath, config.ZpeConfig.Properties.CertFilePath)
	if err != nil {
		return nil, common.Errorf("unable to load TLS Config, error: %s", err.Error())
	}

	minExpiryTime := config.ZpeConfig.Properties.TokenExpirationMin * 60
	maxExpiryTime := config.ZpeConfig.Properties.TokenExpirationMax * 60

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := zts.NewClient(config.KeyStore.Properties.ZtsUrl, transport)

	roleToken, err := client.GetRoleToken(zts.DomainName(config.ZpeConfig.Properties.DomainName),
		zts.EntityList(config.ZpeConfig.Properties.RoleNames), &minExpiryTime, &maxExpiryTime, "")
	if err != nil {
		return nil, common.Errorf("unable to get roleToken, error: %s", err.Error())
	}

	return &v1.ServiceTokenResponse{Token: roleToken.Token}, nil
}

func allowAction(action, resource, domain string, roles []string) (*v1.AccessCheckResponse, error) {

	// check parameters to not be empty
	if roles == nil || len(roles) == 0 {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyRoleTokenInvalid}, nil
	}

	if domain == "" {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyRoleTokenInvalid}, nil
	}

	if action == "" {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyInvalidParameters}, nil
	}

	if resource == "" {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyInvalidParameters}, nil
	}

	action = strings.ToLower(action)
	resource = strings.ToLower(resource)
	resource = common.StripDomainPrefix(resource, domain, "")

	// Note: if domain in token doesn't match
	// domain in resource then there will be
	// no match of any resource in the assertions
	// - so deny immediately
	if resource == "" {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyDomainMismatch}, nil
	}

	now := time.Now().UnixNano()

	var status int32
	status = DenyDomainNotFound

	// first hunt by role for deny assertions since
	// deny takes precedence over allow assertions
	roleMap, ok := cache.DomainStandardRoleDenyMap[domain]
	if ok && roleMap.Expiry < now {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyDomainExpired}, nil
	}
	if ok && len(roleMap.RoleDataMap) > 0 {
		if actionByRole(action, resource, roles, roleMap.RoleDataMap) {
			return &v1.AccessCheckResponse{AccessCheckStatus: Deny}, nil
		} else {
			status = DenyNoMatch
		}
	} else if ok {
		status = DenyDomainEmpty
	}

	// if the check was not explicitly denied by a
	// standard role, then let's process our wildcard
	// roles for deny assertions
	roleMap, ok = cache.DomainWildcardRoleDenyMap[domain]
	if ok && roleMap.Expiry < now {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyDomainExpired}, nil
	}
	if ok && len(roleMap.RoleDataMap) > 0 {
		if actionByWildCardRole(action, resource, roles, roleMap.RoleDataMap) {
			return &v1.AccessCheckResponse{AccessCheckStatus: Deny}, nil
		} else {
			status = DenyNoMatch
		}
	} else if ok {
		status = DenyDomainEmpty
	}

	// so far it did not match any deny assertions so now let's
	// process our allow assertions
	roleMap, ok = cache.DomainStandardRoleAllowMap[domain]
	if ok && roleMap.Expiry < now {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyDomainExpired}, nil
	}
	if ok && len(roleMap.RoleDataMap) > 0 {
		if actionByRole(action, resource, roles, roleMap.RoleDataMap) {
			return &v1.AccessCheckResponse{AccessCheckStatus: Allow}, nil
		} else {
			status = DenyNoMatch
		}
	} else if ok {
		status = DenyDomainEmpty
	}

	// at this point we either got an allow or didn't match anything so we're
	// going to try the wildcard roles
	roleMap, ok = cache.DomainWildcardRoleAllowMap[domain]
	if ok && roleMap.Expiry < now {
		return &v1.AccessCheckResponse{AccessCheckStatus: DenyDomainExpired}, nil
	}
	if ok && len(roleMap.RoleDataMap) > 0 {
		if actionByWildCardRole(action, resource, roles, roleMap.RoleDataMap) {
			return &v1.AccessCheckResponse{AccessCheckStatus: Allow}, nil
		} else {
			status = DenyNoMatch
		}
	} else if ok {
		status = DenyDomainEmpty
	}

	return &v1.AccessCheckResponse{AccessCheckStatus: status}, nil
}

func actionByRole(action, resource string, roles []string,
	roleMap map[string][]map[string]interface{}) bool {

	var asserts []map[string]interface{}
	var ok bool
	for _, role := range roles {
		asserts, ok = roleMap[role]
		if !ok {
			continue
		}

		// see if any of its assertions match the action and resource
		// the assert action value does not have the domain prefix
		// ex: "Modify"
		// the assert resource value has the domain prefix
		// ex: "angler:angler.stuff"
		if matchAssertions(asserts, action, resource) {
			return true
		}
	}
	return false
}

func matchAssertions(asserts []map[string]interface{}, action, resource string) bool {

	var match matcher.ZpeMatch
	for _, strAssert := range asserts {

		// ex: "mod*"
		match = reflect.ValueOf(strAssert[common.ZpeActionMatchStruct]).Interface().(matcher.ZpeMatch)
		if !match.Match(action) {
			continue
		}

		// ex: "weather:service.storage.tenant.sports.*"
		match = reflect.ValueOf(strAssert[common.ZpeResourceMatchStruct]).Interface().(matcher.ZpeMatch)
		if !match.Match(resource) {
			continue
		}

		return true
	}

	return false
}

func actionByWildCardRole(action, resource string, roles []string,
	roleMap map[string][]map[string]interface{}) bool {

	// find policy matching resource and action
	// get assertions for given domain+role
	// then cycle thru those assertions looking
	// for matching action and resource.
	// we will visit each of the wildcard roles
	keys := make([]string, 0, len(roleMap))
	for key := range roleMap {
		keys = append(keys, key)
	}

	var asserts []map[string]interface{}
	var assert map[string]interface{}
	var match matcher.ZpeMatch
	var ok bool

	for _, role := range roles {
		for _, roleName := range keys {

			asserts, ok = roleMap[roleName]
			if !ok {
				continue
			}

			assert = asserts[0]
			match = reflect.ValueOf(assert[common.ZpeRoleMatchStruct]).Interface().(matcher.ZpeMatch)
			if !match.Match(role) {
				continue
			}

			// HAVE: matched the role with the wildcard

			// see if any of its assertions match the action and resource
			// the assert action value does not have the domain prefix
			// ex: "Modify"
			// the assert resource value has the domain prefix
			// ex: "angler:angler.stuff"
			if matchAssertions(asserts, action, resource) {
				return true
			}
		}
	}

	return false
}

// accept keyFile and certFile address and read content
// in byte array format and pass them getTLSConfig to
// create tls config
func getTLSConfigFromFiles(keyFilePath, certFilePath string) (*tls.Config, error) {
	keyPem, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, common.Errorf("unable to read keyFile, error: %s", err.Error())
	}

	certPem, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return nil, common.Errorf("unable to read certFile, error: %s", err.Error())
	}

	return getTLSConfig(keyPem, certPem)
}

// use key and cert to create tls config
func getTLSConfig(keyPem, certPem []byte) (*tls.Config, error) {
	clientCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, common.Errorf("unable to formulate clientCert "+
			"from keyPem and certPem, error: %s", err.Error())
	}

	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0] = clientCert

	return tlsConfig, nil
}
