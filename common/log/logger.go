/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Date: 3/31/21
 * Time: 8:15 PM
 *
 * Description:
 *
 */

package log

import (
	"github.com/hamed-yousefi/athenz-agent/common"
)

const (
	// Fatal level.  Logs and then calls `log.Exit(1)`.
	Fatal Level = 1
	// Error level. Used for errors that should definitely be noted.
	Error Level = 2
	// Info level. General operational entries about what's going on inside the application.
	Info Level = 3
	// Debug level. Usually only enabled when debugging. Very verbose logging.
	Debug Level = 4
	// Trace level. Designates finer-grained informational events than the Debug.
	Trace Level = 5
)

var (
	string2Level = map[string]Level{
		"fatal": Fatal,
		"error": Error,
		"info":  Info,
		"debug": Debug,
		"trace": Trace,
	}

	level2String = map[Level]string{
		Fatal: "fatal",
		Error: "error",
		Info:  "info",
		Debug: "debug",
		Trace: "trace",
	}
)

type (
	Level uint32

	// Logger is a general interface for logging.
	Logger interface {
		Fatal(msg string)
		Fatalf(format string, params ...interface{})
		Info(msg string)
		Error(msg string)
		Debug(msg string)
		Trace(msg string)
	}

	// Initializer the interface that wrap log init function.
	Initializer interface {
		// InitialLog creates a log object internally and
		// returns a log rotator object for optional extra configuration.
		InitialLog(level Level) Rotator
	}

	// Rotator the interface to wrap log rotation config.
	Rotator interface {
		SetupRotation(provider common.LogConfigProvider)
	}
)

func (l Level) String() string {
	str, ok := level2String[l]
	if !ok {
		return ""
	}
	return str
}

func GetLevel(in string) Level {
	level, ok := string2Level[in]
	if !ok {
		common.Fatalf("invalid input, level: %s", in)
	}
	return level
}
