/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/1/21
 * Time: 9:57 PM
 *
 * Description:
 *
 */

package log

import (
	"github.com/hamed-yousefi/athenz-agent/common"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	logPath = "logs/app"
)

type (
	logrusConfigProviderTest struct {
		level           string
		path            string
		maxAge          time.Duration
		rotationTime    time.Duration
		maxSize         int64
		filenamePattern string
	}
)

func newLogConfigProviderTest(path string) common.LogConfigProvider {
	return logrusConfigProviderTest{
		level:           Info.String(),
		path:            path,
		maxAge:          time.Duration(600) * time.Second,
		rotationTime:    time.Duration(1) * time.Hour,
		maxSize:         167772160,
		filenamePattern: ".%Y-%m-%dT%H:%M",
	}
}

func (l logrusConfigProviderTest) GetLevel() string {
	return l.level
}

func (l logrusConfigProviderTest) GetPath() string {
	return l.path
}

func (l logrusConfigProviderTest) GetMaxAge() time.Duration {
	return l.maxAge
}

func (l logrusConfigProviderTest) GetRotationTime() time.Duration {
	return l.rotationTime
}

func (l logrusConfigProviderTest) GetMaxSize() int64 {
	return l.maxSize
}

func (l logrusConfigProviderTest) GetFilenamePattern() string {
	return l.filenamePattern
}

func tearDown() {
	path := strings.Split(logPath, string(os.PathSeparator))
	err := common.RemoveAll(path[0])
	if err != nil {
		common.Fatal(err.Error())
	}
}

func TestLogrusLogRotator_SetupRotation(t *testing.T) {
	a := assert.New(t)
	logInit := NewLogrusInitializer()
	rotator := logInit.InitialLog(Info)
	rotator.SetupRotation(newLogConfigProviderTest(logPath))
	_, err := os.Stat(logPath)
	a.NoError(err)
	logger := GetLogger(common.GolangFileName())
	logger.Info("Hello World!")
	tearDown()
}

func TestLevels(t *testing.T) {
	logInit := NewLogrusInitializer()
	logInit.InitialLog(Trace)
	logger := GetLogger(common.GolangFileName())
	logger.Info("test info")
	logger.Error("test error")
	logger.Debug("test debug")
	logger.Trace("test trace")
}
