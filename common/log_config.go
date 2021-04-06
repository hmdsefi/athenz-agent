/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/1/21
 * Time: 3:28 PM
 *
 * Description:
 *
 */

package common

import "time"

type (

	// LogConfigProvider is the interface that wraps log configs. It helps config
	// package to prevent using log package objects and vice versa. This interface
	// provides one directional dependency and OCP principle.
	LogConfigProvider interface {
		// GetLevel returns log level.
		GetLevel() string
		// GetPath returns the path that log files must be stored there.
		GetPath() string
		// GetMaxAge returns the max age of a log file before it gets purged from the file system.
		GetMaxAge() time.Duration
		// GetRotationTime return the time between rotation.
		GetRotationTime() time.Duration
		// GetMaxSize returns the log file size between rotation.
		GetMaxSize() int64
		// GetFilenamePattern returns filename pattern.
		GetFilenamePattern() string
	}
)
