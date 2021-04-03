/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 2/12/19
 * Time: 8:37 AM
 *
 * Description:
 *
 */

package matcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizePattern(t *testing.T) {
	a := assert.New(t)

	a.Equal("^abc$", normalizePattern("abc"))
	a.Equal("^abc.*$", normalizePattern("abc*"))
	a.Equal("^abc.$", normalizePattern("abc?"))
	a.Equal("^.*abc.$", normalizePattern("*abc?"))
	a.Equal("^.abc(c|d)$", normalizePattern("?abc(c|d)"))
	a.Equal("^abc\\.abc:.*$", normalizePattern("abc.abc:*"))
	a.Equal("^ab[a-c]c$", normalizePattern("ab[a-c]c"))
	a.Equal("^abc\\\\test\\\\$", normalizePattern("abc\\test\\"))
}

func TestNormalizePattern_regex(t *testing.T) {
	a := assert.New(t)

	a.Regexp(normalizePattern("?bcd(e|f)"), "abcde")
	a.NotRegexp(normalizePattern("?bcd(e|f)"), "bcde")
	a.Regexp(normalizePattern("*"), "abc12312")
	a.Regexp(normalizePattern("abc*"), "abcde")
}
