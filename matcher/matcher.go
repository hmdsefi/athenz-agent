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

// ZpeMatch is a mechanism to see if a string match with a pattern or not.
type ZpeMatch interface {
	// Match checks if the input string is a match with inner pattern or not.
	Match(value string) bool
}

// ZpeMatchAll is an implementation of ZpeMatch. It matches with all input strings.
type ZpeMatchAll struct {
}

// Match returns true for any input string
func (zma ZpeMatchAll) Match(value string) bool {
	return true
}

// ZpeMatchEqual is an implementation of ZpeMatch.
type ZpeMatchEqual struct {
	MatchValue string
}

// Match returns true if input string is equal to MatchValue.
func (zme ZpeMatchEqual) Match(value string) bool {
	return zme.MatchValue == value
}

// ZpeMatchStartsWith is an implementation of ZpeMatch for matching strings
// with a specific prefix.
type ZpeMatchStartsWith struct {
	Prefix string
}

// Match returns true if input string starts with Prefix.
func (zms ZpeMatchStartsWith) Match(value string) bool {
	return strings.HasPrefix(value, zms.Prefix)
}

// ZpeMatchRegex is an implementation of ZpeMatch for matching regex pattern.
type ZpeMatchRegex struct {
	Pattern *regexp.Regexp
}

// Match returns true if input string matches with Pattern.
func (zmr ZpeMatchRegex) Match(value string) bool {
	return zmr.Pattern.MatchString(value)
}

// NewZpeMatchRegex returns new instance of ZpeMatch type. It would return an
// error if the regex pattern were malformed.
func NewZpeMatchRegex(in string) (ZpeMatch, error) {
	patter, err := regexp.Compile(normalizePattern(in))
	if err != nil {
		return nil, err
	}

	return ZpeMatchRegex{Pattern: patter}, nil
}

func isRegexMetaCharacter(regexChar string) bool {
	if regexChar == "^" || regexChar == "$" ||
		regexChar == "." || regexChar == "\\" {
		return true
	}
	return false
}

func normalizePattern(in string) string {
	var buffer bytes.Buffer
	buffer.WriteString("^")
	for _, char := range in {
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
