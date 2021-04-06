/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/6/21
 * Time: 7:38 AM
 *
 * Description:
 *
 */

package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevel_String(t *testing.T) {
	a := assert.New(t)
	a.Equal(Fatal.String(), "fatal")
	a.Equal(Error.String(), "error")
	a.Equal(Info.String(), "info")
	a.Equal(Debug.String(), "debug")
	a.Equal(Trace.String(), "trace")
}

func TestGetLevel(t *testing.T) {
	a := assert.New(t)
	a.Equal(Fatal, GetLevel("fatal"))
	a.Equal(Error, GetLevel("error"))
	a.Equal(Info, GetLevel("info"))
	a.Equal(Debug, GetLevel("deBug"))
	a.Equal(Trace, GetLevel("tracE"))
}
