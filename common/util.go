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
 * Utility file, Use this method every where you want.
 *
 */

package common

import (
	"github.com/yahoo/athenz/libs/go/zmssvctoken"
	"github.com/yahoo/athenz/utils/zpe-updater/util"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	unknown = "unknown"
)

// check if file path is exists or not
func Exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

// return list of files from given directory
func LoadFileStatus(dirName string) ([]os.FileInfo, error) {
	if len(dirName) <= 0 || dirName == "" {

	}

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// this function return a timestamp by unix millisecond
func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Verify input json with zms or zts public keys
func Verify(input, signature, publicKey string) error {
	verifier, err := zmssvctoken.NewVerifier([]byte(publicKey))
	if err != nil {
		return err
	}

	err = verifier.Verify(input, signature)
	return err
}

func StripDomainPrefix(assertString, domain, defaultValue string) string {
	index := strings.Index(assertString, ":")
	if index == -1 {
		return assertString
	}
	if assertString[0:index] != domain {
		return defaultValue
	}
	return assertString[index+1:]
}

func CreateFile(fileName, content string) error {
	if util.Exists(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			return Errorf("unable to remove file: %v, Error: %s", fileName, err.Error())
		}
	}

	err := ioutil.WriteFile(fileName, []byte(content), 0755)
	if err != nil {
		return Errorf("unable to write file: %v, Error: %s", fileName, err.Error())
	}

	return nil
}

func RemoveAll(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Println(err)
	}
}

// CreateMetricDirectory makes new directory for metric file, if it doesn't exist
// CreateAllDirectories makes directory with all sub directories
func CreateAllDirectories(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return Error("CreateAllDirectories: cannot create path, error: " + err.Error())
		}
	}
	return nil
}

// GolangFileName returns the caller function file name.
// The return string is empty if it was not possible to recover information.
//
// File name is useful in logs and error messages. GolangFileName helps to
// get the file name in runtime and prevents using hardcoded strings in code.
func GolangFileName() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	path := strings.Split(filePath, string(os.PathSeparator))
	return path[len(path)-1]
}

// CallerFuncName returns caller function name. If there's no caller function then
// returns 'unknown'.
//
// Caller name is important in logs and errors messages. CallerFuncName helps to
// get the caller function name in runtime and prevents hardcoded string in code.
func CallerFuncName() string {
	return getFuncName(4)
}

// FuncName returns function name. It returns 'unknown' if there's no function.
//
// Function name is important in logs and errors messages. FuncName helps to
// get the function name in runtime and prevents hardcoded string in code.
func FuncName() string {
	return getFuncName(3)
}

//
func getFuncName(skip int) string {
	pc := make([]uintptr, 1)

	n := runtime.Callers(skip, pc)
	if n == 0 {
		return unknown
	}

	caller := runtime.FuncForPC(pc[0]-1)
	if caller == nil {
		return unknown
	}

	path := strings.Split(caller.Name(), string(os.PathSeparator))
	return path[len(path)-1]
}
