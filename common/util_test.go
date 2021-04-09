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

package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
)

var (
	testFilename = GolangFileName()
)

func TestExists(t *testing.T) {
	a := assert.New(t)

	_, file, _, ok := runtime.Caller(0)
	a.True(ok)
	a.True(Exists(file))

	a.False(Exists("this/path/does/not/exists"))
}

func TestLoadFileStatusNull(t *testing.T) {

	a := assert.New(t)
	files, err := LoadFileStatus("./somewhere")
	a.Error(err)
	a.Nil(files)
}

func TestLoadFileStatusEmptyPath(t *testing.T) {
	a := assert.New(t)
	files, err := LoadFileStatus("")
	a.Error(err)
	a.True(strings.HasSuffix(err.Error(), "invalid input: directory name is cannot be empty"))
	a.Nil(files)
}

func TestLoadFileStatus(t *testing.T) {
	a := assert.New(t)
	// load current directory
	files, err := LoadFileStatus(".")
	a.NoError(err)
	a.Equal(7, len(files))
}

func TestStripDomainPrefix(t *testing.T) {
	a := assert.New(t)
	// domain exists in assertString
	a.Equal("books", StripDomainPrefix("domain:books", "domain", "some"))
	// domain doesn't exist in assertString
	a.Equal("books", StripDomainPrefix("books", "domain", "some"))
	// domain doesn't exist in assertString
	a.Equal("some", StripDomainPrefix("angler:books", "domain", "some"))
}

func TestCreateFile(t *testing.T) {
	a := assert.New(t)

	content := "You're a Yahoo Athenz fan, so this app is created for you. athenz-agent contains athenz ZPE and ZPU" +
		" utilities in Go language. ZPU will download the domains' policy files and store them into the filesystem." +
		" In other side, ZPE will use that policy files, and it will cache them into memory to use them as fast as possible."

	tmpDir, err := ioutil.TempDir("./", "tmp")
	a.NoError(err)
	filename := "test-file.txt"
	err = CreateFile(tmpDir+string(os.PathSeparator)+filename, content)
	a.NoError(err)

	err = CreateFile(tmpDir+string(os.PathSeparator)+filename, content)
	a.NoError(err)

	err = RemoveAll(tmpDir)
	a.NoError(err)
}

func TestVerifierPositiveTest(t *testing.T) {
	a := assert.New(t)
	publicKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTHpmU09UUUpmRW0xZW00TDNza3lOVlEvYngwTU9UcQphK1J3T0gzWmNNS3lvR3hPSm85QXllUmE2RlhNbXZKSkdZczVQMzRZc3pGcG5qMnVBYmkyNG5FQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"
	input := `{"expires":"2017-06-09T06:11:12.125Z","modified":"2017-06-02T06:11:12.125Z","policyData":{"domain":"sys.auth","policies":[{"assertions":[{"action":"*","effect":"ALLOW","resource":"*","role":"sys.auth:role.admin"},{"action":"*","effect":"DENY","resource":"*","role":"sys.auth:role.non-admin"}],"name":"sys.auth:policy.admin"}]},"zmsKeyId":"0","zmsSignature":"Y2HuXmgL86PL1WnleGFHwPmNEqUdWgDxmmIsDnF5f5oqakacqTtwt9JNqDV9nuJ7LnKl3zsZoDQSAtcHMu4IGA--"}`
	signature := "XJnQ4t33D4yr7NtUjLaWhXULFr76z.z0p3QV4uCkA5KR9L4liVRmICYwVmnXxvHAlImKlKLv7sbIHNsjBfGfCw--"
	key, err := new(zmssvctoken.YBase64).DecodeString(publicKey)
	if err != nil {
		Fatalf("Failed to decode key to Verify data , error: %s", err.Error())
	}
	err = Verify(input, signature, string(key))
	a.Nil(err, "Verifier failed for valid data")
}

func TestVerifierTamperedInput(t *testing.T) {
	a := assert.New(t)
	publicKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTHpmU09UUUpmRW0xZW00TDNza3lOVlEvYngwTU9UcQphK1J3T0gzWmNNS3lvR3hPSm85QXllUmE2RlhNbXZKSkdZczVQMzRZc3pGcG5qMnVBYmkyNG5FQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"
	input := `{"expires":"2017-06-09T06:11:12.125Z","modified" : "2017-06-02T06:11:12.125Z","policyData":{"domain":"sys.auth","policies":[{"assertions":[{"action":"*","effect":"ALLOW","resource":"*","role":"sys.auth:role.admin"},{"action":"*","effect":"DENY","resource":"*","role":"sys.auth:role.non-admin"}],"name":"sys.auth:policy.admin"}]},"zmsKeyId":"0","zmsSignature":"Y2HuXmgL86PL1WnleGFHwPmNEqUdWgDxmmIsDnF5f5oqakacqTtwt9JNqDV9nuJ7LnKl3zsZoDQSAtcHMu4IGA--"}`
	signature := "XJnQ4t33D4yr7NtUjLaWhXULFr76z.z0p3QV4uCkA5KR9L4liVRmICYwVmnXxvHAlImKlKLv7sbIHNsjBfGfCw--"
	key, err := new(zmssvctoken.YBase64).DecodeString(publicKey)
	a.NoError(err)
	err = Verify(input, signature, string(key))
	a.NotNil(err, "Verifier validated for invalid data")
}

func TestVerifierTamperedKey(t *testing.T) {
	a := assert.New(t)
	publicKey := "LS1tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTHpmU09UUUpmRW0xZW00TDNza3lOVlEvYngwTU9UcQphK1J3T0gzWmNNS3lvR3hPSm85QXllUmE2RlhNbXZKSkdZczVQMzRZc3pGcG5qMnVBYmkyNG5FQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"
	input := `{"expires":"2017-06-09T06:11:12.125Z","modified":"2017-06-02T06:11:12.125Z","policyData":{"domain":"sys.auth","policies":[{"assertions":[{"action":"*","effect":"ALLOW","resource":"*","role":"sys.auth:role.admin"},{"action":"*","effect":"DENY","resource":"*","role":"sys.auth:role.non-admin"}],"name":"sys.auth:policy.admin"}]},"zmsKeyId":"0","zmsSignature":"Y2HuXmgL86PL1WnleGFHwPmNEqUdWgDxmmIsDnF5f5oqakacqTtwt9JNqDV9nuJ7LnKl3zsZoDQSAtcHMu4IGA--"}`
	signature := "XJn4t33D4yr7NtUjLaWhXULFr76z.z0p3QV4uCkA5KR9L4liVRmICYwVmnXxvHAlImKlKLv7sbIHNsjBfGfCw--"
	key, err := new(zmssvctoken.YBase64).DecodeString(publicKey)
	a.NoError(err)
	err = Verify(input, signature, string(key))
	a.NotNil(err, "Verifier validated data with tampered key")
}

func TestVerifierTamperedSignature(t *testing.T) {
	a := assert.New(t)
	publicKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTHpmU09UUUpmRW0xZW00TDNza3lOVlEvYngwTU9UcQphK1J3T0gzWmNNS3lvR3hPSm85QXllUmE2RlhNbXZKSkdZczVQMzRZc3pGcG5qMnVBYmkyNG5FQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo-"
	input := `{"expires":"2017-06-09T06:11:12.125Z","modified":"2017-06-02T06:11:12.125Z","policyData":{"domain":"sys.auth","policies":[{"assertions":[{"action":"*","effect":"ALLOW","resource":"*","role":"sys.auth:role.admin"},{"action":"*","effect":"DENY","resource":"*","role":"sys.auth:role.non-admin"}],"name":"sys.auth:policy.admin"}]},"zmsKeyId":"0","zmsSignature":"Y2HuXmgL86PL1WnleGFHwPmNEqUdWgDxmmIsDnF5f5oqakacqTtwt9JNqDV9nuJ7LnKl3zsZoDQSAtcHMu4IGA--"}`
	signature := "XJpQ4t33D4yr7NtUjLaWhXULFr76z.z0p3QV4uCkA5KR9L4liVRmICYwVmnXxvHAlImKlKLv7sbIHNsjBfGfCw--"
	key, err := new(zmssvctoken.YBase64).DecodeString(publicKey)
	a.NoError(err)
	err = Verify(input, signature, string(key))
	a.NotNil(err, "verifier validated data with tampered signature")
}

func TestCreateAllDirectories(t *testing.T) {
	a := assert.New(t)
	path := "tmp/metric"
	err := CreateAllDirectories(path)
	a.NoError(err)
	_, err = os.Stat(path)
	a.False(os.IsNotExist(err))
	_ = os.RemoveAll("tmp")
}

func TestCreateAllDirectoriesExist(t *testing.T) {
	a := assert.New(t)
	path := "tmp/metric"
	err := CreateAllDirectories(path)
	a.NoError(err)
	_, err = os.Stat(path)
	a.False(os.IsNotExist(err))
	err = CreateAllDirectories(path)
	a.NoError(err)
	_ = os.RemoveAll("tmp")
}

func TestGetGolangFileName(t *testing.T) {
	a := assert.New(t)
	filename := GolangFileName()
	fmt.Println(filename)
	a.Equal("util_test.go", filename)
}

func TestGetGolangFileNamePackageLevel(t *testing.T) {
	a := assert.New(t)
	a.Equal("util_test.go", testFilename)
	fmt.Println(testFilename)
}

func TestFuncName(t *testing.T) {
	a := assert.New(t)
	funcName := FuncName()
	fmt.Println(funcName)
	a.Equal("common.TestFuncName", funcName)
}

func TestCallerFuncName(t *testing.T) {
	a := assert.New(t)
	funcName := callerName()
	fmt.Println(funcName)
	a.Equal("common.TestCallerFuncName", funcName)
}

func callerName() string {
	return CallerFuncName()
}
