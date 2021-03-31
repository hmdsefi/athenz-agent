/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Date: 3/31/21
 * Time: 8:15 PM
 *
 * Description:
 *
 */

package log

import (
	"errors"
)

const (
	// Fatal level.  Logs and then calls `logger.Exit(1)`.
	Fatal = 1
	// Error level. Used for errors that should definitely be noted.
	Error = 2
	// Info level. General operational entries about what's going on inside the application.
	Info = 3
	// Debug level. Usually only enabled when debugging. Very verbose logging.
	Debug = 4
	// Trace level. Designates finer-grained informational events than the Debug.
	Trace = 5
)

var (
	string2Level = map[string]Level{
		"fatal": 1,
		"error": 2,
		"info":  3,
		"debug": 4,
		"trace": 5,
	}

	level2String = map[Level]string{
		1: "fatal",
		2: "error",
		3: "info",
		4: "debug",
		5: "trace",
	}
)

type (
	Level uint32

	// Logger is a general interface for logging.
	Logger interface {
		Fatal(funcName string, msg string)
		Info(funcName string, msg string)
		Error(funcName string, msg string)
		Debug(funcName string, msg string)
		Trace(funcName string, msg string)
	}

	// Initializer the interface that wrap log init function.
	Initializer interface {
		// InitialLog creates a logger object internally and
		// returns a log rotator object for optional extra configuration.
		InitialLog(level Level) Rotator
	}

	// Rotator the interface to wrap log rotation config.
	Rotator interface {
		SetupRotation(properties Properties) error
	}

	// Properties is log configuration properties
	Properties struct {
		Level           Level
		Path            string
		MaxAge          uint32
		MaxSize         uint32
		FilenamePattern string
		RotationTime    uint32
	}
)

func (l Level) String() string {
	str, ok := level2String[l]
	if !ok {
		return ""
	}
	return str
}

func GetLevel(in string) (Level, error) {
	level, ok := string2Level[in]
	if !ok {
		return 0, errors.New("invalid input")
	}
	return level, nil
}
