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
		GetLevel() string
		GetPath() string
		GetMaxAge() time.Duration
		GetRotationTime() time.Duration
		GetMaxSize() int64
		GetFilenamePattern() string
	}
)
