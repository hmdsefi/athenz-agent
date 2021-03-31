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
 * This file used for finding values, role and resource values.
 * This value can be in four different shape:
 * 		- exact or complete words: angler:role.public
 *		  for this type we need to use ZpeMatchEqual type
 *		- values that start with a specific characters: angler:role.pub*
 *		  for this type we will use ZpeMatchStartsWith type
 *		- using `*`, it means everything can be: angler:role.*
 *		  for this type we need to use ZpeMatchAll type
 *		- using regex for a value: angler:role.p?ub[a-z]
 *		  In this case we will use ZpeMatchRegex type
 *
 */

package matcher

import (
	"bytes"
	"regexp"
	"strings"
)

type ZpeMatch interface {
	Match(value string) bool
}

type ZpeMatchAll struct {
}

func (zma ZpeMatchAll) Match(value string) bool {
	return true
}

type ZpeMatchEqual struct {
	MatchValue string
}

func (zme ZpeMatchEqual) Match(value string) bool {
	return zme.MatchValue == value
}

type ZpeMatchStartsWith struct {
	Prefix string
}

func (zms ZpeMatchStartsWith) Match(value string) bool {
	return strings.HasPrefix(value, zms.Prefix)
}

type ZpeMatchRegex struct {
	Pattern *regexp.Regexp
}

func (zmr ZpeMatchRegex) Match(value string) bool {
	return zmr.Pattern.MatchString(value)
}

func isRegexMetaCharacter(regexChar string) bool {
	if regexChar == "^" || regexChar == "$" ||
		regexChar == "." || regexChar == "|" ||
		regexChar == "[" || regexChar == "+" ||
		regexChar == "(" || regexChar == ")" ||
		regexChar == "{" || regexChar == "\\" {
		return true
	}
	return false
}

func PatternFromGlob(glob string) string {
	var buffer bytes.Buffer
	buffer.WriteString("^")
	for _, char := range glob {
		strChar := string(char)
		if strChar == "*" {
			buffer.WriteString(".*")
		} else if strChar == "?" {
			buffer.WriteString(".")
		} else {
			if isRegexMetaCharacter(strChar) {
				buffer.WriteString("\\")
			}
			buffer.WriteString(strChar)
		}
	}
	buffer.WriteString("$")
	str := buffer.String()
	return str
}
