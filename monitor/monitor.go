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
 * Time: 4:44 AM
 *
 * Description:
 *
 */

package monitor

type (

	// Monitor monitors a process.
	Monitor interface {
		// Start starts a process and monitor it. Most of the time this function
		// runs in a separate goroutine, because of that it accept a channel as
		// input argument.
		Start(chan<- string)
	}
)

