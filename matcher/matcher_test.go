/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
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

func TestPatternFromGlob(t *testing.T) {
	a := assert.New(t)

	a.Equal(PatternFromGlob("abc"), "^abc$")
	a.Equal(PatternFromGlob("abc*"), "^abc.*$")
	a.Equal(PatternFromGlob("abc?"), "^abc.$")
	a.Equal(PatternFromGlob("*abc?"), "^.*abc.$")
	a.Equal(PatternFromGlob("abc.abc:*"), "^abc\\.abc:.*$")
	a.Equal(PatternFromGlob("ab[a-c]c"), "^ab\\[a-c]c$")
	a.Equal(PatternFromGlob("ab*.()^$c"), "^ab.*\\.\\(\\)\\^\\$c$")
	a.Equal(PatternFromGlob("abc\\test\\"), "^abc\\\\test\\\\$")
	a.Equal(PatternFromGlob("ab{|c+"), "^ab\\{\\|c\\+$")
	a.Equal(PatternFromGlob("^$[()\\+{.*?|"), "^\\^\\$\\[\\(\\)\\\\\\+\\{\\..*.\\|$")
}
