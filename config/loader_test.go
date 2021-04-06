/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/1/21
 * Time: 9:54 PM
 *
 * Description:
 *
 */

package config

import (
	"github.com/fsnotify/fsnotify"
	"testing"
)

func TestNotify(t *testing.T) {
	notify(fsnotify.Event{Name: "agent.conf"})
}
