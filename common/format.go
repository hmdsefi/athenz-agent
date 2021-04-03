/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/1/21
 * Time: 4:38 PM
 *
 * Description:
 *
 */

package common

import (
	"fmt"
	"log"
)

// Fatal wraps log.Fatal function to prevent importing golang log package.
func Fatal(in string)  {
	log.Fatalf("%s-> %s", CallerFuncName(),in)
}

// Fatalf wraps log.Fatalf function to prevent importing golang log package.
func Fatalf(format string, params ...interface{}){
	msg := fmt.Sprintf(format, params...)
	log.Fatalf("%s-> %s", CallerFuncName(), msg)
}

// Errorf wraps fmt.Errorf function to customize error message. Errorf adds its
// caller function name to the error message.
func Errorf(format string, params ...interface{}) error {
	msg := fmt.Sprintf(format, params...)
	return fmt.Errorf("%s-> %s", CallerFuncName(), msg)
}

// Error wraps errors.New function to customize error message. Error adds its
// caller function to the error message.
func Error(in string) error {
	return fmt.Errorf("%s-> %s", CallerFuncName(), in)
}
