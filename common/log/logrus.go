/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Date: 3/31/21
 * Time: 9:08 PM
 *
 * Description:
 *
 */

package log

import (
	"github.com/hamed-yousefi/athenz-agent/common"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

const (
	logFilename = "log"
)

var (
	singleton sync.Once
	log       = new(logrusInitializer)
)

type (
	// logrusLogger is an implementation of Logger that using
	// logrus library underlying.
	logrusLogger struct {
		// the '*.go' file that logging is happening there.
		fileName string
		log      *logrusInitializer
	}

	// logrusInitializer is an implementation of Initializer for logrus.
	// It holds global logger instance. logrusLogger uses the global logger
	// that stored in logrusInitializer
	logrusInitializer struct {
		logger *logrus.Logger
	}

	// logrusLogRotator is an implementation of Rotator for logrus.
	logrusLogRotator struct {
		logrusInit *logrusInitializer
	}
)

// GetLogger returns a Logger for a specific `.go` file.
func GetLogger(fileName string) Logger {
	return &logrusLogger{fileName: fileName, log: log}
}

func (l *logrusLogger) Fatal(msg string) {
	l.log.logger.WithFields(logrus.Fields{
		"FileName": l.fileName,
		"Func":     common.CallerFuncName(),
	}).Fatal(msg)
}

func (l *logrusLogger) Fatalf(format string, params ...interface{}) {
	l.log.logger.WithFields(logrus.Fields{
		"FileName": l.fileName,
		"Func":     common.CallerFuncName(),
	}).Fatalf(format, params...)
}

func (l *logrusLogger) Info(msg string) {
	l.log.logger.WithFields(logrus.Fields{
		"FileName": l.fileName,
		"Func":     common.CallerFuncName(),
	}).Info(msg)
}

func (l *logrusLogger) Error(msg string) {
	l.log.logger.WithFields(logrus.Fields{
		"FileName": l.fileName,
		"Func":     common.CallerFuncName(),
	}).Error(msg)
}

func (l *logrusLogger) Debug(msg string) {
	l.log.logger.WithFields(logrus.Fields{
		"FileName": l.fileName,
		"Func":     common.CallerFuncName(),
	}).Debug(msg)
}

func (l *logrusLogger) Trace(msg string) {
	l.log.logger.WithFields(logrus.Fields{
		"FileName": l.fileName,
		"Func":     common.CallerFuncName(),
	}).Trace(msg)
}

func NewLogrusInitializer() Initializer {
	return log
}

// InitialLog initializes logrus log.
//
// InitialLog Accepts  log level as input param and will creates
// a logrus instance. This initialization happens just once in entire
// application lifecycle.
func (l *logrusInitializer) InitialLog(level Level) Rotator {
	// initial logrus log just once in entire application lifecycle
	singleton.Do(func() {
		log.logger = logrus.New()
		log.logger.SetFormatter(&logrus.JSONFormatter{})
		lvl, err := logrus.ParseLevel(level.String())
		if err != nil {
			log.logger.Fatal("unable to parse input log level to logrus log level")
		}
		log.logger.SetLevel(lvl)

		// set default output
		log.logger.SetOutput(os.Stdout)
	})

	// create and return log rotator
	return &logrusLogRotator{
		logrusInit: log,
	}
}

// SetupRotation creates a custom output writer for log.
func (r *logrusLogRotator) SetupRotation(provider common.LogConfigProvider) {

	// create log directory
	if err := common.CreateAllDirectories(provider.GetPath()); err != nil {
		r.logrusInit.logger.Fatal(err)
	}

	// log rotation config
	writer, err := rotateLogs.New(
		provider.GetPath()+string(os.PathSeparator)+logFilename+provider.GetFilenamePattern(),
		rotateLogs.WithLinkName(provider.GetPath()+string(os.PathSeparator)+logFilename),
		rotateLogs.WithMaxAge(provider.GetMaxAge()),
		rotateLogs.WithRotationSize(provider.GetMaxSize()),
		rotateLogs.WithRotationTime(provider.GetRotationTime()),
	)
	if err != nil {
		r.logrusInit.logger.Fatal(err)
	}

	// set new output writer, consul logs are crucial so lets
	// write logs in both file and stout
	multiWriter := io.MultiWriter(writer, os.Stdout)
	r.logrusInit.logger.SetOutput(multiWriter)
}
