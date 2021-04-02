/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/1/21
 * Time: 10:19 PM
 *
 * Description:
 *
 */

package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorf(t *testing.T) {
	a := assert.New(t)
	params := []string{"v1", "v2"}
	msg := "this is test error, value1: %s, value2: %s"
	expected := fmt.Sprintf("common.TestErrorf-> %s", fmt.Sprintf(msg, params[0], params[1]))
	a.Equal(expected, Errorf("this is test error, value1: %s, value2: %s", params[0], params[1]).Error())
}

func TestError(t *testing.T) {
	a := assert.New(t)
	msg := "this is test error"
	expected := fmt.Sprintf("common.TestError-> %s", msg)
	a.Equal(expected, Error(msg).Error())
}
