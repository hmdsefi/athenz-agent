/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/10/21
 * Time: 3:52 PM
 *
 * Description:
 *
 */

package mock

import (
	"context"
	v1 "github.com/hamed-yousefi/athenz-agent/.gen/proto/api/message/v1"
)

type (
	AthenzAgentService struct{}
)

func (m AthenzAgentService) CheckAccessWithToken(ctx context.Context, request *v1.AccessCheckRequest) (*v1.AccessCheckResponse, error) {
	return &v1.AccessCheckResponse{AccessCheckStatus: v1.AccessStatus_DENY_DOMAIN_EMPTY}, nil
}

func (m AthenzAgentService) GetServiceToken(ctx context.Context, request *v1.ServiceTokenRequest) (*v1.ServiceTokenResponse, error) {
	panic("implement me")
}
